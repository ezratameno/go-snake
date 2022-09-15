package main

import (
	"fmt"
	"os"

	"github.com/etameno/go-snake/internal/snake"
	"github.com/etameno/go-snake/pkg/terminal"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func run() error {
	t := terminal.NewTerminal()
	game, err := snake.NewGame(t)
	if err != nil {
		return err
	}
	game.BeforeGame()
	for {
		maxX, maxY, err := t.GetSize()
		if err != nil {
			return err
		}
		// calculate new head position
		newHeadPos := game.Snake.Body[0]

		// note the top left corner is 1,1.
		switch game.Snake.Direction {
		case snake.North:
			newHeadPos[1]--
		case snake.East:
			newHeadPos[0]++
		case snake.South:
			newHeadPos[1]++
		case snake.West:
			newHeadPos[0]--
		}

		// if you hit the wall, game over
		hitWall := newHeadPos[0] < 1 || newHeadPos[1] < 1 || newHeadPos[0] > maxX ||
			newHeadPos[1] > maxY
		if hitWall {
			game.Over()
		}

		// if you run into yourself, game over
		// go through all the points of the snake to see if the new head pos
		// touch any of them
		for _, pos := range game.Snake.Body {
			if snake.PositionsAreSame(newHeadPos, pos) {
				game.Over()
			}
		}

		// add the new head to the body
		game.Snake.Body = append([]snake.Position{newHeadPos}, game.Snake.Body...)

		ateFood := snake.PositionsAreSame(game.Food, newHeadPos)
		if ateFood {
			game.Score++
			game.PlaceNewFood()
		} else {
			// if the snake didn't eat we remove the last elment in the snake so it will seems like it's moving.
			game.Snake.Body = game.Snake.Body[:len(game.Snake.Body)-1]
		}

		err = game.Draw()
		if err != nil {
			return err
		}
	}
}
