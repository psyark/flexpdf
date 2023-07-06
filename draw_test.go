package flexpdf

import (
	_ "embed"
	"image/color"
	"os"
	"testing"

	"github.com/signintech/gopdf"
)

var (
	//go:embed "testdata/fonts/ipaexg.ttf"
	ipaexgBytes []byte
	//go:embed "testdata/fonts/ipaexm.ttf"
	ipaexmBytes []byte
)

func TestText(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})

	if err := pdf.AddTTFFontData("ipaexg", ipaexgBytes); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFontData("", ipaexgBytes); err != nil {
		t.Fatal(err)
	}

	root := NewBox(
		DirectionColumn,
		NewBox(
			DirectionRow,
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF}),
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(UniformedSpacing(5)),
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xFF, B: 0xCC, A: 0xFF}).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 5)),
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xFF, A: 0xFF}).SetPadding(UniformedSpacing(5)),
		).SetMargin(UniformedSpacing(30)),
		NewBox(
			DirectionRow,
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF}),
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetAlign(TextAlignBegin),
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xFF, B: 0xCC, A: 0xFF}).SetAlign(TextAlignCenter),
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xFF, A: 0xFF}).SetAlign(TextAlignEnd),
		).SetMargin(UniformedSpacing(30)).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
	).SetPadding(UniformedSpacing(50))

	t.Log("draw start")
	if err := Draw(pdf, root, gopdf.PageSizeA4); err != nil {
		t.Fatal(err)
	}

	t.Log("draw end")
	data, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("write start")
	if err := os.WriteFile("text.pdf", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func TestXxx(t *testing.T) {
	t.Skip()

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})

	if err := pdf.AddTTFFontData("ipaexg", ipaexgBytes); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFontData("ipaexm", ipaexmBytes); err != nil {
		t.Fatal(err)
	}

	root := NewBox(
		DirectionColumn,
		createJustifyContentExamples(DirectionColumn, DirectionRow),
		createJustifyContentExamples(DirectionRow, DirectionColumn),
	)

	t.Log("draw start")
	if err := Draw(pdf, root, gopdf.PageSizeA4); err != nil {
		t.Fatal(err)
	}

	t.Log("draw end")
	data, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("write start")
	if err := os.WriteFile("out.pdf", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func createJustifyContentExamples(dir1, dir2 Direction) *Box {
	return NewBox(
		dir1,
		createJustifyContentExample(dir2, JustifyContentFlexStart),
		createJustifyContentExample(dir2, JustifyContentFlexEnd),
		createJustifyContentExample(dir2, JustifyContentCenter),
		createJustifyContentExample(dir2, JustifyContentSpaceBetween),
		createJustifyContentExample(dir2, JustifyContentSpaceAround),
	).SetMargin(
		UniformedSpacing(20),
	).SetBorder(
		UniformedBorder(color.Black, BorderStyleDashed, 10),
	).SetPadding(
		UniformedSpacing(20),
	).SetHeight(
		300,
	).SetBackgroundColor(
		color.RGBA{0x00, 0x00, 0x00, 0x22},
	)
}

func createJustifyContentExample(dir Direction, jc JustifyContent) *Box {
	return NewBox(
		DirectionColumn,
		NewText("ipaexg", 20, string(jc)+":"),
		NewBox(
			dir,
			NewText("ipaexg", 24, "あいうえお").SetSize(80, 40).SetBackgroundColor(color.RGBA{0xFF, 0xCC, 0xCC, 0xFF}),
			NewText("ipaexg", 24, "かきくけこ").SetBackgroundColor(color.RGBA{0xCC, 0xFF, 0xCC, 0xFF}),
			NewText("ipaexg", 24, "さしすせそ たちつてと").SetBackgroundColor(color.RGBA{0xCC, 0xCC, 0xFF, 0xFF}),
		).SetBackgroundColor(
			color.RGBA{0x88, 0x88, 0x88, 0xFF},
		).SetBorder(
			UniformedBorder(color.RGBA{A: 0xFF}, BorderStyleSolid, 2),
		).SetJustifyContent(
			jc,
		),
	)
}
