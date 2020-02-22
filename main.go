package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"
)

var (
	proteinChain = flag.String("protein", "", "character stream over the alphabet of {h, p}")
	p1           = 0.2
	p2           = 0.3
)

func main() {
	flag.Parse()
	protein, err := New(*proteinChain)
	if err != nil {
		log.Fatalf("An error occured: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	posX := rand.Intn(len(protein.Table))
	posY := rand.Intn(len(protein.Table))
	protein.Table[posX][posY] = []byte(*proteinChain)[0]

	for {
		nPosX, nPosY, err := Move(rand.Intn(4), posX, posY, len(protein.Table))
		if err == nil {
			posX = nPosX
			posY = nPosY
			break
		}
	}
	protein.Table[posX][posY] = []byte(*proteinChain)[1]
	results := make(chan int)
	go Search(results, protein.Table, protein.Chain, posX, posY, 2, 0, 0, 0)
	for runtime.NumGoroutine() > 0 {
		select {
		case res, ok := <-results:
			if ok {
				fmt.Println(res)
			}
		default:
			continue
		}
	}
}

func Search(results chan int, matrix [][]byte, chain string, posX, posY, k, e, min int, avg float64) {
	availableMoves := GetAvailableMoves(matrix, posX, posY)

	// duplicating a table
	duplicate := Duplicate(matrix)

	for _, move := range availableMoves {
		// TODO new logic
		x, y, err := Move(move, posX, posY, len(duplicate))
		if err != nil {
			continue
		}
		duplicate[x][y] = chain[k]
		energy := CalculateEnergy(matrix, posX, posY, move, e)
		avg := (float64(e + energy)) / float64(k)
		if energy < min {
			min = energy
		}
		if k >= len(chain)-1 {
			//fmt.Println(energy)
			results <- energy
		} else if chain[k] == 'h' {
			if energy < min {
				go Search(results, duplicate, chain, x, y, k+1, energy, min, avg)
			}
			if float64(energy) > avg {
				r := rand.Float64()
				if r > p1 {
					go Search(results, duplicate, chain, x, y, k+1, energy, min, avg)
				}
			}
			if min <= energy && float64(energy) <= avg {
				r := rand.Float64()
				if r > p2 {
					go Search(results, duplicate, chain, x, y, k+1, energy, min, avg)
				}
			}
		} else {
			go Search(results, duplicate, chain, x, y, k+1, energy, min, avg)
		}
	}
}

func Duplicate(matrix [][]byte) [][]byte {
	duplicate := make([][]byte, len(matrix))
	for i := range matrix {
		duplicate[i] = make([]byte, len(matrix[i]))
		copy(duplicate[i], matrix[i])
	}
	return duplicate
}

func CalculateEnergy(proteinTable [][]byte, posX, posY, from, prevEnergy int) (e int) {
	var moves []int
	switch from {
	case 0: //from left
		moves = []int{0, 2, 3}
		break
	case 1: //from right
		moves = []int{1, 2, 3}
		break
	case 2: //from up
		moves = []int{0, 1, 2}
		break
	case 3: //from down
		moves = []int{0, 1, 3}
		break
	}

	for _, move := range moves {
		if x, y, err := Move(move, posX, posY, len(proteinTable)); err != nil {
			continue
		} else if proteinTable[x][y] == 'h' && proteinTable[posX][posY] == 'h' {
			e--
		}
	}
	e += prevEnergy
	return
}

func GetAvailableMoves(proteinTable [][]byte, posX, posY int) (moves []int) {
	for i := 0; i < 4; i++ {
		if x, y, err := Move(i, posX, posY, len(proteinTable)); err != nil {
			continue
		} else if proteinTable[x][y] != 0 {
			continue
		}
		moves = append(moves, i)
	}
	return
}
