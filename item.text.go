package flexpdf

import (
	"image/color"

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

	LineHeight float64
	Align      TextAlign
	Runs       []*TextRun
}

type TextRun struct {
	Color      color.Color
	FontSize   float64
	FontFamily string
	Text       string
}

func NewRun(text string) *TextRun {
	return &TextRun{
		Color:      color.Black,
		FontSize:   10,
		FontFamily: "",
		Text:       text,
	}
}
func (r *TextRun) SetColor(c color.Color) *TextRun {
	r.Color = c
	return r
}
func (r *TextRun) SetFontSize(s float64) *TextRun {
	r.FontSize = s
	return r
}
func (r *TextRun) SetFontFamily(f string) *TextRun {
	r.FontFamily = f
	return r
}
func (r *TextRun) SetText(t string) *TextRun {
	r.Text = t
	return r
}

func (t *Text) AddRun(run *TextRun) *Text {
	t.Runs = append(t.Runs, run)
	return t
}
func (t *Text) SetLineHeight(lineHeight float64) *Text {
	t.LineHeight = lineHeight
	return t
}
func (t *Text) SetAlign(align TextAlign) *Text {
	t.Align = align
	return t
}

func NewText(runs ...*TextRun) *Text {
	t := &Text{
		LineHeight: 1,
		Runs:       runs,
	}
	t.flexItemCommon.init(t)
	return t
}

func (t *Text) drawContent(pdf *gopdf.GoPdf, r rect, depth int) (err error) {
	defer wrap(&err, "text.drawContent")

	for _, run := range t.Runs {
		if err := pdf.SetFont(run.FontFamily, "", run.FontSize); err != nil {
			return err
		}

		{
			c := run.Color
			if c == nil {
				c = color.Black
			}
			if err := setColor(pdf, c); err != nil {
				return err
			}
		}

		pdf.SetXY(r.x, r.y)
		// TODO 幅が小さすぎる場合に無限ループになるのを抑制
		if r.w < 20 {
			return nil
		}

		lines, err := pdf.SplitText(run.Text, r.w)
		if err != nil {
			return err
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
				return err
			}
			pdf.Br((t.LineHeight - 1) * run.FontSize)
			pdf.SetX(r.x)
		}
	}

	return nil
}

func (t *Text) getContentSize(pdf *gopdf.GoPdf, width float64) (s size, err error) {
	defer wrap(&err, "text.getContentSize")

	run := t.Runs[0]

	if err := pdf.SetFont(run.FontFamily, "", run.FontSize); err != nil {
		return size{}, err
	}

	if width < 0 { // 負の場合、幅が制限されないときのサイズを調べる
		width = 10000000
	}

	lines := []string{}

	if run.Text != "" {
		if lines_, err := pdf.SplitText(run.Text, width); err != nil {
			return size{}, err
		} else {
			lines = lines_
		}
	}

	cs := size{
		h: run.FontSize * (float64(len(lines)) + (t.LineHeight-1)*float64(len(lines)-1)),
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
