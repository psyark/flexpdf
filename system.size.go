package flexpdf

// sizeは幅と高さを表現する値です
type size struct {
	w float64
	h float64
}

// get は指定した軸の値を返します
func (s size) get(a axis) float64 {
	if a == horizontal {
		return s.w
	} else {
		return s.h
	}
}

// set は指定した軸の値を設定し、新たなsizeを返します
func (s size) set(a axis, value float64) size {
	if a == horizontal {
		s.w = value
	} else {
		s.h = value
	}
	return s
}

// add は指定した軸の値を加算し、新たなsizeを返します
func (s size) add(a axis, value float64) size {
	if a == horizontal {
		s.w += value
	} else {
		s.h += value
	}
	return s
}

// update は値を更新するためのコールバックを用いて指定した軸の値を更新し、新たなsizeを返します
func (s size) update(a axis, fn func(float64) float64) size {
	return s.set(a, fn(s.get(a)))
}

func (s size) expand(spacing Spacing) size {
	s.w += spacing.Left + spacing.Right
	s.h += spacing.Top + spacing.Bottom
	// TODO negative
	return s
}
func (s size) shrink(spacing Spacing) size {
	s.w -= spacing.Left + spacing.Right
	s.h -= spacing.Top + spacing.Bottom
	// TODO negative
	return s
}
