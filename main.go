package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var (
	proteinChain = flag.String("protein", "", "character stream over the alphabet of {h, p}")
	threads      = flag.Int("threads", 4, "maximum concurrent threads (default: 4)")
	p1           = 0.5
	p2           = 0.8
)

func main() {
	flag.Parse()
	protein, err := New(*proteinChain)
	if err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	posX := rand.Intn(len(protein.Table))
	posY := rand.Intn(len(protein.Table))
	protein.Table[posX][posY] = []byte(*proteinChain)[0]
	for {
		done := false
		switch rand.Intn(4) {
		case 0: //up
			if posX-1 < 0 {
				continue
			}
			posX -= 1
			done = true
			break
		case 1: //left
			if posY-1 < 0 {
				continue
			}
			posY -= 1
			done = true
			break
		case 2: //right
			if posY+1 > len(protein.Table)-1 {
				continue
			}
			posY += 1
			done = true
			break
		case 3: //down
			if posX+1 > len(protein.Table)-1 {
				continue
			}
			posX += 1
			done = true
			break
		}
		if done {
			protein.Table[posX][posY] = []byte(*proteinChain)[1]
			break
		}
	}
	e := Search(*protein, [2]int{posX, posY}, 2, 0)
	fmt.Println(e)
}

func Search(protein Protein, pos [2]int, k int, partialEnergy int) int {
	positions := GetAvailableMoves(protein.Table, pos)
	var results []int
	for _, i := range positions {
		var res int
		posX := pos[0]
		posY := pos[1]
		pt := append(protein.Table[:0:0], protein.Table...)
		switch i {
		case 0: //up
			posX--
			break
		case 1: //left
			posY--
			break
		case 2: //right
			posX++
			break
		case 3: //down
			posY++
			break
		}

		pt[posX][posY] = protein.Chain[k]
		proteinCopy := protein
		proteinCopy.Table = pt

		e, min, avg := CalculateEnergy(pt, [2]int{posX, posY}, i)
		fmt.Println(e, min, avg)
		partialNextEnergy := partialEnergy + e
		if k == len(protein.Table)-1 {
			return partialNextEnergy
		} else {
			if protein.Chain[k] == 'h' {
				if partialNextEnergy <= min {
					res = Search(proteinCopy, [2]int{posX, posY}, k+1, partialNextEnergy)
					results = append(results, res)
				}
				if float64(partialNextEnergy) > avg {
					if float64(partialNextEnergy) > p1 {
						res = Search(proteinCopy, [2]int{posX, posY}, k+1, partialNextEnergy)
						results = append(results, res)
					}
				}
				if min <= partialNextEnergy && float64(partialNextEnergy) <= avg {
					if float64(partialNextEnergy) > p2 {
						res = Search(proteinCopy, [2]int{posX, posY}, k+1, partialNextEnergy)
						results = append(results, res)
					}
				}
			} else {
				res = Search(proteinCopy, [2]int{posX, posY}, k+1, partialNextEnergy)
				results = append(results, res)
			}
		}
		//return res
	}
	min := 0
	for _, v := range results {
		if v < min {
			min = v
		}
	}
	return min
}
func CalculateEnergy(proteinTable [][]byte, pos [2]int, from int) (e int, min int, avg float64) {
	var moves [3]int
	switch from {
	case 0:
		moves = [3]int{0, 1, 2}
		break
	case 1:
		moves = [3]int{0, 1, 3}
		break
	case 2:
		moves = [3]int{0, 2, 3}
		break
	case 3:
		moves = [3]int{1, 2, 3}
		break
	}
	check := func() []int {
		var e []int
		for _, move := range moves {
			var next [2]int
			switch move {
			case 0:
				next = [2]int{pos[0] - 1, pos[1]}
				break
			case 1:
				next = [2]int{pos[0], pos[1] - 1}
				break
			case 3:
				next = [2]int{pos[0] + 1, pos[1]}
				break
			case 2:
				next = [2]int{pos[0], pos[1] + 1}
				break
			}
			if next[0] >= 0 && next[0] < len(proteinTable) && next[1] >= 0 && next[1] < len(proteinTable) {
				if proteinTable[next[0]][next[1]] == proteinTable[pos[0]][pos[1]] {
					e = append(e, -1)
				} else {
					e = append(e, 0)
				}
			}
		}
		return e
	}
	energies := check()
	for _, v := range energies {
		e += v
		if v < min {
			min = v
		}
	}
	avg = math.Abs(float64(e)) / float64(len(energies))
	return
}

func GetAvailableMoves(proteinTable [][]byte, pos [2]int) (moves []int) {
	for i := 0; i < 4; i++ {
		switch i {
		case 0: //up
			if pos[0]-1 >= 0 {
				if proteinTable[pos[0]-1][pos[1]] == 0 {
					moves = append(moves, i)
				}
			}
			break
		case 1: //left
			if pos[1]-1 >= 0 {
				if proteinTable[pos[0]][pos[1]-1] == 0 {
					moves = append(moves, i)
				}
			}
			break
		case 2: //right
			if pos[1]+1 < len(proteinTable) {
				if proteinTable[pos[0]][pos[1]+1] == 0 {
					moves = append(moves, i)
				}
			}
			break
		case 3: //down
			if pos[0]+1 < len(proteinTable) {
				if proteinTable[pos[0]+1][pos[1]] == 0 {
					moves = append(moves, i)
				}
			}
			break
		}
	}
	return
}
