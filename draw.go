package flexpdf

import (
	"github.com/signintech/gopdf"
)

func Draw(pdf *gopdf.GoPdf, box *Box, pageSize *gopdf.Rect) error {
	pdf.AddPageWithOption(gopdf.PageOption{PageSize: pageSize})

	if err := box.draw(pdf, rect{0, 0, pageSize.W, pageSize.H}); err != nil {
		return err
	}

	return nil
}
