package flexpdf

import "github.com/signintech/gopdf"

type FlexItem interface {
	draw(*gopdf.GoPdf, rect) error
}
