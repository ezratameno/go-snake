package snake

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/etameno/go-snake/pkg/terminal"
	"github.com/mattn/go-tty"
)

type game struct {
	Score    int
	Snake    *snake
	Food     Position
	terminal *terminal.Terminal
}

type snake struct {
	Body []Position
	// the direction the snake is traveling
	Direction direction
}

// x and y corordinates
type Position [2]int

// enum
type direction int

const (
	North direction = iota
	East
	South
	West
)

func newSnake(t *terminal.Terminal) (*snake, error) {
	maxX, maxY, err := t.GetSize()
	if err != nil {
		return nil, err
	}
	pos := Position{maxX / 2, maxY / 2}

	return &snake{
		Body:      []Position{pos},
		Direction: North,
	}, nil
}

func NewGame(t *terminal.Terminal) (*game, error) {
	rand.Seed(time.Now().UnixNano())
	snake, err := newSnake(t)
	if err != nil {
		return nil, err
	}
	foodRandPos, err := RandomPosition(t)
	if err != nil {
		return nil, err
	}
	game := &game{
		Score:    0,
		Snake:    snake,
		Food:     foodRandPos,
		terminal: t,
	}

	// listen in the bg for key press.
	go game.listenForKeyPress()

	return game, nil
}

// BeforeGame hides the cursor and listen for Interrupt.
func (g *game) BeforeGame() {
	g.terminal.HideCursor()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			g.Over()
		}
	}()
}

// Over will clear the screen and print the user his score.
func (g *game) Over() {
	g.terminal.ClearScreen()

	g.terminal.MoveCursor(Position{1, 1})
	g.terminal.ShowCursor()

	g.terminal.Draw("Game over. score: " + strconv.Itoa(g.Score) + "\n")
	g.terminal.Render()
	time.Sleep(3 * time.Second)
	// reset the terminal
	cmd := exec.Command("reset")
	cmd.Run()
	os.Exit(0)
}

func (g *game) listenForKeyPress() {
	// open the terminal.
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		// read from the terminal the input.
		char, err := tty.ReadRune()
		if err != nil {
			panic(err)
		}
		// UP, DOWN, RIGHT, LEFT == [A, [B, [C, [D
		// we ignore the escape character [
		switch char {
		case 'A':
			g.Snake.Direction = North
		case 'B':
			g.Snake.Direction = South
		case 'C':
			g.Snake.Direction = East
		case 'D':
			g.Snake.Direction = West
		}
	}
}

// Draw the screen.
func (g *game) Draw() error {
	g.terminal.ClearScreen()
	maxX, _, err := g.terminal.GetSize()
	if err != nil {
		return err
	}

	// ===========================================
	// draw the score

	status := "score: " + strconv.Itoa(g.Score)
	// place in the middle
	statusXPos := maxX/2 - len(status)/2

	g.terminal.MoveCursor(Position{statusXPos, 0})
	g.terminal.Draw(status)
	// ===========================================
	// draw the food

	g.terminal.MoveCursor(g.Food)
	g.terminal.Draw("*")

	// ==========================================
	// draw the snake
	for i, pos := range g.Snake.Body {
		g.terminal.MoveCursor(pos)
		if i == 0 {
			g.terminal.Draw("O")
		} else {
			g.terminal.Draw("o")
		}
	}

	g.terminal.Render()
	// the speed of the game.
	time.Sleep(time.Millisecond * 60)
	return nil
}

// PlaceNewFood will place new food in a random location
func (g *game) PlaceNewFood() error {
	for {
		newFoodPosition, err := RandomPosition(g.terminal)
		if err != nil {
			return err
		}

		if PositionsAreSame(newFoodPosition, g.Food) {
			continue
		}

		for _, pos := range g.Snake.Body {
			if PositionsAreSame(newFoodPosition, pos) {
				continue
			}
		}

		g.Food = newFoodPosition

		break
	}

	return nil
}
