package flexpdf

import (
	"image/color"

	"github.com/signintech/gopdf"
)

type Text struct {
	FontFamily string
	FontSize   float64
	Text       string
	Color      color.Color
	LineHeight float64
}

func (t *Text) draw(pdf *gopdf.GoPdf, r rect) error {
	w, h, err := t.getSize(pdf)
	if err != nil {
		return err
	}

	{
		c := t.Color
		if c == nil {
			c = color.Black
		}
		if err := setColor(pdf, c); err != nil {
			return err
		}
	}

	pdf.SetXY(r.x, r.y)
	lines, err := pdf.SplitTextWithWordWrap(t.Text, 10000000)
	if err != nil {
		return err
	}

	for _, line := range lines {
		pdf.MultiCell(&gopdf.Rect{W: 10000000, H: 100000}, line)
		if t.LineHeight != 0 && t.LineHeight != 1 {
			pdf.Br((t.LineHeight - 1) * t.FontSize)
			pdf.SetX(r.x)
		}
	}

	pdf.Rectangle(r.x, r.y, r.x+w, r.y+h, "D", 0, 0)

	return nil
}

func (t *Text) getSize(pdf *gopdf.GoPdf) (float64, float64, error) {
	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return 0, 0, err
	}

	// pdf.MeasureTextWidth()

	lines, err := pdf.SplitTextWithWordWrap(t.Text, 10000000)
	if err != nil {
		return 0, 0, err
	}

	h := t.FontSize * float64(len(lines))
	if t.LineHeight != 0 {
		h += (t.FontSize * (t.LineHeight - 1)) * float64(len(lines)-1)
	}

	var mw float64
	for _, line := range lines {
		w, err := pdf.MeasureTextWidth(line)
		if err != nil {
			return 0, 0, err
		}
		if mw < w {
			mw = w
		}
	}
	return mw, h, nil
}

type TextSpan struct {
}
