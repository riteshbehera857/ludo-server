package rng

import "math/rand"

var rollToggle = false

type RNG struct{}

func (rng *RNG) GenerateRandomNumber() int {
	rollToggle = !rollToggle
	origVal := rand.Intn(6) + 1
	val := origVal
	if rollToggle && origVal < 4 {
		val = origVal + 3
	}

	return val
}
