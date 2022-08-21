package internal

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

//******************************************************************
//		Banner stuff
//******************************************************************

// A banner is for strings. Style how you want.
type Banner struct {
	// The raw text in the banner
	text string
	// The style to apply to the text
	style lipgloss.Style
}

// Return the stylized view of the banner
func (b Banner) View() string {
	return b.style.Render(b.text)
}

//******************************************************
//		Notice stuff
//	This area displays game messages to the player
//******************************************************

var noticeStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(Colors["Text"]))

func NewNotice() Banner {
	return Banner{
		text:  "",
		style: noticeStyle,
	}
}

//******************************************************************
//		Tile stuff
//  A Tile can draw 1 letter and be stylized. Very flexible.
//******************************************************************

type Tile struct {
	// The single "letter" in this tile
	// The definition of letter is real loose...
	letter string
	// The style to apply to the letter
	style lipgloss.Style
}

// Return a new stylized Tile
func NewTile(l string, s lipgloss.Style) Tile {
	return Tile{
		letter: l,
		style:  s,
	}
}

// ******************************************************
//
//			Board stuff
//	 Where the correct letters are hidden and revealed on
//	 correct player guesses
//
// ******************************************************
// What to show as "blank" before the tile has been guessed
var blankBoardTile = " "

var boardTileStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(Colors["Text"])).
	Background(lipgloss.Color(Colors["StrongMauve"])).
	Width(5).
	Align(lipgloss.Center)

// The Board is just a collection of Tiles
type Board []Tile

// Make a new Board of n blankBoardTile
func NewBoard(n int) Board {
	b := make([]Tile, n)
	for i := 0; i < n; i++ {
		b[i] = NewTile(blankBoardTile, boardTileStyle)
	}
	return Board(b)
}

// Return the stylized view of the board
// Choose how you want the Tiles to separated from each other
func (b Board) View(sep string) string {
	// Render each board tile and stick in a list
	var result []string
	for _, tile := range b {
		result = append(result, tile.style.Render(tile.letter))
	}
	// Return one giant string that is the board
	return strings.Join(result, sep)
}

// Check if a Tile is in the Board
func (b Board) Contains(s string) bool {
	// Crawl through the board and see if there's a hit
	for _, tile := range b {
		if tile.letter == s {
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
type Keyboard struct {
	// The keyboard alphabet to show. Each row will be a Board
	alphabet []Board
	// The styles to apply when the letter has been used or not
	onStyle  lipgloss.Style
	offStyle lipgloss.Style
}

var letterOffStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(Colors["Base"])).
	Background(lipgloss.Color(Colors["Mauve"])).
	Width(3).
	Align(lipgloss.Center)

var letterOnStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(Colors["Mauve"])).
	Background(lipgloss.Color(Colors["Base"])).
	Width(3).
	Align(lipgloss.Center).
	Bold(true)

var keyboardRows = [][]string{
	{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"},
	{"A", "S", "D", "F", "G", "H", "J", "K", "L"},
	{"Z", "X", "C", "V", "B", "N", "M"},
}

func NewKeyboard() Keyboard {
	// Build up a 2d array of Tiles
	var alphabetTiles = make([]Board, 3)
	for i, row := range keyboardRows {
		alphabetTiles[i] = make(Board, len(row))
		for j, letter := range row {
			// Create the new Tile, ensuring initial style is off
			tile := NewTile(letter, letterOffStyle)
			alphabetTiles[i][j] = tile
		}
	}
	return Keyboard{
		alphabet: alphabetTiles,
		onStyle:  letterOnStyle,
		offStyle: letterOffStyle,
	}
}

// Call View to see stylized string representation of Keyboard
func (keyboard Keyboard) View() string {
	var result []string
	// Each row is a Board, so it can be easily Viewed
	for _, row := range keyboard.alphabet {
		result = append(result, row.View(""))
	}

	return lipgloss.NewStyle().
		// Give the keyboard some room or it crowds the hangman dude
		MarginLeft(4).
		Render(
			// Combine the keyboard rows into a stack
			lipgloss.JoinVertical(lipgloss.Center, result...),
		)
}

// Find this letter in the Keyboard struct and flip it's style between off/on
func (letters *Keyboard) FlipOn(letter string) {
	for i, row := range letters.alphabet {
		for j, tile := range row {
			if tile.letter == letter {
				letters.alphabet[i][j].style = letters.onStyle
				break
			}
		}
	}
}
