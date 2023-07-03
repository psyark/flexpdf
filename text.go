package flexpdf

import (
	"image/color"

	"github.com/pkg/errors"
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
	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return errors.Wrap(err, "setFont")
	}

	ps, err := t.getPreferredSize(pdf)
	if err != nil {
		return errors.Wrap(err, "getPreferredSize")
	}

	{
		c := t.Color
		if c == nil {
			c = color.Black
		}
		if err := setColor(pdf, c); err != nil {
			return errors.Wrap(err, "setColor")
		}
	}

	pdf.SetXY(r.x, r.y)
	lines, err := pdf.SplitTextWithWordWrap(t.Text, 10000000)
	if err != nil {
		return errors.Wrap(err, "splitTextWithWordWrap")
	}

	for _, line := range lines {
		if err := pdf.MultiCell(&gopdf.Rect{W: 10000000, H: 100000}, line); err != nil {
			return errors.Wrap(err, "multiCell")
		}
		if t.LineHeight != 0 && t.LineHeight != 1 {
			pdf.Br((t.LineHeight - 1) * t.FontSize)
			pdf.SetX(r.x)
		}
	}

	{ // TODO デバッグ用
		pdf.SetLineType("dotted")
		if err := pdf.Rectangle(r.x, r.y, r.x+ps.w, r.y+ps.h, "D", 0, 0); err != nil {
			return errors.Wrap(err, "rectangle")
		}
	}

	return nil
}

func (t *Text) getPreferredSize(pdf *gopdf.GoPdf) (*size, error) {
	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return nil, err
	}

	lines, err := pdf.SplitTextWithWordWrap(t.Text, 10000000)
	if err != nil {
		return nil, err
	}

	ps := &size{
		h: t.FontSize * float64(len(lines)),
	}
	if t.LineHeight != 0 {
		ps.h += (t.FontSize * (t.LineHeight - 1)) * float64(len(lines)-1)
	}

	for _, line := range lines {
		w, err := pdf.MeasureTextWidth(line)
		if err != nil {
			return nil, err
		}
		if ps.w < w {
			ps.w = w
		}
	}
	return ps, nil
}

type TextSpan struct {
}
