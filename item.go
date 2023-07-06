package flexpdf

import "image/color"

var flexItemCommonDefault = flexItemCommon{
	Width:           -1,
	Height:          -1,
	FlexGrow:        0,
	FlexShrink:      1,
	BackgroundColor: nil,
	Border:          UniformedBorder(nil, BorderStyleSolid, 0), // TODO None
}

type flexItemCommon struct {
	Width           float64
	Height          float64
	FlexGrow        float64
	FlexShrink      float64
	BackgroundColor color.Color
	Border          Border
}
