package internal

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"

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
//	Model stuff
//
// ******************************************************************
type model struct {
	// Struct for all things related to the hangman graphic
	graphicView *GraphicView
	// The word the player is trying to guess
	word string
	// The "board" under the graphic where player guesses are shown
	board Board
	// Text area where player types their guesses
	input textinput.Model
	// All the letters the player has guessed
	userGuesses []string
	// All the possible letters that can be guessed
	keyboard     *Keyboard
	showKeyboard bool
	// The notice area thing
	notice PrettyString
	// Did game end?
	gameOver bool
	// Title banner
	title     PrettyString
	showTitle bool
	// Footer banner area
	footer PrettyString
	// Dimensions of terminal windows
	height int
	width  int
	// If the board is cut off, this is how many tiles are being cut
	numCutoffTiles int
	// Any errors caught go here and should be reported somewhere
	err error
}

func initialModel() model {
	// Get random word from dictionary
	word := dictionary[rand.Intn(len(dictionary))]

	// Make a new board based on word length
	board := NewBoard(len(word), boardTileStyle)

	// New input area
	textInput := newInput()

	// Empty list to hold userGuesses
	var userGuesses []string

	// Graphic stuff
	graphicView := NewGraphicView()

	notice := NewNotice()

	keyboard := NewKeyboard()

	title := NewTitle()

	footer := NewFooter()

	return model{
		graphicView:    &graphicView,
		word:           word,
		board:          board,
		input:          textInput,
		userGuesses:    userGuesses,
		keyboard:       &keyboard,
		showKeyboard:   true,
		notice:         notice,
		gameOver:       false,
		title:          title,
		showTitle:      true,
		footer:         footer,
		height:         0,
		width:          0,
		numCutoffTiles: 0,
		err:            nil,
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

// Update model based on user guess
func handleGuess(m *model) {
	// Did the player enter anything before pressing return?
	if m.input.Value() == "" {
		return
	}
	// Reset notice content for next render
	// Putting it here means it only clears when the user guesses again.
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
			// The guess is a hit! Start "flipping" tiles
			for _, id := range ids {
				m.board[id].text = guess
			}
			// Update model to flash for correct guess on next render
			m.graphicView.flash = true
			m.graphicView.flashStyle = flashCorrectStyle
		} else {
			// Wrong guess! increment graphics
			graphic, err := m.graphicView.graphicGenerator()
			// Update model to flash for incorrect guess on next render
			m.graphicView.flash = true
			m.graphicView.flashStyle = flashWrongStyle
			if err != nil {
				// No more graphics to get. Player loses!
				m.notice.text = fmt.Sprintf("You lose :(\nThe hidden word was: %s", m.word)
				m.notice.style = loseNoticeStyle
				m.gameOver = true
			} else {
				m.graphicView.currentGraphic.text = graphic
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
		m.notice.style = winNoticeStyle
		m.gameOver = true
	}
}

// Update model based on terminal resizing.
// Clear the screen if required.
func handleScreenResize(m *model) {
	// Hide keyboard if there isn't enough room
	maxWidth := lipgloss.Width(m.graphicView.View()) + lipgloss.Width(m.keyboard.View())
	if m.width < maxWidth {
		m.showKeyboard = false
		ClearScreen()
	} else {
		m.showKeyboard = true
	}

	// Hide title if there isn't enough room
	maxWidth = lipgloss.Width(m.title.View())
	if m.width < maxWidth {
		m.showTitle = false
		ClearScreen()
	} else {
		m.showTitle = true
	}

	// Count how many tiles are cut off if there isn't enough room
	maxWidth = lipgloss.Width(m.board.View(" "))
	if m.width < maxWidth {
		tileSize := lipgloss.Width(m.board[0].View())
		m.numCutoffTiles = (maxWidth-m.width)/tileSize + 1
		ClearScreen()
	} else {
		m.numCutoffTiles = 0
	}

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Clear out any flash status. This line is what makes it flash!
	if m.graphicView.flash {
		m.graphicView.ResetFlash()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// The player has guessed something. Process it.
			handleGuess(&m)
			if m.gameOver {
				// TODO: Could allow "play again" feature
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Width
		handleScreenResize(&m)

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
	// Build up pieces for top half of view
	// Get the title
	title := ""
	if m.showTitle {
		title = m.title.View()
	}

	keyboardElement := ""
	if m.showKeyboard {
		keyboardElement = m.keyboard.View()
	}

	// Combine the graphic and keyboard components
	midView := lipgloss.JoinHorizontal(lipgloss.Center, m.graphicView.View(), keyboardElement)

	// Format components together to be aligned
	s := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		midView,
	)

	// Render the board where the word is revealed as player makes correct guess
	// Wrap the tiles if the window is too small
	board := m.board.View(" ")
	if m.numCutoffTiles > 0 {
		// Wrap effect is done by inserting newlines at the cutoff point
		wrappedBoard := slices.Insert(m.board, len(m.board)-m.numCutoffTiles, PrettyString{"\n\n", lipgloss.NewStyle()})
		board = wrappedBoard.View(" ")
	}
	s += "\n\n" + board

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

	s += "\n"

	// footer
	s += m.footer.View()

	return s
}

// ******************************************************************
//
//		Clear screen logic
//	Thanks: https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
//	TODO: Library termenv which we already import seems to support this: https://pkg.go.dev/github.com/muesli/termenv#readme-screen
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
