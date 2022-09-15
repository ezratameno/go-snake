// package terminal will be in charge with interacting with the terminal
package terminal

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	screen *bufio.Writer
}

func NewTerminal() *Terminal {
	return &Terminal{
		screen: bufio.NewWriter(os.Stdout),
	}
}

func (t *Terminal) HideCursor() {
	fmt.Fprint(t.screen, "\033[?25l")
}

func (t *Terminal) ShowCursor() {
	fmt.Fprint(t.screen, "\033[?25h")
}

func (t *Terminal) ClearScreen() {
	fmt.Fprint(t.screen, "\033[2J")

}

// MoveCursor recive x,y cordinates where to move the cursor.
func (t *Terminal) MoveCursor(pos [2]int) {
	fmt.Fprintf(t.screen, "\033[%d;%dH", pos[1], pos[0])
}

func (t *Terminal) Draw(str string) {
	fmt.Fprint(t.screen, str)
}

// Render will Render to the screen
func (t *Terminal) Render() {
	t.screen.Flush()
}

// GetSize returns the size of the screen.
// note: top left of the screen is 1,1 not 0,0.
func (t *Terminal) GetSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}
