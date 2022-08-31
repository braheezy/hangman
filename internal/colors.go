package internal

import "github.com/charmbracelet/lipgloss"

// https://github.com/catppuccin/catppuccin
// Mostly the Latte (Light) and Macchiato (Dark) flavors.
var DarkColors = map[string]string{
	"Rosewater":   "#f5e0dc",
	"Flamingo":    "#f2cdcd",
	"Pink":        "#f5c2e7",
	"Mauve":       "#cba6f7",
	"StrongMauve": "#8839ef",
	"Red":         "#f38ba8",
	"Maroon":      "#eba0ac",
	"Peach":       "#fab387",
	"Yellow":      "#f9e2af",
	"Green":       "#a6e3a1",
	"Teal":        "#8bd5ca",
	"Sky":         "#91d7e3",
	"Sapphire":    "#7dc4e4",
	"Blue":        "#8aadf4",
	"Lavender":    "#b7bdf8",
	"Text":        "#cad3f5",
	"Subtext1":    "#b8c0e0",
	"Subtext0":    "#a5adcb",
	"Overlay2":    "#939ab7",
	"Overlay1":    "#8087a2",
	"Overlay0":    "#6e738d",
	"Surface2":    "#5b6078",
	"Surface1":    "#494d64",
	"Surface0":    "#363a4f",
	"Base":        "#24273a",
	"DarkText":    "#24273a",
	"Mantle":      "#1e2030",
	"Crust":       "#181926",
}

var LightColors = map[string]string{
	"Rosewater": "#dc8a78",
	"Flamingo":  "#dd7878",
	"Pink":      "#ea76cb",
	"Mauve":     "#8839ef",
	"Red":       "#d20f39",
	"Maroon":    "#e64553",
	"Peach":     "#fe640b",
	"Yellow":    "#df8e1d",
	"Green":     "#40a02b",
	"Teal":      "#179299",
	"Sky":       "#04a5e5",
	"Sapphire":  "#209fb5",
	"Blue":      "#1e66f5",
	"Lavender":  "#7287fd",
	"Text":      "#4c4f69",
	"Subtext1":  "#5c5f77",
	"Subtext0":  "#6c6f85",
	"Overlay2":  "#7c7f93",
	"Overlay1":  "#8c8fa1",
	"Overlay0":  "#9ca0b0",
	"Surface2":  "#acb0be",
	"Surface1":  "#bcc0cc",
	"Surface0":  "#ccd0da",
	"Base":      "#eff1f5",
	"Mantle":    "#e6e9ef",
	"Crust":     "#dce0e8",
}

/*
	I have no UI skills and have no idea what best practices are
	when it comes to color code management.

	Below has the appearance of a good scheme. Just change the colors here
	and changes magically propagate where needed.

	This all falls apart when the aspiring designer realizes there is no rhyme
	or reason for why one color is named such, then applied to some elements.

	Some colors here are only used once! The madness...
*/
// Controls the base color used in the app:
//   - Text of letters in keyboard during On state
//   - Background of letters in keyboard during Off state
//   - Hangman graphic text and border color
//   - Background to title text area
//   - Footer text
//   - Notice area text (neutral messages)
var primaryColor = lipgloss.AdaptiveColor{
	Light: LightColors["Lavender"],
	Dark:  DarkColors["Mauve"],
}

// Secondary complimentary color:
//   - Board text color
//   - Player input area
var secondaryColor = lipgloss.AdaptiveColor{
	Light: LightColors["Pink"],
	Dark:  DarkColors["Pink"],
}

// Third complimentary color:
//   - Keyboard letter background during On state
var tertiaryColor = lipgloss.AdaptiveColor{
	Light: LightColors["Overlay0"],
	Dark:  DarkColors["Overlay0"],
}

// Controls the strong color used in the app:
//   - The background of the board where word is revealed
var strongColor = lipgloss.AdaptiveColor{
	Light: DarkColors["StrongMauve"],
	Dark:  DarkColors["StrongMauve"],
}

// Color for positive feedback:
//   - The color of the win game text
//   - The graphic flash color on successful guess
var successColor = lipgloss.AdaptiveColor{
	Light: LightColors["Green"],
	Dark:  DarkColors["Green"],
}

// Color for negative feedback:
//   - The color of the lose game text
//   - The graphic flash color on incorrect guess
var failColor = lipgloss.AdaptiveColor{
	Light: LightColors["Red"],
	Dark:  DarkColors["Red"],
}

// Controls the text color for:
//   - Title text
//   - Text of letters in keyboard during Off state
var textColor = lipgloss.AdaptiveColor{
	Light: LightColors["Text"],
	Dark:  DarkColors["DarkText"],
}

// Controls color for background of hangman graphic
var backgroundColor = lipgloss.AdaptiveColor{
	Light: LightColors["Base"],
	Dark:  DarkColors["Base"],
}
