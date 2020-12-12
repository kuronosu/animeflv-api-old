package utils

import (
	"log"
	"strings"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"

// ColoredText return text with the specified color
func ColoredText(text string, color string) string {
	switch strings.ToLower(color) {
	case "red":
		return string(colorRed) + text + string(colorReset)
	case "green":
		return string(colorGreen) + text + string(colorReset)
	case "yellow":
		return string(colorYellow) + text + string(colorReset)
	case "blue":
		return string(colorBlue) + text + string(colorReset)
	case "purple":
		return string(colorPurple) + text + string(colorReset)
	case "cyan":
		return string(colorCyan) + text + string(colorReset)
	case "white":
		return string(colorWhite) + text + string(colorReset)
	default:
		return text
	}
}

// InfoLog print log in blue color
func InfoLog(text string) {
	log.Println(ColoredText(text, "blue"))
}

// SuccessLog print log in blue color
func SuccessLog(text string) {
	log.Println(ColoredText(text, "green"))
}

// WarningLog print log in yellow color
func WarningLog(text string) {
	log.Println(ColoredText(text, "yellow"))
}

// ErrorLog print log in red color
func ErrorLog(text string) {
	log.Println(ColoredText(text, "red"))
}

// FatalLog call log.Fatal with red color
func FatalLog(text string) {
	log.Fatal(ColoredText(text, "red"))
}
