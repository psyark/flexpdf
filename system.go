package flexpdf

import (
	"fmt"
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
	// draw はこのFlexItemを与えられた矩形内に描画します。
	draw(pdf *gopdf.GoPdf, r rect) error
	getPreferredSize(pdf *gopdf.GoPdf, marginBoxMax size) (size, error)
	getFlexGrow() float64
}

func setColor(pdf *gopdf.GoPdf, col color.Color) (err error) {
	defer wrap(&err, "setColor")

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

// func (s Spacing) w() float64 {
// 	return s.Right + s.Left
// }
// func (s Spacing) h() float64 {
// 	return s.Top + s.Bottom
// }

func wrap(errp *error, format string, args ...any) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
	}
}
