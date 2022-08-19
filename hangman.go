package main

import (
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

type errMsg error

//go:embed dictionary.txt
var f embed.FS
var DictionaryFile, _ = f.ReadFile("dictionary.txt")

func LoadWords() (words []string, err error) {
	// Load dictionary into a list and return list
	words = strings.Fields(string(DictionaryFile))
	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}

	return words, nil
}

var dictionary, _ = LoadWords()

//******************************************************************
//		Handle player input
//******************************************************************
// A 1-character text input area for the player to make letter guesses
func newInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Guess a letter!"
	ti.Focus()
	ti.CharLimit = 1
	ti.Width = 1

	ti.Validate = validateInput()

	return ti
}

// Only allow letter inputs
func validateInput() textinput.ValidateFunc {
	var err error
	return func(s string) error {
		if unicode.IsDigit(rune(s[0])) {
			return err
		}
		return nil
	}
}

//******************************************************************
//		Model stuff
//******************************************************************
type model struct {
	// Call this repeatedly to get the next graphic
	graphicGenerator func() string
	// The graphic to show. Changes when player is wrong
	currentGraphic string
	// The word the player is trying to guess
	word string
	// The "board" under the graphic where player guesses are shown
	board Board
	// Text area where player types their guesses
	input textinput.Model
	// All the letters the player has guessed
	guesses []string
	err     error
}

func initialModel() model {
	// Get random word from dictionary
	rand.Seed(time.Now().UnixNano())
	w := dictionary[rand.Intn(len(dictionary))]

	// Make a new board based on word length
	b := NewBoard(len(w))

	// New input area
	ti := newInput()

	// Empty list to hold guesses
	var g []string

	// Graphic stuff
	gg := Graphics()
	cg := gg()

	return model{
		graphicGenerator: gg,
		currentGraphic:   cg,
		word:             w,
		board:            b,
		input:            ti,
		guesses:          g,
		err:              nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

//******************************************************************
//		Update stuff
//******************************************************************
// Return list of indexes where letters occur in string
func Indexes(s string, letter string) []int {
	var indexes []int
	for i, c := range s {
		if string(c) == letter {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		// The player has guessed something. Process it.
		case "enter":
			// Reset err state
			m.err = nil
			// Pull out the letter
			guess := strings.ToUpper(m.input.Value())
			// Can't guess letters already guessed
			if slices.Contains(m.guesses, guess) {
				m.err = errors.New("you already guessed that. try again")
			} else {
				// See if the guess is one of the letters in the word
				ids := Indexes(m.word, guess)
				if len(ids) > 0 {
					// The guess is a hit! Start flipping tiles
					for _, id := range ids {
						m.board[id] = NewTile(guess)
					}
				} else {
					// Wrong guess! increment graphics
					m.currentGraphic = m.graphicGenerator()
				}
				// Remember guesses for next loop
				m.guesses = append(m.guesses, guess)
			}
			// Clear the input area
			m.input.Reset()

			// If there aren't any more blank tiles then word is filled! Winner!
			if !m.board.Contains(NewTile(BlankSpace)) {
				m.err = errors.New("you win! feel free to re-run the program to play again")
				return m, tea.Quit
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

//******************************************************************
//		View stuff
//******************************************************************
func (m model) View() string {
	// Title area
	s := "Play Hangman!"

	// Current hangman graphic is replaced as player makes incorrect guesses
	s += "\n\n" + m.currentGraphic

	// Render the board where the word is revealed as player makes correct guess
	s += "\n\n" + m.board.String()

	// Render the little input area for player guesses
	s += "\n\n" + m.input.View()

	// DEV: This is for debug :)
	// s += fmt.Sprintf("\n\nPsst the word is %s\n\n", m.word)

	s += "\n"

	// TODO: Don't do this lol we're reusing the err thing as a changeable text area to tell the user stuff
	if m.err != nil {
		s += fmt.Sprintf("%v\n", m.err)
	}

	// footer
	s += "\nPress ESC or Ctrl+C to quit.\n"

	return s
}

//******************************************************************
//		Run stuff
//******************************************************************
func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
