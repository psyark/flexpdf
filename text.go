package flexpdf

type Text struct {
	Text string
}

func (*Text) flexItem() {}

type TextSpan struct {
}
