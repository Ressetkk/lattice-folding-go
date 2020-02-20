package main

import (
	"flag"
	"log"
	"math/rand"
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
	for {
	}
}

func Search(results chan int, matrix [][]byte, chain string, posX, posY, k, e, min int, avg float64) {
	availableMoves := GetAvailableMoves(matrix, [2]int{posX, posY})
	// duplicating a table
	duplicate := make([][]byte, len(matrix))
	for i := range matrix {
		duplicate[i] = make([]byte, len(matrix[i]))
		copy(duplicate[i], matrix[i])
	}

	for _, _ = range availableMoves {
		// TODO new logic
	}
}

func CalculateEnergy(proteinTable [][]byte, pos [2]int, from int) (e int) {
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
	check := func() int {
		e := 0
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
				if proteinTable[next[0]][next[1]] == 'h' && proteinTable[pos[0]][pos[1]] == 'h' {
					e--
				}
			}
		}
		return e
	}
	e = check()
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
