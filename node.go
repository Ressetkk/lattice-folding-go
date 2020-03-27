package main

import "fmt"

type Node struct {
	x int
	h bool
}

func (n Node) String() string {
	var sh string
	if n.h {
		sh = "h"
	} else {
		sh = "p"
	}
	return fmt.Sprintf("(%v, %v) %v", (n.x&0xFFFF)-0x8000, ((n.x>>16)&0xFFFF)-0x8000, sh)
}

func NewNode(x int, h bool) *Node {
	return &Node{x, h}
}
