package flexpdf

// GoでPDFを画像に変換する
// https://qiita.com/toshikitsubouchi/items/51c3268185cdc976a52f

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"gopkg.in/gographics/imagick.v3/imagick"
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
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(5),
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xFF, B: 0xCC, A: 0xFF}).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 5)),
			NewText("ipaexg", 30, "Text").SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xFF, A: 0xFF}).SetPadding(5),
		).SetMargin(30),
		NewBox(
			DirectionRow,
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF}),
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetAlign(TextAlignBegin),
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xFF, B: 0xCC, A: 0xFF}).SetAlign(TextAlignCenter),
			NewText("ipaexg", 30, "Text").SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xFF, A: 0xFF}).SetAlign(TextAlignEnd),
		).SetMargin(30).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
		NewBox(
			DirectionRow,
			NewText("ipaexg", 30, "あいうえおかきくけこさしすせそたちつてと").SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(10).SetBorder(UniformedBorder(color.Black, BorderStyleDotted, 1)),
		).SetMargin(30).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
		NewBox(
			DirectionRow,
			NewText("ipaexg", 30, "あいうえおかきくけこさしすせそたちつてと").SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(10).SetBorder(UniformedBorder(color.Black, BorderStyleDotted, 1)),
			NewText("ipaexg", 30, "あいうえおかきくけこさしすせそたちつてと").SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(10).SetBorder(UniformedBorder(color.Black, BorderStyleDotted, 1)),
		).SetMargin(30).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
	).SetPadding(50)

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

	if err := hoge(data); err != nil {
		t.Fatal(err)
	}
}

func hoge(pdfBytes []byte) error {
	// ImageMagic を初期化する。
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// 解像度を設定する。
	err := mw.SetResolution(150, 150)
	if err != nil {
		return errors.Wrap(err, "failed at SetResolution")
	}

	// 変換元のPDFを読み込む。
	err = mw.ReadImageBlob(pdfBytes)
	if err != nil {
		return errors.Wrap(err, "failed at ReadImage")
	}

	// ページ数を取得する。
	n := mw.GetNumberImages()
	log.Println("number image: ", n)

	// 出力フォーマットをPNGに設定する。
	err = mw.SetImageFormat("png")
	if err != nil {
		return errors.Wrap(err, "failed at SetImageFormat")
	}

	// １ページずつ変換して出力する。
	for i := 0; i < int(n); i++ {
		// ページ番号を設定する。
		if !mw.SetIteratorIndex(i) {
			break
		}

		// 画像を出力する。
		err = mw.WriteImage(fmt.Sprintf("test_%02d.png", i))
		if err != nil {
			return errors.Wrap(err, "failed at WriteImage")
		}
	}

	return nil
}

func TestXxx(t *testing.T) {
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

	if err := Draw(pdf, root, gopdf.PageSizeA4); err != nil {
		t.Fatal(err)
	}

	data, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

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
		20,
	).SetBorder(
		UniformedBorder(color.Black, BorderStyleDashed, 10),
	).SetPadding(
		20,
	).SetHeight(
		0,
	).SetFlexGrow(1).SetBackgroundColor(
		color.RGBA{0x00, 0x00, 0x00, 0x22},
	).SetJustifyContent(JustifyContentSpaceBetween)
}

func createJustifyContentExample(dir Direction, jc JustifyContent) *Box {
	return NewBox(
		DirectionColumn,
		NewText("ipaexg", 14, string(jc)),
		NewBox(
			dir,
			NewText("ipaexg", 15, "A").SetSize(20, 20).SetBackgroundColor(color.RGBA{0xFF, 0xCC, 0xCC, 0xFF}),
			NewText("ipaexg", 15, "B").SetSize(20, 20).SetBackgroundColor(color.RGBA{0xCC, 0xFF, 0xCC, 0xFF}),
			NewText("ipaexg", 15, "C").SetSize(20, 20).SetBackgroundColor(color.RGBA{0xCC, 0xCC, 0xFF, 0xFF}),
		).SetBackgroundColor(
			color.RGBA{0x88, 0x88, 0x88, 0xFF},
		).SetBorder(
			UniformedBorder(color.RGBA{A: 0xFF}, BorderStyleSolid, 2),
		).SetJustifyContent(
			jc,
		).SetFlexGrow(1),
	)
}
