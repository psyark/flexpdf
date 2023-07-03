package flexpdf

type Direction int

const (
	DirectionRow Direction = iota
	DirectionColumn
)

func (d Direction) String() string {
	switch d {
	case DirectionRow:
		return "row"
	case DirectionColumn:
		return "column"
	default:
		panic(d)
	}
}

func (d Direction) mainAxis() axis {
	switch d {
	case DirectionRow:
		return horizontal
	case DirectionColumn:
		return vertical
	default:
		panic(d)
	}
}

// https://www.w3.org/TR/css-flexbox/#justify-content-property
type JustifyContent string

const (
	JustifyContentFlexStart    JustifyContent = "flex-start"
	JustifyContentFlexEnd      JustifyContent = "flex-end"
	JustifyContentCenter       JustifyContent = "center"
	JustifyContentSpaceBetween JustifyContent = "space-between"
	JustifyContentSpaceAround  JustifyContent = "space-around"
	// JustifyContentSpaceEvenly  JustifyContent = "space-evenly"
)

// https://www.w3.org/TR/css-flexbox/#propdef-align-items
type AlignItems string

const (
	AlignItemsFlexStart AlignItems = "flex-start"
	AlignItemsFlexEnd   AlignItems = "flex-end"
	AlignItemsCenter    AlignItems = "center"
	AlignItemsStretch   AlignItems = "stretch"
)
