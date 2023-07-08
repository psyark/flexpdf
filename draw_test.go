package flexpdf

// GoでPDFを画像に変換する
// https://qiita.com/toshikitsubouchi/items/51c3268185cdc976a52f

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/signintech/gopdf"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	//go:embed "testdata/fonts/ipaexg.ttf"
	ipaexgBytes []byte
	//go:embed "testdata/fonts/ipaexm.ttf"
	ipaexmBytes []byte
)

var errUnmatch = errors.New("unmatch")

func TestMain(m *testing.M) {
	code := (func() int {
		// 同じスコープで os.Exit すると 待機中のdefer が呼ばれないため
		imagick.Initialize()
		defer imagick.Terminate()

		return m.Run()
	})()

	os.Exit(code)
}

func TestDraw(t *testing.T) {
	for name, box := range cases {
		name, box := name, box
		t.Run(name, func(t *testing.T) {
			pdf := &gopdf.GoPdf{}
			pdf.Start(gopdf.Config{})

			if err := pdf.AddTTFFontData("ipaexg", ipaexgBytes); err != nil {
				t.Fatal(err)
			}
			if err := pdf.AddTTFFontData("ipaexm", ipaexmBytes); err != nil {
				t.Fatal(err)
			}
			if err := pdf.AddTTFFontData("", ipaexgBytes); err != nil {
				t.Fatal(err)
			}

			if err := Draw(pdf, box, gopdf.PageSizeA4); err != nil {
				t.Fatal(err)
			}

			data, err := pdf.GetBytesPdfReturnErr()
			if err != nil {
				t.Fatal(err)
			}

			if err := os.WriteFile(fmt.Sprintf("testdata/out/%s.pdf", name), data, 0666); err != nil {
				t.Fatal(err)
			}

			if err := compareImage(data, fmt.Sprintf("testdata/out/%s.png", name)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

var cases = map[string]*Box{
	"text": NewColumnBox(
		// 1 run for 1 text
		NewRowBox(
			NewText(NewRun("normal").SetFontSize(30)),
			NewText(NewRun("size").SetFontSize(60)),
			NewText(NewRun("color").SetFontSize(30).SetColor(color.RGBA{R: 0xFF, A: 0xFF})),
			NewText(NewRun("family").SetFontSize(30).SetFontFamily("ipaexm")),
		),
		// many run for 1 text
		NewRowBox(
			NewText(
				NewRun("normal").SetFontSize(30),
				NewRun("size").SetFontSize(60),
				NewRun("color").SetFontSize(30).SetColor(color.RGBA{R: 0xFF, A: 0xFF}),
				NewRun("family").SetFontSize(30).SetFontFamily("ipaexm"),
			),
		),
		NewRowBox(
			NewText(NewRun("Text").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF}),
			NewText(NewRun("Text").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(5),
			NewText(NewRun("Text").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xFF, B: 0xCC, A: 0xFF}).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 5)),
			NewText(NewRun("Text").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xFF, A: 0xFF}).SetPadding(5),
		).SetMargin(30),
		NewRowBox(
			NewText(NewRun("Text").SetFontSize(30)).SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF}),
			NewText(NewRun("Text").SetFontSize(30)).SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetAlign(TextAlignBegin),
			NewText(NewRun("Text").SetFontSize(30)).SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xFF, B: 0xCC, A: 0xFF}).SetAlign(TextAlignCenter),
			NewText(NewRun("Text").SetFontSize(30)).SetFlexGrow(1).SetBackgroundColor(color.RGBA{R: 0xCC, G: 0xCC, B: 0xFF, A: 0xFF}).SetAlign(TextAlignEnd),
		).SetMargin(30).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
		NewRowBox(
			NewText(NewRun("あいうえおかきくけこさしすせそたちつてと").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(10).SetBorder(UniformedBorder(color.Black, BorderStyleDotted, 1)),
		).SetMargin(30).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
		NewRowBox(
			NewText(NewRun("あいうえおかきくけこさしすせそたちつてと").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(10).SetBorder(UniformedBorder(color.Black, BorderStyleDotted, 1)),
			NewText(NewRun("あいうえおかきくけこさしすせそたちつてと").SetFontSize(30)).SetBackgroundColor(color.RGBA{R: 0xFF, G: 0xCC, B: 0xCC, A: 0xFF}).SetMargin(10).SetBorder(UniformedBorder(color.Black, BorderStyleDotted, 1)),
		).SetMargin(30).SetBorder(UniformedBorder(color.Black, BorderStyleDashed, 1)),
	).SetPadding(50),

	"justifycontent": NewColumnBox(
		createJustifyContentExamples(DirectionColumn, DirectionRow),
		createJustifyContentExamples(DirectionRow, DirectionColumn),
	),
}

