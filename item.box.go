package flexpdf

import (
	"math"

	"github.com/pkg/errors"
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
func (b *Box) drawContent(pdf *gopdf.GoPdf, r rect, depth int) error {
	// log.Printf("%sBox.draw(r=%v, d=%v jc=%v ai=%v)\n", strings.Repeat("  ", depth), r, b.Direction, b.JustifyContent, b.AlignItems)

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
			maxWidth := -1.0 // -1 -> 自然なサイズ
			if mainAxis == vertical {
				maxWidth = r.w
			}
			ps, err := item.getPreferredSize(pdf, maxWidth)
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

		// 2パス目は幅を制限したときのサイズ
		for i, item := range b.Items {
			ps := prefSizes[i]

			// グロー
			if growTotal != 0 {
				ps = ps.update(mainAxis, func(f float64) float64 {
					return f + growing*item.getFlexGrow()/growTotal
				})
			}

			if mainAxis == horizontal {
				ps2, err := item.getPreferredSize(pdf, ps.w) // 指定したサイズ
				if err != nil {
					return err
				}
				ps.h = ps2.h // 高さだけ更新
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
		if err := item.draw(pdf, itemRect, depth+1); err != nil {
			return errors.Wrap(err, "item.draw")
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
func (b *Box) getContentSize(pdf *gopdf.GoPdf, _ float64) (size, error) {
	cs := size{}
	for _, item := range b.Items {
		ips, err := item.getPreferredSize(pdf, -1)
		if err != nil {
			return size{}, err
		}
		if b.Direction.mainAxis() == horizontal {
			cs.w += ips.w
			cs.h = math.Max(cs.h, ips.h)
		} else {
			cs.w = math.Max(cs.w, ips.w)
			cs.h += ips.h
		}
	}
	return cs, nil
}
