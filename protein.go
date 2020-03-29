package main

type Folding struct {
	parent  *Folding
	x, y, e int
	h       bool
}

func NewRootFolding(h bool) *Folding {
	return &Folding{parent: nil, x: x0, y: y0, h: h}
}

func NewFolding(x, y int, parent *Folding, h bool) *Folding {
	e := parent.e
	if !h {
		return &Folding{parent: parent, x: x, y: y, h: h, e: e}
	}
	sum := 1
	var temp *Folding
	temp = parent.parent
	L, R, U, D := x-ex, x+ex, y+ey, y-ey

	for temp != nil && sum < 4 {
		if (temp.x == L && temp.y == y) || (temp.x == R && temp.y == y) || (temp.y == U && temp.x == x) || (temp.y == D && temp.x == x) {
			if temp.h {
				e--
			}
			sum++
		}
		temp = temp.parent
	}
	return &Folding{parent: parent, x: x, y: y, h: h, e: e}
}

func (f *Folding) isEmpty(x, y int) bool {
	t := f
	for t != nil {
		if x == t.x && t.y == y {
			return false
		}
		t = t.parent
	}
	return true
}
