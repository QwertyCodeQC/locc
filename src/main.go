package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/qwertycodeqc/locc/src/helpers"
)

const VERSION = "1.0"
var IGNORE = []string{".git", "node_modules", "dist", "coverage"}
var BINARY_EXTENSIONS = []string{"*.exe", "*.dll", "*.so", "*.dylib", "*.o", "*.a", "*.lib", "*.class", "*.jar", "*.pyc", "*.pyo"}

func main() {
	fmt.Println(helpers.Colorize("locc", helpers.Cyan, helpers.Bold), "v"+VERSION)
	loadIgnore()
}

func loadIgnore() {
	path, err := helpers.FindUp(".loccignore")
	if err != nil {
		fmt.Println(helpers.Colorize("Error finding .loccignore file:", helpers.Red), err)
		return
	}
	if path != "" {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(helpers.Colorize("Error opening .loccignore file:", helpers.Red), err)
			return
		}
		lines := bytes.SplitSeq(content, []byte{'\n'})
		for line := range lines {
			lineString := string(bytes.TrimSpace(line))
			// Add include (!) handling
			if lineString != "" && !bytes.HasPrefix(line, []byte("#")) && bytes.HasPrefix(line, []byte("!")) {
				var newIgnore []string
				for _, i := range IGNORE {
					if strings.TrimLeft(lineString, "!") != i {
						newIgnore = append(newIgnore, i)
					}
				}
				IGNORE = newIgnore
			// Add ignore handling
			} else if lineString != "" && !bytes.HasPrefix(line, []byte("#")) {
				IGNORE = append(IGNORE, lineString)
			}
		}
	}
}