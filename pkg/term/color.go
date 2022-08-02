package term

import "github.com/fatih/color"

type ColoredOutput interface {
	PrintRedString(string) string
	PrintBlueString(string) string
	PrintGreenString(string) string
}

type Color struct{}

func (c Color) PrintRedString(content string) string {
	return color.RedString(content)
}

func (c Color) PrintBlueString(content string) string {
	return color.BlueString(content)
}

func (c Color) PrintGreenString(content string) string {
	return color.GreenString(content)
}
