package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const origin, ex, ey = (1 << 31) + (1 << 15), 1, 1 << 16

var (
	chain = flag.String("protein", "", "character stream over the alphabet of {h, p}")
	p1    = flag.Float64("-p1", 0.6, "first probability [Dafault 0.4]")
	p2    = flag.Float64("-p2", 0.4, "second probability [Dafault 0.2]")
)

func main() {
	flag.Parse()
	*chain = strings.ToLower(*chain)
	if check, _ := regexp.MatchString("[ph]*", *chain); !check {
		fmt.Println("Chain does not match - must be Pp/Hh")
		os.Exit(1)
	}

	if len(*chain) == 0 {
		fmt.Println("Chain is empty. Provide valid chain.")
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	folding := Fold(ctx, *chain)
	fmt.Println(folding)
}

func Fold(ctx context.Context, chain string) *[]*Node {
	check := func(chain byte) bool {
		return chain == 'h'
	}

	switch len(chain) {
	case 1:
		return &[]*Node{NewNode(origin, check(chain[0]))}
	case 2:
		return &[]*Node{NewNode(origin, check(chain[0])), NewNode(origin+ex, check(chain[1]))}
	}
	min, k, minIndex := 0, 2, 0
	var avg float64
	branches := &[]*Folding{NewFolding(origin+ex, NewRootFolding(check(chain[0])), check(chain[1]))}
	for k < len(chain) {
		ifH := check(chain[k])
		temp := new([]*Folding)
		var c, localMin, localSum int
		for _, b := range *branches {
			r := rand.Float64()
			foldIfFree := func(pos int) *Folding {
				if b.isEmpty(pos) {
					return NewFolding(pos, b, ifH)
				}
				return nil
			}

			L, R, U, D := b.pos-ex, b.pos+ex, b.pos+ey, b.pos-ey
			next := &[]*Folding{foldIfFree(L), foldIfFree(R), foldIfFree(U), foldIfFree(D)} // TODO make it parallel
			for _, nb := range *next {
				if nb != nil && (nb.e <= min || r >= *p1 || (float64(nb.e) <= avg && r >= *p2)) { // TODO maybe wrong logic here
					if localMin > nb.e {
						localMin = nb.e
						minIndex = c
					}
					localSum += nb.e
					*temp = append(*temp, nb)
					c++
				}
			}
		}
		if c == 0 {
			fmt.Println("All paths exhausted...")
			return &[]*Node{}
		}
		min = localMin
		avg = float64(localSum) / float64(c)
		branches = temp
		k++
	}
	if k < len(chain) {
		return &[]*Node{}
	}
	ret := make([]*Node, len(chain))
	temp := (*branches)[minIndex]
	i := k - 1
	for i >= 0 {
		ret[i] = NewNode(temp.pos, temp.h)
		temp = temp.parent
		i--
	}
	println(min)
	return &ret
}
