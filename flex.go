package flexpdf

type Direction string

const (
	DirectionRow    Direction = "row"
	DirectionColumn Direction = "column"
)

type JustifyContent string

const (
	JustifyContentFlexStart    JustifyContent = "flex-start"
	JustifyContentFlexEnd      JustifyContent = "flex-end"
	JustifyContentCenter       JustifyContent = "center"
	JustifyContentSpaceBetween JustifyContent = "space-between"
	JustifyContentSpaceAround  JustifyContent = "space-around"
	JustifyContentSpaceEvenly  JustifyContent = "space-evenly"
)

type AlignItems string

const (
	AlignItemsFlexStart AlignItems = "flex-start"
	AlignItemsFlexEnd   AlignItems = "flex-end"
	AlignItemsCenter    AlignItems = "center"
	AlignItemsStretch   AlignItems = "stretch"
)
