package internal

import (
	"errors"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// ******************************************************************
//
//	PrettyString stuff
//
// Take one or more characters and style with lipgloss
// ******************************************************************
type PrettyString struct {
	// The raw text in the banner
	text string
	// The style to apply to the text
	style lipgloss.Style
}

// Return a new stylized Tile
func NewPrettyString(t string, s lipgloss.Style) PrettyString {
	return PrettyString{
		text:  t,
		style: s,
	}
}

// Return the stylized view of the banner
func (s PrettyString) View() string {
	return s.style.Render(s.text)
}

// ******************************************************
//
//		Title stuff
//	The top greeter
//
// ******************************************************
var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(textColor).
	Background(primaryColor).
	PaddingLeft(2).
	PaddingRight(2).
	MarginBottom(1)

func NewTitle() PrettyString {
	return PrettyString{
		text:  "Hangman\nCan you save this criminal?",
		style: titleStyle,
	}
}

// ******************************************************
//
//		Footer stuff
//	The top greeter
//
// ******************************************************
var footerStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Underline(true)

func NewFooter() PrettyString {
	return PrettyString{
		text:  "Press ESC or Ctrl+C to quit.",
		style: footerStyle,
	}
}

// ******************************************************
//
//		Notice stuff
//	This area displays game messages to the player
//
// ******************************************************
var noticeStyle = lipgloss.NewStyle().
	Bold(true).
	Italic(true).
	Foreground(primaryColor)

var loseNoticeStyle = lipgloss.NewStyle().
	Inherit(noticeStyle).
	Foreground(failColor)

var winNoticeStyle = lipgloss.NewStyle().
	Inherit(noticeStyle).
	Foreground(successColor)

func NewNotice() PrettyString {
	return PrettyString{
		text:  "",
		style: noticeStyle,
	}
}

// ******************************************************
//
//			Board stuff
//	A horizontal grouping of PrettyStrings
//
// ******************************************************
// What to show as "blank" before the tile has been guessed
var blankBoardTile = " "

var boardTileStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(secondaryColor).
	Background(strongColor).
	Width(5).
	Align(lipgloss.Center)

// The Board is just a collection of PrettyStrings
type Board []PrettyString

// Make a new Board of n blank PrettyStrings
func NewBoard(n int, s lipgloss.Style) Board {
	b := make([]PrettyString, n)
	for i := 0; i < n; i++ {
		b[i] = NewPrettyString(blankBoardTile, s)
	}
	return Board(b)
}

// Return the stylized view of the board
// Choose how you want the Tiles to separated from each other
func (b Board) View(sep string) string {
	// Render each board tile and stick in a list
	var result []string
	for _, tile := range b {
		result = append(result, tile.View())
	}
	// Return one giant string that is the board
	return strings.Join(result, sep)
}

// Check if a Tile is in the Board
func (b Board) Contains(s string) bool {
	// Crawl through the board and see if there's a hit
	for _, tile := range b {
		if tile.text == s {
			return true
		}
	}
	return false
}

// ******************************************************************
//
//			Letters view stuff
//	 A view into the letters the player has already guessed.
//	 It's a "keyboard" and the letters are marked off
//
// ******************************************************************
type Keyboard struct {
	// The keyboard alphabet to show. Each row will be a Board
	alphabet []Board
	// The styles to apply when the letter has been used or not
	onStyle  lipgloss.Style
	offStyle lipgloss.Style
}

var letterOffStyle = lipgloss.NewStyle().
	Foreground(textColor).
	Background(primaryColor).
	Width(3).
	Align(lipgloss.Center)

var letterOnStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Background(tertiaryColor).
	Width(3).
	Align(lipgloss.Center).
	Bold(true)

var keyboardRows = [][]string{
	{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"},
	{"A", "S", "D", "F", "G", "H", "J", "K", "L"},
	{"Z", "X", "C", "V", "B", "N", "M"},
}

func NewKeyboard() Keyboard {
	// Each row in the Keyboard is a "Board"
	var alphabetTiles = make([]Board, 3)
	for i, row := range keyboardRows {
		alphabetTiles[i] = NewBoard(len(row), letterOffStyle)
		for j, letter := range row {
			// Create the new Tile, ensuring initial style is off
			// tile := NewTile(letter, letterOffStyle)
			alphabetTiles[i][j].text = letter
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
			if tile.text == letter {
				letters.alphabet[i][j].style = letters.onStyle
				break
			}
		}
	}
}

// ******************************************************************
//
//	Graphic view
//
// The hangman character.
// A "picture" from the graphics file, stuffed into a Tile
// ******************************************************************
type GraphicView struct {
	// The graphic to show. Changes when player is wrong
	currentGraphic PrettyString
	// Call this repeatedly to get the next graphic
	graphicGenerator func() (string, error)
	// If true, flash the graphic
	flash bool
	// The style to apply when flashing
	flashStyle lipgloss.Style
}

var baseGraphicStyle = lipgloss.NewStyle().
	Bold(true).
	Border(lipgloss.RoundedBorder()).
	BorderBackground(backgroundColor).
	Background(backgroundColor)

var graphicStyle = lipgloss.NewStyle().
	Inherit(baseGraphicStyle).
	Padding(1, 3, 1, 3).
	Foreground(primaryColor).
	BorderForeground(primaryColor)

var flashWrongStyle = lipgloss.NewStyle().
	Inherit(baseGraphicStyle).
	Padding(1, 3, 1, 3).
	Foreground(failColor).
	BorderForeground(failColor)

var flashCorrectStyle = lipgloss.NewStyle().
	Inherit(baseGraphicStyle).
	Padding(1, 3, 1, 3).
	Foreground(successColor).
	BorderForeground(successColor)

func NewGraphicView() GraphicView {
	// Set up the generator can call it once to get first graphic.
	graphicGen := Graphics()
	currentGraphic, err := graphicGen()
	if err != nil {
		panic(err)
	}

	return GraphicView{
		currentGraphic: PrettyString{
			text:  currentGraphic,
			style: graphicStyle,
		},
		graphicGenerator: graphicGen,
		// "nil" style
		flashStyle: lipgloss.NewStyle(),
	}
}

func (g *GraphicView) View() string {
	if g.flash {
		g.currentGraphic.style = g.flashStyle
	}
	return g.currentGraphic.View()
}

func (g *GraphicView) ResetFlash() {
	g.currentGraphic.style = graphicStyle
	g.flash = false
}

// ******************************************************************
//
//		Player input area
// A 1-character text input area for the player to make letter guesses
// ******************************************************************

func newInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Guess a letter!"
	ti.Focus()
	ti.CharLimit = 1
	ti.Width = 1

	ti.Validate = validateInput()

	// TODO: Is there a better place to put this?
	ti.Prompt = "─> "
	ti.PromptStyle = lipgloss.NewStyle().
		Foreground(secondaryColor)
	ti.PlaceholderStyle = lipgloss.NewStyle().
		Italic(true).
		Faint(true).
		Foreground(secondaryColor)

	return ti
}

// Only allow letter inputs
func validateInput() textinput.ValidateFunc {
	return func(s string) error {
		letter := rune(s[0])
		if !unicode.IsLetter(letter) {
			return errors.New("not valid input")
		}
		return nil
	}
}
