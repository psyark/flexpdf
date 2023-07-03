package flexpdf

import (
	"image/color"
	"log"

	"github.com/signintech/gopdf"
)

type Box struct {
	Direction       Direction
	JustifyContent  JustifyContent
	AlignItems      AlignItems
	BackgroundColor color.Color
	Items           []FlexItem
}

func (b *Box) draw(pdf *gopdf.GoPdf, r rect) error {
	if err := pdf.SetFont("ipaexg", "", 20); err != nil {
		return err
	}

	// 背景色
	if b.BackgroundColor != nil {
		if err := setColor(pdf, b.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(r.x, r.y, r.x+r.w, r.y+r.h, "F", 0, 0); err != nil {
			return err
		}
	}

	log.Println("dir:", b.Direction)
	// 子孫

	x := r.x
	for _, item := range b.Items {
		ps, err := item.getPreferredSize(pdf)
		if err != nil {
			return err
		}

		log.Println(ps.w, ps.h)

		if err := item.draw(pdf, rect{x: x, y: r.y, w: ps.w, h: ps.h}); err != nil {
			return err
		}

		x += ps.w
	}

	return nil
}
