package flexpdf

import (
	"image/color"

	"github.com/signintech/gopdf"
)

type axis bool

const (
	horizontal axis = false
	vertical   axis = true
)

var (
	_ FlexItem = &Text{}
	_ FlexItem = &Box{}
)

type FlexItem interface {
	draw(pdf *gopdf.GoPdf, r rect, depth int) error
	getPreferredSize(pdf *gopdf.GoPdf, width float64) (size, error)
	getFlexGrow() float64
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

type TRBL[T any] struct {
	Top    T
	Right  T
	Bottom T
	Left   T
}

type Spacing TRBL[float64]
