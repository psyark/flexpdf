package flexpdf

import (
	"image/color"
	"log"
	"math"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

// FlexBasis は持たず、 Width/Heightでサイズが指定してあればそのサイズ（basis=auto同等）
// Width/Heightが指定してなければ子要素のサイズ(basis=content)となる
type Box struct {
	Direction       Direction
	Width           *float64
	Height          *float64
	Border          Border
	FlexGrow        float64
	FlexShrink      float64
	JustifyContent  JustifyContent
	AlignItems      AlignItems
	BackgroundColor color.Color
	Items           []FlexItem
}

func NewBox() *Box {
	return &Box{FlexGrow: 0, FlexShrink: 1}
}

func (b *Box) draw(pdf *gopdf.GoPdf, r rect) error {
	// 背景色
	if b.BackgroundColor != nil && r.w != 0 && r.h != 0 {
		if err := setColor(pdf, b.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(r.x, r.y, r.x+r.w, r.y+r.h, "F", 0, 0); err != nil {
			return errors.Wrap(err, "rectangle")
		}
	}

	log.Printf("Direction=%q JustifyContent=%q AlignItems=%q\n", b.Direction, b.JustifyContent, b.AlignItems)

	// 子孫
	itemRect := r
	prefSizes := make([]*size, len(b.Items))
	mainAxisRemains := r.getLength(b.Direction.mainAxis())
	log.Println(mainAxisRemains)
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

		if err := item.draw(pdf, itemRect); err != nil {
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
	if b.Width != nil {
		ps.w = *b.Width
	}
	if b.Height != nil {
		ps.h = *b.Height
	}
	return ps, nil
}
