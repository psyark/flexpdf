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

func (b *Box) draw(pdf *gopdf.GoPdf, rect rect) error {
	if err := pdf.SetFont("ipaexg", "", 20); err != nil {
		return err
	}

	// 背景色
	if b.BackgroundColor != nil {
		if err := setColor(pdf, b.BackgroundColor); err != nil {
			return err
		}
		if err := pdf.Rectangle(rect.x, rect.y, rect.x+rect.w, rect.y+rect.h, "F", 0, 0); err != nil {
			return err
		}
	}

	log.Println("dir:", b.Direction)
	// 子孫
	for _, item := range b.Items {
		if err := item.draw(pdf, rect); err != nil {
			return err
		}
	}

	return nil
}

func setColor(pdf *gopdf.GoPdf, col color.Color) error {
	r, g, b, a := col.RGBA()
	if a == 0xFFFF {
		pdf.ClearTransparency()
	} else if err := pdf.SetTransparency(gopdf.Transparency{Alpha: float64(a) / 0xFFFF, BlendModeType: gopdf.NormalBlendMode}); err != nil {
		return err
	}

	pdf.SetFillColor(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return nil
}
