package internal

import (
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"
)

type errMsg error

// ******************************************************************
//
//	Dictionary stuff
//
// ******************************************************************
//
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

// ******************************************************************
//
//	Handle player input
//
// ******************************************************************
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
	return func(s string) error {
		letter := rune(s[0])
		if !unicode.IsLetter(letter) {
			return errors.New("not valid input")
		}
		return nil
	}
}

// ******************************************************************
//
//	Model stuff
//
// ******************************************************************
type model struct {
	// Call this repeatedly to get the next graphic
	graphicGenerator func() (string, error)
	// The graphic to show. Changes when player is wrong
	currentGraphic string
	// The word the player is trying to guess
	word string
	// The "board" under the graphic where player guesses are shown
	board Board
	// Text area where player types their guesses
	input textinput.Model
	// All the letters the player has guessed
	userGuesses []string
	// All the possible letters that can be guessed
	keyboard *Keyboard
	// The notice area thing
	notice Banner
	err    error
}

func initialModel() model {
	// Get random word from dictionary
	w := dictionary[rand.Intn(len(dictionary))]

	// Make a new board based on word length
	b := NewBoard(len(w))

	// New input area
	ti := newInput()

	// Empty list to hold userGuesses
	var g []string

	// Graphic stuff
	gg := Graphics()
	cg, _ := gg()

	notice := NewNotice()

	keyboard := NewKeyboard()

	return model{
		graphicGenerator: gg,
		currentGraphic:   cg,
		word:             w,
		board:            b,
		input:            ti,
		userGuesses:      g,
		keyboard:         &keyboard,
		notice:           notice,
		err:              nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// ******************************************************************
//
//	Update stuff
//
// ******************************************************************
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
			// Did the player enter anything?
			if m.input.Value() == "" {
				break
			}
			// Reset notice content
			m.notice.text = ""
			// Pull out the letter
			guess := strings.ToUpper(m.input.Value())
			// Can't guess letters already guessed
			if slices.Contains(m.userGuesses, guess) {
				m.notice.text = "Silly, you already guessed that! Try again"
			} else {
				// See if the guess is one of the letters in the word
				ids := Indexes(m.word, guess)
				if len(ids) > 0 {
					// The guess is a hit! Start flipping tiles
					for _, id := range ids {
						m.board[id].letter = guess
					}
				} else {
					// Wrong guess! increment graphics
					graphic, err := m.graphicGenerator()
					if err != nil {
						// No more graphics to get. Player loses!
						m.notice.text = fmt.Sprintf("You lose :(\nThe word we were looking for: %s", m.word)
						// Looks tacky to leave the last character typed
						// TODO: Surely this can be refactored
						m.input.Reset()
						return m, tea.Quit
					} else {
						m.currentGraphic = graphic
					}
				}
				// Remember userGuesses for next loop
				m.userGuesses = append(m.userGuesses, guess)
				m.keyboard.FlipOn(guess)
			}
			// Clear the input area
			m.input.Reset()

			// If there aren't any more blank tiles then word is filled! Winner!
			if !m.board.Contains(blankBoardTile) {
				m.notice.text = "Woo you win! Feel free to re-run the program to play again!"
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

// ******************************************************************
//
//	View stuff
//
// ******************************************************************
func (m model) View() string {
	// Title area
	s := "Play Hangman!\n\n"

	// Current hangman graphic is replaced as player makes incorrect guesses
	s += lipgloss.JoinHorizontal(lipgloss.Center, m.currentGraphic, m.keyboard.View())

	// Render the board where the word is revealed as player makes correct guess
	s += "\n\n" + m.board.View()

	// Render the little input area for player guesses
	s += "\n\n" + m.input.View()

	// !: This is for debug :)
	// s += fmt.Sprintf("\n\nPsst the word is %s\n\n", m.word)

	s += "\n"

	if m.err != nil {
		s += fmt.Sprintf("%v\n", m.err)
	}
	if m.notice.text != "" {
		s += m.notice.View()
	}

	// footer
	s += "\n\nPress ESC or Ctrl+C to quit.\n"

	return s
}

// ******************************************************************
//
//		Clear screen logic
//	Thanks: https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
//
// ******************************************************************
// Store clear commands for different OSes
var clearFuncs map[string]func() = initClearMap()

func initClearMap() map[string]func() {
	clearMap := make(map[string]func())

	// For each operating system, define a function that will clear the screen.
	clearMap["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["darwin"] = clearMap["linux"]

	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return clearMap
}

func ClearScreen() {
	clearFunction, ok := clearFuncs[runtime.GOOS]
	if ok {
		clearFunction()
	}
}

// ******************************************************************
//
//	Run stuff
//
// ******************************************************************
func Run() {
	// Wipe the current terminal of content for fresh play
	ClearScreen()

	// Start BubbleTea runtime
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
