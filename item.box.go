package flexpdf

import (
	"math"

	"github.com/signintech/gopdf"
)

// FlexBasis は持たず、 Width/Heightでサイズが指定してあればそのサイズ（basis=auto同等）
// Width/Heightが指定してなければ子要素のサイズ(basis=content)となる
type Box struct {
	// 共通フィールド
	flexItemCommon[*Box]

	Direction      Direction
	JustifyContent JustifyContent
	AlignItems     AlignItems
	Items          []FlexItem
}

func NewColumnBox(items ...FlexItem) *Box {
	return NewBox(DirectionColumn, items...)
}
func NewRowBox(items ...FlexItem) *Box {
	return NewBox(DirectionRow, items...)
}

func NewBox(dir Direction, items ...FlexItem) *Box {
	b := &Box{
		Direction:  dir,
		Items:      items,
		AlignItems: AlignItemsStretch,
	}
	b.flexItemCommon.init(b)
	return b
}
func (b *Box) SetJustifyContent(jc JustifyContent) *Box {
	b.JustifyContent = jc
	return b
}
func (b *Box) SetAlignItems(aa AlignItems) *Box {
	b.AlignItems = aa
	return b
}
func (b *Box) drawContent(pdf *gopdf.GoPdf, r rect) (err error) {
	defer wrap(&err, "box.drawContent")

	mainAxis := b.Direction.mainAxis()
	counterAxis := !mainAxis

	// 子孫
	itemRect := r
	prefSizes := make([]size, len(b.Items))

	var spacing float64
	{
		var growing, growTotal float64
		mainAxisRemains := r.getLength(mainAxis)

		// 1パス目は自然なサイズ
		for i, item := range b.Items {
			ps, err := item.getPreferredSize(pdf, size{w: r.w, h: r.h})
			if err != nil {
				return err
			}
			growTotal += item.getFlexGrow()
			prefSizes[i] = ps
			mainAxisRemains -= ps.get(mainAxis)
		}

		if mainAxisRemains < 0 {
			mainAxisRemains = 0
		}
		if growTotal >= 1 {
			growing = mainAxisRemains
			spacing = 0
		} else {
			growing = mainAxisRemains * growTotal
			spacing = mainAxisRemains - growing
		}

		// log.Println(mainAxisRemains, growTotal, growing, spacing)

		// 2パス目はグロー・シュリンクを考慮したサイズ
		for i, item := range b.Items {
			ps := prefSizes[i]

			// グロー
			if growTotal != 0 {
				ps = ps.update(mainAxis, func(f float64) float64 {
					return f + growing*item.getFlexGrow()/growTotal
				})
			}

			ps_, err := item.getPreferredSize(pdf, ps) // グロー・シュリンクしたサイズ
			if err != nil {
				return err
			}
			{
				// TODO ここの仕様を明確にする
				// なぜmainAxisで分岐している？
				// なぜMaxしている？
				// グローの再計算は不要？
				if mainAxis == horizontal {
					ps.w = math.Max(ps.w, ps_.w)
				} else {
					ps.h = math.Max(ps.h, ps_.h)
				}
			}

			prefSizes[i] = ps
		}
	}

	// 開始位置
	switch b.JustifyContent {
	case JustifyContentFlexEnd:
		itemRect = itemRect.updatePos(mainAxis, func(v float64) float64 {
			return v + spacing
		})
	case JustifyContentCenter:
		itemRect = itemRect.updatePos(mainAxis, func(v float64) float64 {
			return v + spacing/2
		})
	case JustifyContentSpaceAround:
		itemRect = itemRect.updatePos(mainAxis, func(v float64) float64 {
			return v + spacing/float64(len(b.Items)*2)
		})
	}

	for i, item := range b.Items {
		ps := prefSizes[i]
		itemRect.w = ps.w
		itemRect.h = ps.h

		// ストレッチ
		if b.AlignItems == AlignItemsStretch {
			itemRect = itemRect.setSize(counterAxis, r.getSize(counterAxis))
		}

		// 描画
		if err := item.draw(pdf, itemRect); err != nil {
			return err
		}

		// アイテム間の余白
		itemRect = itemRect.updatePos(mainAxis, func(v float64) float64 {
			return v + ps.get(mainAxis)
		})
		switch b.JustifyContent {
		case JustifyContentSpaceBetween:
			itemRect = itemRect.updatePos(mainAxis, func(v float64) float64 {
				return v + spacing/float64(len(b.Items)-1)
			})
		case JustifyContentSpaceAround:
			itemRect = itemRect.updatePos(mainAxis, func(v float64) float64 {
				return v + spacing/float64(len(b.Items))
			})
		}
	}

	return nil
}
func (b *Box) getContentSize(pdf *gopdf.GoPdf, contentBoxMax size) (size, error) {

	cs := size{}
	for _, item := range b.Items {
		ips, err := item.getPreferredSize(pdf, contentBoxMax)
		if err != nil {
			return size{}, err
		}

		mainAxis := b.Direction.mainAxis()
		counterAxis := !mainAxis

		cs = cs.add(mainAxis, ips.get(mainAxis))
		cs = cs.update(counterAxis, func(ov float64) float64 {
			return math.Max(ov, ips.get(counterAxis))
		})
	}

	return cs, nil
}
