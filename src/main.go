package main

import (
	"fmt"

	"github.com/qwertycodeqc/locc/src/colors"
)

const VERSION = "1.0"

func main() {
	fmt.Println(colors.Colorize("locc", colors.Cyan, colors.Bold), "v"+VERSION)
}
