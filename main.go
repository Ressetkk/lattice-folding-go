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
	Search(results, protein.Table, protein.Chain, posX, posY, 2, 0, 0, 0)

}

func Search(results chan int, matrix [][]byte, chain string, posX, posY, k, e, min int, avg float64) {
	availableMoves := GetAvailableMoves(matrix, posX, posY)

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

func CalculateEnergy(proteinTable [][]byte, posX, posY, from int) (e int) {

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
