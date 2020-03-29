package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

const x0, y0, ex, ey = 0, 0, 1, 1
const p1, p2 = 0.7, 0.5

var (
	chain  = flag.String("protein", "", "character stream over the alphabet of {h, p}")
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
	ret, e := Fold(*chain)
	if err := GenerateImage(ret, e, *output); err != nil {
		fmt.Printf("An error occurred during image generation: %v\n", err)
	} else {
		fmt.Printf("Successfully generated image %s\n", *output)
	}
}

func Fold(chain string) (*[]*Node, int) {
	check := func(chain byte) bool {
		return chain == 'h'
	}

	switch len(chain) {
	case 1:
		return &[]*Node{NewNode(x0, y0, check(chain[0]))}, 0
	case 2:
		return &[]*Node{NewNode(x0, y0, check(chain[0])), NewNode(x0+ex, y0, check(chain[1]))}, 0
	}
	min, k, minIndex := 0, 2, 0
	var avg float64
	branches := &[]*Folding{NewFolding(x0+ex, y0, NewRootFolding(check(chain[0])), check(chain[1]))}
	for k < len(chain) {
		ifH := check(chain[k])
		var temp []*Folding
		var c, localSum int
		for _, b := range *branches {
			r := rand.Float64()
			foldIfFree := func(x, y int) *Folding {
				if b.isEmpty(x, y) {
					return NewFolding(x, y, b, ifH)
				}
				return nil
			}
			L, R, U, D := b.x-ex, b.x+ex, b.y+ey, b.y-ey
			res := make(chan *Folding, 4)
			go func(x, y int) {
				res <- foldIfFree(x, y)
			}(L, b.y)
			go func(x, y int) {
				res <- foldIfFree(x, y)
			}(R, b.y)
			go func(x, y int) {
				res <- foldIfFree(x, y)
			}(b.x, U)
			go func(x, y int) {
				res <- foldIfFree(x, y)
			}(b.x, D)

			for i := 0; i < 4; i++ {
				select {
				case nb, ok := <-res:
					if ok && nb != nil && (nb.e <= min || r >= p1 || (float64(nb.e) <= avg && r >= p2)) {
						if min > nb.e {
							min = nb.e
							c = 0
							minIndex = c
							temp = []*Folding{}
						}
						localSum += nb.e
						temp = append(temp, nb)
						c++
					}
				}
			}
		}
		if c == 0 {
			fmt.Println("All paths exhausted...")
			return &[]*Node{}, 0
		}
		avg = float64(localSum) / float64(c)
		branches = &temp
		k++
	}
	if k < len(chain) {
		return &[]*Node{}, 0
	}
	ret := make([]*Node, len(chain))
	temp := (*branches)[minIndex]
	i := k - 1
	for i >= 0 {
		ret[i] = NewNode(temp.x, temp.y, temp.h)
		temp = temp.parent
		i--
	}
	return &ret, min
}
