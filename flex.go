package flexpdf

type Direction int

const (
	Row Direction = iota
	Column
)

type JustifyContent int

const (
	JustifyContentFlexStart JustifyContent = iota
	JustifyContentFlexEnd
	JustifyContentCenter
	JustifyContentSpaceBetween
	JustifyContentSpaceAround
	JustifyContentSpaceEvenly
)

type AlignItems int

const (
	AlignItemsFlexStart AlignItems = iota
	AlignItemsFlexEnd
	AlignItemsCenter
)
