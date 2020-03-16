package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Protein struct {
	Table  [][]byte
	Chain  string
	Result int
	Move   int
	h      bool
	parent *Protein
}

var ChainNotMatchError = errors.New("chain does not match - must be Pp/Hh")

func New(chain string) (*Protein, error) {
	if check, _ := regexp.MatchString("[ph]*", strings.ToLower(chain)); !check {
		return nil, ChainNotMatchError
	}
	size := len(chain)
	if size%2 == 0 {
		size++
	}
	arr := make([][]byte, size)
	for i := range arr {
		arr[i] = make([]byte, size)
	}
	return &Protein{
		Table:  arr,
		Chain:  strings.ToLower(chain),
		Result: 0,
		h:      chain[0] == 'h',
	}, nil
}

func NewChild(move int, parent *Protein, result int, table [][]byte, h bool) *Protein {
	return &Protein{Table: table, Chain: parent.Chain, Result: result, Move: move, parent: parent, h: h}
}

func (protein *Protein) String() string {
	var sp, h string
	switch protein.Move {
	case 0:
		sp = "LEFT"
		break
	case 1:
		sp = "RIGHT"
		break
	case 2:
		sp = "UP"
		break
	case 3:
		sp = "DOWN"
		break
	case -1:
		sp = ""
	}
	if protein.h {
		h = "h"
	} else {
		h = "p"
	}
	return fmt.Sprintf("%s %v", h, sp)
}
