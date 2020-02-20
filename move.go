package main

import "fmt"

func MoveLeft(posX, posY, k int) (int, int, error) {
	posY--
	return posX, posY, checkIfInBounds(posY, k)
}

func MoveRight(posX, posY, k int) (int, int, error) {
	posY++
	return posX, posY, checkIfInBounds(posY, k)
}

func MoveUp(posX, posY, k int) (int, int, error) {
	posX--
	return posX, posY, checkIfInBounds(posX, k)
}

func MoveDown(posX, posY, k int) (int, int, error) {
	posX++
	return posX, posY, checkIfInBounds(posX, k)
}

func checkIfInBounds(pos, k int) (err error) {
	if !(pos < k && pos >= 0) {
		err = fmt.Errorf("element out of bound")
	}
	return
}

func Move(move, posX, posY, k int) (int, int, error) {
	switch move {
	case 0:
		return MoveLeft(posX, posY, k)
	case 1:
		return MoveRight(posX, posY, k)
	case 2:
		return MoveUp(posX, posY, k)
	case 3:
		return MoveDown(posX, posY, k)
	}
	return posX, posY, nil
}
