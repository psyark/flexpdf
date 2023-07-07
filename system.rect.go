package flexpdf

// rect はページ中の矩形を表します
type rect struct {
	x, y, w, h float64
}

// getSize は指定した軸の値を返します
func (s rect) getSize(a axis) float64 {
	if a == horizontal {
		return s.w
	} else {
		return s.h
	}
}

// setSize は指定した軸の値を設定し、新たなrectを返します
func (s rect) setSize(a axis, value float64) rect {
	if a == horizontal {
		s.w = value
	} else {
		s.h = value
	}
	return s
}

// updateSize は値を更新するためのコールバックを用いて指定した軸の値を更新し、新たなrectを返します
func (s rect) updateSize(a axis, fn func(float64) float64) rect {
	return s.setSize(a, fn(s.getSize(a)))
}

// getPos は指定した軸の値を返します
func (s rect) getPos(a axis) float64 {
	if a == horizontal {
		return s.x
	} else {
		return s.y
	}
}

// setPos は指定した軸の値を設定し、新たなrectを返します
func (s rect) setPos(a axis, value float64) rect {
	if a == horizontal {
		s.x = value
	} else {
		s.y = value
	}
	return s
}

// updatePos は値を更新するためのコールバックを用いて指定した軸の値を更新し、新たなrectを返します
func (s rect) updatePos(a axis, fn func(float64) float64) rect {
	return s.setPos(a, fn(s.getPos(a)))
}

func (s rect) getLength(a axis) float64 {
	if a == horizontal {
		return s.w
	} else {
		return s.h
	}
}
func (s rect) shrink(spacing Spacing) rect {
	s.x += spacing.Left
	s.w -= spacing.Left + spacing.Right
	s.y += spacing.Top
	s.h -= spacing.Top + spacing.Bottom
	// TODO negative
	return s
}
