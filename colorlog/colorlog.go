package colorlog

import (
	"fmt"
	"github.com/fatih/color"
)

func Success(msg string) {
	Green(msg)
}

func Warn(msg string) {
	Yellow(msg)
}

func Error(msg string) {
	Red(msg)
}

func Black(msg string) {
	color.Set(color.FgBlack, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func Red(msg string) {
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func Green(msg string) {
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func Yellow(msg string) {
	color.Set(color.FgYellow, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func Blue(msg string) {
	color.Set(color.FgHiBlue, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func Magenta(msg string) {
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func Cyan(msg string) {
	color.Set(color.FgCyan, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}

func White(msg string) {
	color.Set(color.FgWhite, color.Bold)
	defer color.Unset()
	fmt.Printf("%v\n", msg)
}
