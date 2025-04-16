package main

import (
	"fmt"
	"strings"
)

// This is just to print with pretty format for console using ASCII UniCode

// ANSI color codes
const (
	Reset        = "\033[0m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Purple       = "\033[35m"
	Cyan         = "\033[36m"
	Gray         = "\033[37m"
	DarkGray     = "\033[90m"
	LightRed     = "\033[91m"
	LightGreen   = "\033[92m"
	LightYellow  = "\033[93m"
	LightBlue    = "\033[94m"
	LightMagenta = "\033[95m"
	LightCyan    = "\033[96m"
	White        = "\033[97m"
)

type Color struct {
	Code string
}

var (
	ColorRed          = Color{Red}
	ColorGreen        = Color{Green}
	ColorYellow       = Color{Yellow}
	ColorBlue         = Color{Blue}
	ColorPurple       = Color{Purple}
	ColorCyan         = Color{Cyan}
	ColorGray         = Color{Gray}
	ColorDarkGray     = Color{DarkGray}
	ColorLightRed     = Color{LightRed}
	ColorLightGreen   = Color{LightGreen}
	ColorLightYellow  = Color{LightYellow}
	ColorLightBlue    = Color{LightBlue}
	ColorLightMagenta = Color{LightMagenta}
	ColorLightCyan    = Color{LightCyan}
	ColorWhite        = Color{White}
	ColorNothing      = Color{""}
)

type BoxChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

var UnicodeBox = BoxChars{
	TopLeft:     "┌",
	TopRight:    "┐",
	BottomLeft:  "└",
	BottomRight: "┘",
	Horizontal:  "─",
	Vertical:    "│",
}

func colored(text string, color Color) string {
	return color.Code + text + Reset
}

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value string
}

func printBoxedUnicode(heading []string, headingColors []Color, content []KeyValue, keyColor Color, valueColor Color) {
	// Find the longest key for consistent spacing.
	longestKey := 0
	for _, kv := range content {
		if len(kv.Key) > longestKey {
			longestKey = len(kv.Key)
		}
	}

	// Calculate the box width.
	boxWidth := 80 // Fixed box width

	// Print the top border.
	fmt.Println(UnicodeBox.TopLeft + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.TopRight)

	// Print the heading.
	for i, line := range heading {
		headingPadding := (boxWidth - len(line) - 2) / 2
		if headingPadding < 0 {
			headingPadding = 0
		}
		var currentColor Color
		if i < len(headingColors) {
			currentColor = headingColors[i]
		} else {
			currentColor = ColorWhite // Default color
		}
		fmt.Printf("%s%s%s%s%s\n",
			UnicodeBox.Vertical,
			strings.Repeat(" ", headingPadding),
			colored(line, currentColor),
			strings.Repeat(" ", boxWidth-len(line)-2-headingPadding),
			UnicodeBox.Vertical)
	}

	// extra line between heading and content
	// fmt.Println(UnicodeBox.Vertical + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.Vertical) // add line
	fmt.Println(UnicodeBox.Vertical + strings.Repeat(" ", boxWidth-2) + UnicodeBox.Vertical) // just space

	// Calculate fixed positions
	firstColStart := 2               // After border and space
	secondColStart := boxWidth/2 - 3 // Fixed position for second column

	// Print content in pairs
	for i := 0; i < len(content); i += 2 {
		kv1 := content[i]
		line := UnicodeBox.Vertical + " " // Start with left border and space

		// First column
		line += colored(kv1.Key, keyColor) + strings.Repeat(" ", longestKey-len(kv1.Key)) + " : " + colored(kv1.Value, valueColor)

		if i+1 < len(content) {
			// Second column exists
			kv2 := content[i+1]

			// Calculate padding between first column and second column
			currentPos := firstColStart + longestKey + 3 + len(kv1.Value) // 3 for " : "
			middlePadding := secondColStart - currentPos
			line += strings.Repeat(" ", middlePadding)

			// Second column
			line += colored(kv2.Key, keyColor) + strings.Repeat(" ", longestKey-len(kv2.Key)) + " : " + colored(kv2.Value, valueColor)

			// Add padding to right border
			currentPos = secondColStart + longestKey + 3 + len(kv2.Value) // 3 for " : "
			rightPadding := boxWidth - currentPos - 1                     // -1 for border
			line += strings.Repeat(" ", rightPadding)
		} else {
			// Single column - pad to right border
			remainingSpace := boxWidth - firstColStart - longestKey - len(kv1.Value) - 3 - 1 // -3 for " : ", -1 for border
			line += strings.Repeat(" ", remainingSpace)
		}

		line += UnicodeBox.Vertical
		fmt.Println(line)
	}

	// Print the bottom border.
	fmt.Println(UnicodeBox.BottomLeft + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.BottomRight)
}

// Here is the example to call this.
func TestBoxPrint() {
	appName := []string{"My Awesome App", "Configuration", "More Settings"}
	headingColors := []Color{
		ColorCyan,
		ColorGreen,
		ColorNothing,
	}

	// Content defined in order
	appSettings := []KeyValue{
		{"Version", "1.2.3"},
		{"Environment", "Production"},
		{"Database", "PostgreSQL"},
		{"API Key", "********************"},
		{"Server IP", "192.168.1.100"},
		{"Port", "8080"},
		{"Debug", "true"},
		{"Timeout", "10s"},
		{"Extra", "value"},
	}

	keyColor := ColorYellow
	valueColor := ColorWhite

	printBoxedUnicode(appName, headingColors, appSettings, keyColor, valueColor)
}
