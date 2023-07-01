package flexpdf

import (
	"github.com/signintech/gopdf"
)

type Box struct {
	Direction      Direction
	JustifyContent JustifyContent
	AlignItems     AlignItems
	Items          []FlexItem
}

func (b *Box) draw(pdf *gopdf.GoPdf, rect rect) error {
	if err := pdf.SetFont("ipaexg", "", 20); err != nil {
		return err
	}

	pdf.SetXY(rect.x, rect.y)
	pdf.SetLineWidth(4)
	return pdf.CellWithOption(&gopdf.Rect{W: rect.w, H: rect.h}, "HOGE", gopdf.CellOption{Border: gopdf.AllBorders})
}
