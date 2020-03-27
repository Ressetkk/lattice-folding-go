package main

type Folding struct {
	parent *Folding
	pos, e int
	h      bool
}

func NewRootFolding(h bool) *Folding {
	return &Folding{parent: nil, pos: origin, h: h}
}

func NewFolding(pos int, parent *Folding, h bool) *Folding {
	if !h {
		return &Folding{parent: parent, pos: pos, h: h}
	}
	e := parent.e
	sum := 1
	var temp *Folding
	if parent != nil {
		temp = parent.parent
	}
	L, R, U, D := pos-ex, pos+ex, pos+ey, pos-ey

	for temp != nil && sum < 4 {
		if temp.pos == L || temp.pos == R || temp.pos == U || temp.pos == D {
			if temp.h {
				e--
			}
			sum++
		}
		temp = temp.parent
	}
	return &Folding{parent: parent, pos: pos, h: h, e: e}
}

func (f *Folding) isEmpty(pos int) bool {
	t := f
	for t != nil {
		if pos == t.pos {
			return false
		}
		t = t.parent
	}
	return true
}
