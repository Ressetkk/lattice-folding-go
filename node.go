package main

import (
	"fmt"
)

type Node struct {
	x, y int
	h    bool
}

func (n Node) String() string {
	var sh string
	if n.h {
		sh = "h"
	} else {
		sh = "p"
	}
	return fmt.Sprintf("(%v, %v) %v", n.x, n.y, sh)
}

func NewNode(x, y int, h bool) *Node {
	return &Node{x, y, h}
}
