package helpers

const Red = "\033[31m"
const Green = "\033[32m"
const GreenBg = "\033[42m"
const Yellow = "\033[33m"
const YellowBg = "\033[43m"
const Blue = "\033[34m"
const Black = "\033[30m"
const Magenta = "\033[35m"
const Cyan = "\033[36m"
const CyanBg = "\033[46m"
const White = "\033[37m"
const Gray = "\033[90m"
const GrayBg = "\033[100m"
const Orange = "\033[38;5;208m"
const Reset = "\033[0m"
const Bold = "\033[1m"
const Underline = "\033[4m"

func Colorize(text string, colors ...string) string {
	for _, color := range colors {
		text = color + text + Reset
	}
	return text
}