func compareImage(pdfBytes []byte, fileName string) (err error) {
	defer wrap(&err, "compareImage")

	images, err := getImages(pdfBytes)
	if err != nil {
		return err
	}

	if len(images) != 1 {
		return fmt.Errorf("len(images) = %d", len(images))
	}

	imgGot := images[0]
	var bytesGot []byte

	{
		buf := bytes.NewBuffer(nil)
		if err := png.Encode(buf, imgGot); err != nil {
			return err
		}
		bytesGot = buf.Bytes()
	}

	bytesWant, err := os.ReadFile(fileName)

	if errors.Is(err, os.ErrNotExist) {
		// 存在しない場合は保存する。比較はしない
		return os.WriteFile(fileName, bytesGot, 0666)
	} else if err != nil {
		return err
	}

	diffFileName := strings.TrimSuffix(fileName, ".png") + "_diff.png"
	gotFileName := strings.TrimSuffix(fileName, ".png") + "_got.png"
	if bytes.Equal(bytesGot, bytesWant) {
		_ = os.Remove(diffFileName)
		_ = os.Remove(gotFileName)
		return nil
	} else {
		imgWant, err := png.Decode(bytes.NewReader(bytesWant))
		if err != nil {
			return err
		}

		imgDiff := image.NewGray(imgWant.Bounds().Union(imgGot.Bounds()))
		for y := imgDiff.Rect.Min.Y; y < imgDiff.Rect.Max.Y; y++ {
			for x := imgDiff.Rect.Min.X; x < imgDiff.Rect.Max.X; x++ {
				if imgWant.At(x, y) == imgGot.At(x, y) {
					imgDiff.Set(x, y, color.White)
				} else {
					imgDiff.Set(x, y, color.Black)
				}
			}
		}

		buf := bytes.NewBuffer(nil)
		if err := png.Encode(buf, imgDiff); err != nil {
			return err
		}

		if err := os.WriteFile(diffFileName, buf.Bytes(), 0666); err != nil {
			return err
		}
		if err := os.WriteFile(gotFileName, bytesGot, 0666); err != nil {
			return err
		}

		return fmt.Errorf("%s: %w", fileName, errUnmatch)
	}
}

func getImages(pdfBytes []byte) (images []image.Image, err error) {
	defer wrap(&err, "getImages")

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(150, 150); err != nil {
		return nil, err
	}
	if err := mw.ReadImageBlob(pdfBytes); err != nil {
		return nil, err
	}
	if err := mw.SetImageFormat("png"); err != nil {
		return nil, err
	}

	for i := 0; i < int(mw.GetNumberImages()); i++ {
		if !mw.SetIteratorIndex(i) {
			break
		}

		imageBytes := mw.GetImageBlob()
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return images, nil
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
	return NewColumnBox(
		NewText(NewRun(string(jc)).SetFontSize(14)),
		NewBox(
			dir,
			NewText(NewRun("A").SetFontSize(15)).SetSize(20, 20).SetBackgroundColor(color.RGBA{0xFF, 0xCC, 0xCC, 0xFF}),
			NewText(NewRun("B").SetFontSize(15)).SetSize(20, 20).SetBackgroundColor(color.RGBA{0xCC, 0xFF, 0xCC, 0xFF}),
			NewText(NewRun("C").SetFontSize(15)).SetSize(20, 20).SetBackgroundColor(color.RGBA{0xCC, 0xCC, 0xFF, 0xFF}),
		).SetBackgroundColor(
			color.RGBA{0x88, 0x88, 0x88, 0xFF},
		).SetBorder(
			UniformedBorder(color.RGBA{A: 0xFF}, BorderStyleSolid, 2),
		).SetJustifyContent(
			jc,
		).SetFlexGrow(1),
	)
}
