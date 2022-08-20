package main

import "github.com/charmbracelet/lipgloss"

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
