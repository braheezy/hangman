package internal

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

//******************************************************************
//		Banner stuff
//******************************************************************

type Banner struct {
	content string
	style   lipgloss.Style
}

func NewBanner() Banner {
	return Banner{
		content: "",
		style: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")),
	}
}

// Return a string representation of the banner
func (b Banner) String() string {
	return string(b.style.Render(b.content))
}

//******************************************************************
//		Board & Tile stuff
//******************************************************************

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
var BlankSpace = "_____"

type Tile struct {
	content string
	style   lipgloss.Style
}

var tileStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")). // white-ish
	Background(lipgloss.Color("#7D56F4")). // purple-ish
	Width(5).
	Align(lipgloss.Center)

// Return a new stylized Tiles
func NewTile() Tile {
	return Tile{
		content: BlankSpace,
		style:   tileStyle,
	}
}

type Board []Tile

// Make a new Board of n BlankSpaces
func NewBoard(n int) Board {
	b := make([]Tile, n)
	for i := 0; i < n; i++ {
		b[i] = NewTile()
	}
	return Board(b)
}

// Return the stylized view of the board
func (b Board) View() string {
	var result []string
	for _, tile := range b {
		result = append(result, tile.style.Render(tile.content))
	}
	return strings.Join(result, " ")
}

// Check if a Tile is in the Board
func (b Board) Contains(s string) bool {
	for _, tile := range b {
		if tile.content == s {
			return true
		}
	}
	return false
}

//******************************************************************
//		Letters view stuff
//******************************************************************

/*
A view into the letters the player has already guessed.

Currently, it's a "keyboard" and the letters are marked off as they are guessed.
*/
type Letters struct {
	alphabet []Board
	onStyle  lipgloss.Style
	offStyle lipgloss.Style
}

var letterOffStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#1e1e2e")).
	Background(lipgloss.Color("#89b4fa")).
	Width(3).
	Align(lipgloss.Center)

var letterOnStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#b4befe")).
	Background(lipgloss.Color("#313244")).
	Width(3).
	Align(lipgloss.Center).
	Bold(true)

var keyboardRows = [][]string{
	{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"},
	{"A", "S", "D", "F", "G", "H", "J", "K", "L"},
	{"Z", "X", "C", "V", "B", "N", "M"},
}

func NewLetters() Letters {
	var alphabetTiles = make([]Board, 3)
	for i, row := range keyboardRows {
		alphabetTiles[i] = make(Board, len(row))
		for j, letter := range row {
			tile := NewTile()
			tile.content = letter
			tile.style = letterOffStyle
			alphabetTiles[i][j] = tile
		}
	}
	return Letters{
		alphabet: alphabetTiles,
		onStyle:  letterOnStyle,
		offStyle: letterOffStyle,
	}
}

func (letters Letters) View() string {
	var result []string
	for _, row := range letters.alphabet {
		result = append(result, row.View())
	}
	return lipgloss.NewStyle().
		MarginLeft(4).
		Render(
			lipgloss.JoinVertical(lipgloss.Center, result...),
		)
}

// Find this letter in the Letters struct and flip it's style between off/on
func (letters *Letters) FlipOn(letter string) {
	for i, row := range letters.alphabet {
		for j, tile := range row {
			if tile.content == letter {
				letters.alphabet[i][j].style = letters.onStyle
				break
			}
		}
	}
}
