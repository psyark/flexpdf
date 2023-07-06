package flexpdf

import (
	"log"
	"math"
	"strings"

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
func (b *Box) drawContent(pdf *gopdf.GoPdf, r rect, depth int) error {
	log.Printf("%sBox.draw(r=%v, d=%v jc=%v ai=%v)\n", strings.Repeat("  ", depth), r, b.Direction, b.JustifyContent, b.AlignItems)

	// 子孫
	itemRect := r
	prefSizes := make([]*size, len(b.Items))

	var spacing float64
	{
		var growing, growTotal float64
		mainAxisRemains := r.getLength(b.Direction.mainAxis())

		// 1パス目は自然なサイズ
		for i, item := range b.Items {
			maxWidth := -1.0 // -1 -> 自然なサイズ
			if b.Direction.mainAxis() == vertical {
				maxWidth = r.w
			}
			ps, err := item.getPreferredSize(pdf, maxWidth)
			if err != nil {
				return err
			}
			growTotal += item.getFlexGrow()
			prefSizes[i] = ps
			mainAxisRemains -= ps.getLength(b.Direction.mainAxis())
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

		// 2パス目は幅を制限したときのサイズ
		for i, item := range b.Items {
			ps := prefSizes[i]
			if growTotal != 0 {
				ps.w += growing * item.getFlexGrow() / growTotal
			}

			ps2, err := item.getPreferredSize(pdf, ps.w) // 指定したサイズ
			if err != nil {
				return err
			}
			ps.h = ps2.h // 高さだけ更新
			prefSizes[i] = ps
		}
	}

	if b.Direction.mainAxis() == horizontal {
		switch b.JustifyContent {
		case JustifyContentFlexEnd:
			itemRect.x += spacing
		case JustifyContentCenter:
			itemRect.x += spacing / 2
		case JustifyContentSpaceAround:
			itemRect.x += spacing / float64(len(b.Items)*2)
		}
	} else {
		// TODO
	}

	for i, item := range b.Items {
		ps := prefSizes[i]

		itemRect.w = ps.w
		itemRect.h = ps.h

		if b.Direction.mainAxis() == horizontal {
			if b.AlignItems == AlignItemsStretch {
				itemRect.h = r.h
			}
		} else {
			if b.AlignItems == AlignItemsStretch {
				itemRect.w = r.w
			}
		}

		if err := item.draw(pdf, itemRect, depth+1); err != nil {
			return errors.Wrap(err, "item.draw")
		}

		if b.Direction.mainAxis() == horizontal {
			itemRect.x += ps.w
			switch b.JustifyContent {
			case JustifyContentSpaceBetween:
				itemRect.x += spacing / float64(len(b.Items)-1)
			case JustifyContentSpaceAround:
				itemRect.x += spacing / float64(len(b.Items))
			}
		} else {
			// TODO
			itemRect.y += ps.h
		}
	}

	return nil
}
func (b *Box) getContentSize(pdf *gopdf.GoPdf, _ float64) (*size, error) {
	cs := &size{}
	for _, item := range b.Items {
		ips, err := item.getPreferredSize(pdf, -1)
		if err != nil {
			return nil, err
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
