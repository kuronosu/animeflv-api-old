package utils

import (
	"fmt"
	"log"
)

// NC No color
const NC = "\033[0m"

// Black color
const Black = "\033[0;30m"

// Red color
const Red = "\033[0;31m"

// Green color
const Green = "\033[0;32m"

// BrownOrange color
const BrownOrange = "\033[0;33m"

// Blue color
const Blue = "\033[0;34m"

// Purple color
const Purple = "\033[0;35m"

// Cyan color
const Cyan = "\033[0;36m"

// LightGray color
const LightGray = "\033[0;37m"

// DarkGray color
const DarkGray = "\033[1;30m"

// LightRed color
const LightRed = "\033[1;31m"

// LightGreen color
const LightGreen = "\033[1;32m"

// Yellow color
const Yellow = "\033[1;33m"

// LightBlue color
const LightBlue = "\033[1;34m"

// LightPurple color
const LightPurple = "\033[1;35m"

// LightCyan color
const LightCyan = "\033[1;36m"

// White color
const White = "\033[1;37m"

// Colored return text with the specified Color
func Colored(color string, v ...interface{}) string {
	return color + fmt.Sprint(v...) + NC
}

// ColoredLog print log in specified Color
func ColoredLog(color string, v ...interface{}) {
	log.Println(Colored(color, v...))
}

// InfoLog print log in blue color
func InfoLog(v ...interface{}) {
	log.Println(Colored(LightBlue, v...))
}

// SuccessLog print log in blue color
func SuccessLog(v ...interface{}) {
	log.Println(Colored(LightGreen, v...))
}

// WarningLog print log in yellow color
func WarningLog(v ...interface{}) {
	log.Println(Colored(Yellow, v...))
}

// ErrorLog print log in red color
func ErrorLog(v ...interface{}) {
	log.Println(Colored(LightRed, v...))
}

// FatalLog call log.Fatal with red color
func FatalLog(v ...interface{}) {
	log.Fatal(Colored(Red, v...))
}
