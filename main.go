package main

import (
	"bytes"
	"flag"
	"fmt"
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
	_ = Search(*protein, []int{posX, posY}, 1, 0)
}

func Search(protein Protein, pos []int, k int, partialEnergy int) int {
	positions := GetAvailableMoves(protein.Table, pos)
	for _, v := range protein.Table {
		fmt.Println(v)
	}
	fmt.Println()
	for i := range positions {
		posX := pos[0]
		posY := pos[1]
		pt := protein.Table
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
		e, min, avg := CalculateEnergy(pt, []int{posX, posY})
		partialNextEnergy := partialEnergy + e
		if k == len(protein.Table)-1 {
			return partialNextEnergy
		} else {
			fmt.Println(protein.Chain[k] == 'h', protein.Chain[k], 'h')
			if protein.Chain[k] == 'h' {
				if partialNextEnergy <= min {
					protein.Table[posX][posY] = protein.Chain[k]
					return Search(protein, []int{posX, posY}, k+1, partialNextEnergy)
				}
				if partialNextEnergy > avg {
					if float64(partialNextEnergy) > p1 {
						protein.Table[posX][posY] = protein.Chain[k]
						return Search(protein, []int{posX, posY}, k+1, partialNextEnergy)
					}
				}
				if min <= partialNextEnergy && partialNextEnergy <= avg {
					if float64(partialNextEnergy) > p2 {
						protein.Table[posX][posY] = protein.Chain[k]
						return Search(protein, []int{posX, posY}, k+1, partialNextEnergy)
					}
				}
			} else {
				protein.Table[posX][posY] = protein.Chain[k]
				return Search(protein, []int{posX, posY}, k+1, partialNextEnergy)
			}
		}
	}
	return 1
}
func CalculateEnergy(proteinTable [][]byte, pos []int) (e int, min int, avg int) {
	min = 1
	if pos[0]-1 > 0 {
		if proteinTable[pos[0]-1][pos[1]] != 0 && bytes.Equal([]byte{proteinTable[pos[0]-1][pos[1]], proteinTable[pos[0]][pos[1]]}, []byte("hh")) {
			e += 1
		} else {
			min = 0
		}
	}

	if pos[0]+1 < len(proteinTable) {
		if proteinTable[pos[0]+1][pos[1]] != 0 && bytes.Equal([]byte{proteinTable[pos[0]+1][pos[1]], proteinTable[pos[0]][pos[1]]}, []byte("hh")) {
			e += 1
		} else {
			min = 0
		}
	}

	if pos[1]-1 > 0 {
		if proteinTable[pos[0]][pos[1]-1] != 0 && bytes.Equal([]byte{proteinTable[pos[0]][pos[1]-1], proteinTable[pos[0]][pos[1]]}, []byte("hh")) {
			e += 1
		} else {
			min = 0
		}
	}

	if pos[1]+1 < len(proteinTable) {
		if proteinTable[pos[0]][pos[1]+1] != 0 && bytes.Equal([]byte{proteinTable[pos[0]][pos[1]+1], proteinTable[pos[0]][pos[1]]}, []byte("hh")) {
			e += 1
		} else {
			min = 0
		}
	}
	avg = e / 4.
	return
}

func GetAvailableMoves(proteinTable [][]byte, pos []int) (moves []int) {
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
