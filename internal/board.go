package internal

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

/*
A single letter placement on the game board is a tile.

    ╭───────╮
    │       │
    │       ◯
    │      ╱│╲
    │       │
    │

_ _ _ _ _ _ <- several tiles. Together, they are a Board
*/

type Tile string

var BlankSpace = "_____"

// Return a new stylized Tiles
func NewTile(letter string) Tile {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Width(5).
		Align(lipgloss.Center)

	return Tile(style.Render(letter))
}

type Board []Tile

// Make a new Board of n BlankSpaces
func NewBoard(n int) Board {
	b := make([]Tile, n)
	for i := 0; i < n; i++ {
		b[i] = NewTile(BlankSpace)
	}
	return Board(b)
}

// Return a string representation of the board
func (b Board) String() string {
	var result []string
	for _, tile := range b {
		result = append(result, string(tile))
	}

	return strings.Join(result, " ")
}

// Check if a Tile is in the Board
func (b Board) Contains(t Tile) bool {
	for _, tile := range b {
		if tile == t {
			return true
		}
	}
	return false
}
