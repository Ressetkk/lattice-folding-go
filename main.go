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

const x0, y0, ex, ey = 0, 0, 1, 1

var (
	chain  = flag.String("protein", "", "character stream over the alphabet of {h, p}")
	p1     = flag.Float64("p1", 0.6, "first probability [Default 0.6]")
	p2     = flag.Float64("p2", 0.4, "second probability [Default 0.4]")
	output = flag.String("output", "out.png", "name of the generated image [Default 'out.png']")
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
	if err := GenerateImage(folding, *output); err != nil {
		fmt.Printf("Could not generate image: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Successfully created image %v", *output)
	}
}

func Fold(ctx context.Context, chain string) *[]*Node {
	check := func(chain byte) bool {
		return chain == 'h'
	}

	switch len(chain) {
	case 1:
		return &[]*Node{NewNode(x0, y0, check(chain[0]))}
	case 2:
		return &[]*Node{NewNode(x0, y0, check(chain[0])), NewNode(x0+ex, y0, check(chain[1]))}
	}
	min, k, minIndex := 0, 2, 0
	var avg float64
	branches := &[]*Folding{NewFolding(x0+ex, y0, NewRootFolding(check(chain[0])), check(chain[1]))}
	for k < len(chain) {
		ifH := check(chain[k])
		var temp []*Folding
		var c, localMin, localSum int
		for _, b := range *branches {
			r := rand.Float64()
			foldIfFree := func(x, y int) *Folding {
				if b.isEmpty(x, y) {
					return NewFolding(x, y, b, ifH)
				}
				return nil
			}

			L, R, U, D := b.x-ex, b.x+ex, b.y+ey, b.y-ey
			next := &[]*Folding{foldIfFree(L, b.y), foldIfFree(R, b.y), foldIfFree(b.x, U), foldIfFree(b.x, D)} // TODO make it parallel
			for _, nb := range *next {
				if nb != nil && (nb.e <= min || r >= *p1 || (float64(nb.e) <= avg && r >= *p2)) {
					if localMin > nb.e {
						localMin = nb.e
						minIndex = c
					}
					localSum += nb.e
					temp = append(temp, nb)
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
		branches = &temp
		k++
	}
	if k < len(chain) {
		return &[]*Node{}
	}
	ret := make([]*Node, len(chain))
	temp := (*branches)[minIndex]
	i := k - 1
	for i >= 0 {
		ret[i] = NewNode(temp.x, temp.y, temp.h)
		temp = temp.parent
		i--
	}
	println(min)
	return &ret
}
