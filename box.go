package flexpdf

import (
	"image/color"
	"log"
	"math"

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

	log.Printf("Direction=%q JustifyContent=%q AlignItems=%q\n", b.Direction, b.JustifyContent, b.AlignItems)

	// 子孫
	itemRect := rect{x: r.x, y: r.y}
	for _, item := range b.Items {
		ps, err := item.getPreferredSize(pdf)
		if err != nil {
			return err
		}

		log.Println(ps.w, ps.h)

		itemRect.w = ps.w
		itemRect.h = ps.h

		if err := item.draw(pdf, itemRect); err != nil {
			return err
		}

		if isHorizontal(b.Direction) {
			itemRect.x += ps.w
		} else {
			itemRect.y += ps.h
		}
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
		if isHorizontal(b.Direction) {
			ps.w += ips.w
			ps.h = math.Max(ps.h, ips.h)
		} else {
			ps.w = math.Max(ps.w, ips.w)
			ps.h += ips.h
		}
	}
	return ps, nil
}
