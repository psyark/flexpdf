package flexpdf

import (
	"image/color"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

type TextAlign int

const (
	TextAlignBegin TextAlign = iota
	TextAlignCenter
	TextAlignEnd
)

// Text はテキストを扱うエレメントです
// TODO family, size, color, text は Spanのスライスにする
type Text struct {
	flexItemCommon[*Text]

	FontFamily string
	FontSize   float64
	Text       string
	Color      color.Color
	LineHeight float64
	Align      TextAlign
}

func (t *Text) SetLineHeight(lineHeight float64) *Text {
	t.LineHeight = lineHeight
	return t
}
func (t *Text) SetAlign(align TextAlign) *Text {
	t.Align = align
	return t
}

func NewText(family string, size float64, text string) *Text {
	t := &Text{
		FontFamily: family,
		FontSize:   size,
		Text:       text,
		LineHeight: 1,
	}
	t.flexItemCommon.init(t)
	return t
}

func (t *Text) drawContent(pdf *gopdf.GoPdf, r rect, depth int) error {
	log.Printf("%sText.draw(r=%v, t=%q)\n", strings.Repeat("  ", depth), r, t.Text)

	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return errors.Wrap(err, "setFont")
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
	// TODO 幅が小さすぎる場合に無限ループになるのを抑制
	if r.w < 20 {
		return nil
	}

	lines, err := pdf.SplitText(t.Text, r.w)
	if err != nil {
		return errors.Wrap(err, "splitText")
	}

	for _, line := range lines {
		opt := gopdf.CellOption{}
		switch t.Align {
		case TextAlignBegin:
			opt.Align = gopdf.Left
		case TextAlignCenter:
			opt.Align = gopdf.Center
		case TextAlignEnd:
			opt.Align = gopdf.Right
		}

		if err := pdf.MultiCellWithOption(&gopdf.Rect{W: r.w, H: r.h}, line, opt); err != nil {
			return errors.Wrap(err, "multiCell")
		}
		pdf.Br((t.LineHeight - 1) * t.FontSize)
		pdf.SetX(r.x)
	}

	return nil
}

func (t *Text) getContentSize(pdf *gopdf.GoPdf, width float64) (size, error) {
	if err := pdf.SetFont(t.FontFamily, "", t.FontSize); err != nil {
		return size{}, err
	}

	if width < 0 { // 負の場合、幅が制限されないときのサイズを調べる
		width = 10000000
	}

	lines, err := pdf.SplitText(t.Text, width)
	if err != nil {
		return size{}, err
	}

	cs := size{
		h: t.FontSize * (float64(len(lines)) + (t.LineHeight-1)*float64(len(lines)-1)),
	}

	for _, line := range lines {
		w, err := pdf.MeasureTextWidth(line)
		if err != nil {
			return size{}, err
		}
		if cs.w < w {
			cs.w = w
		}
	}

	return cs, nil
}

type TextSpan struct {
}
