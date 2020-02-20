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
	arr := make([][]byte, size, '-')
	for i := range arr {
		arr[i] = make([]byte, size, '-')
	}
	return &Protein{
		Table:  arr,
		Chain:  strings.ToLower(chain),
		Result: 0,
	}, nil
}

func (protein Protein) String() string {
	stringer := ""
	for _, proteinRow := range protein.Table {
		stringer += fmt.Sprintf("%+v\n", proteinRow)
	}
	return stringer
}
