package flexpdf

import "image/color"

type flexItemCommon struct {
	Width           float64
	Height          float64
	FlexGrow        float64
	FlexShrink      float64
	BackgroundColor color.Color
	Border          Border
}
