package snake

import (
	"math/rand"

	"github.com/etameno/go-snake/pkg/terminal"
)

func PositionsAreSame(a, b Position) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func RandomPosition(t *terminal.Terminal) (Position, error) {
	width, height, err := t.GetSize()
	if err != nil {
		return Position{}, err
	}
	x := rand.Intn(width) + 1
	y := rand.Intn(height) + 2

	return [2]int{x, y}, nil
}
