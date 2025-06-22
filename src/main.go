package main

import (
	"fmt"

	"github.com/qwertycodeqc/locc/src/colors"
)

func main() {
	fmt.Println(colors.Colorize("locc", colors.Cyan, colors.Bold))
}