package flexpdf

import (
	"image/color"

	"github.com/signintech/gopdf"
)

type FlexItem interface {
	draw(*gopdf.GoPdf, rect) error
	getPreferredSize(*gopdf.GoPdf) (*size, error)
}

type size struct {
	w float64
	h float64
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
