package main

import (
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
	results := make(chan int)
	go Search(results, protein.Table, protein.Chain, posX, posY, 2, 0)

	fmt.Println(<-results)
}

func Search(results chan int, matrix [][]byte, chain string, posX int, posY int, k int, e int) {
	availableMoves := GetAvailableMoves(matrix, [2]int{posX, posY})

	// duplicating a table
	duplicate := make([][]byte, len(matrix))
	for i := range matrix {
		duplicate[i] = make([]byte, len(matrix[i]))
		copy(duplicate[i], matrix[i])
	}

	for _, move := range availableMoves {
		switch move {
		case 0: //up
			posX--
		case 1: //left
			posY--
		case 2: //right
			posY++
		case 3: //down
			posX++
		}
		energy, min, avg := CalculateEnergy(duplicate, [2]int{posX, posY}, move)
		duplicate[posX][posY] = chain[k]
		if k == len(chain)-1 {
			results <- e
			return
		} else {
			if chain[k] == 'h' {
				if energy < min {
					go Search(results, duplicate, chain, posX, posY, k+1, energy)
				}
				if float64(energy) > avg {
					go Search(results, duplicate, chain, posX, posY, k+1, energy)
				}
				if energy >= min && float64(energy) <= avg {
					go Search(results, duplicate, chain, posX, posY, k+1, energy)
				}
			} else {
				go Search(results, duplicate, chain, posX, posY, k+1, energy)
			}
		}
	}
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
	avg = float64(e) / float64(len(energies))
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
