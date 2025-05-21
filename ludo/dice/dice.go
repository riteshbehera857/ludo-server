package dice

import (
	"rng"
)

type Dice struct{}

func (d *Dice) Roll() int {
	rng := &rng.RNG{}
	val := rng.GenerateRandomNumber()
	return val
}
