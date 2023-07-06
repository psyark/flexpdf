package flexpdf

import (
	"image/color"
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
	flexItemCommon

	Direction      Direction
	JustifyContent JustifyContent
	AlignItems     AlignItems
	Items          []FlexItem
}

func NewBox(dir Direction, items ...FlexItem) *Box {
	return &Box{
		flexItemCommon: flexItemCommonDefault,
		Direction:      dir,
		Items:          items,
	}
}
func (b *Box) SetWidth(w float64) *Box {
	b.Width = w
	return b
}
func (b *Box) SetHeight(h float64) *Box {
	b.Height = h
	return b
}
func (b *Box) SetSize(w, h float64) *Box {
	return b.SetWidth(w).SetHeight(h)
}
func (b *Box) SetBackgroundColor(c color.Color) *Box {
	b.BackgroundColor = c
	return b
}
func (b *Box) SetBorder(border Border) *Box {
	b.Border = border
	return b
}
func (b *Box) SetJustifyContent(jc JustifyContent) *Box {
	b.JustifyContent = jc
	return b
}
func (b *Box) draw(pdf *gopdf.GoPdf, r rect, depth int) error {
	log.Printf("%sBox.draw(r=%v, d=%v jc=%v ai=%v)\n", strings.Repeat("  ", depth), r, b.Direction, b.JustifyContent, b.AlignItems)

	// 背景色
	if b.BackgroundColor != nil && r.w != 0 && r.h != 0 {
		if err := setColor(pdf, b.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(r.x, r.y, r.x+r.w, r.y+r.h, "F", 0, 0); err != nil {
			return errors.Wrap(err, "rectangle")
		}
	}

	// 子孫
	itemRect := r
	prefSizes := make([]*size, len(b.Items))
	mainAxisRemains := r.getLength(b.Direction.mainAxis())
	for i, item := range b.Items {
		ps, err := item.getPreferredSize(pdf)
		if err != nil {
			return err
		}
		prefSizes[i] = ps
		mainAxisRemains -= ps.getLength(b.Direction.mainAxis())
	}

	if mainAxisRemains < 0 {
		mainAxisRemains = 0
	}

	if b.Direction.mainAxis() == horizontal {
		switch b.JustifyContent {
		case JustifyContentFlexEnd:
			itemRect.x += mainAxisRemains
		case JustifyContentCenter:
			itemRect.x += mainAxisRemains / 2
		case JustifyContentSpaceAround:
			itemRect.x += mainAxisRemains / float64(len(b.Items)*2)
		}
	} else {
		// TODO
	}

	for i, item := range b.Items {
		ps := prefSizes[i]

		if b.Direction.mainAxis() == horizontal {
			itemRect.w = ps.w
		} else {
			itemRect.h = ps.h
		}

		if err := item.draw(pdf, itemRect, depth+1); err != nil {
			return errors.Wrap(err, "item.draw")
		}

		if b.Direction.mainAxis() == horizontal {
			itemRect.x += ps.w
			switch b.JustifyContent {
			case JustifyContentSpaceBetween:
				itemRect.x += mainAxisRemains / float64(len(b.Items)-1)
			case JustifyContentSpaceAround:
				itemRect.x += mainAxisRemains / float64(len(b.Items))
			}
		} else {
			// TODO
			itemRect.y += ps.h
		}
	}

	if err := b.Border.draw(pdf, r); err != nil {
		return err
	}

	return nil
}
func (b *Box) getPreferredSize(pdf *gopdf.GoPdf) (*size, error) {
	ps := &size{}
	for _, item := range b.Items {
		ips, err := item.getPreferredSize(pdf)
		if err != nil {
			return nil, err
		}
		if b.Direction.mainAxis() == horizontal {
			ps.w += ips.w
			ps.h = math.Max(ps.h, ips.h)
		} else {
			ps.w = math.Max(ps.w, ips.w)
			ps.h += ips.h
		}
	}
	if b.Width >= 0 {
		ps.w = b.Width
	}
	if b.Height >= 0 {
		ps.h = b.Height
	}
	return ps, nil
}
