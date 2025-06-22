package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/qwertycodeqc/locc/src/helpers"

	"github.com/bmatcuk/doublestar/v4"
)

const VERSION = "1.0"

var IGNORE = []string{"**/.git/**", "**/node_modules/**", "**/dist/**", "**/coverage/**", "README.md", "LICENSE"}
var BINARY_EXTENSIONS = []string{"**/*.exe", "**/*.dll", "**/*.so", "**/*.dylib", "**/*.o", "**/*.a", "**/*.lib", "**/*.class", "**/*.jar", "**/*.pyc", "**/*.pyo"}

func must(v any, err error) any {
	if err != nil {
		fmt.Println(helpers.Colorize("Error:", helpers.Red), err)
		os.Exit(1)
	}
	return v
}

func main() {
	fmt.Println(helpers.Colorize("locc", helpers.Cyan, helpers.Bold), "v"+VERSION)
	IGNORE = append(IGNORE, BINARY_EXTENSIONS...)
	loadIgnore()
	var totalLines int
	err := filepath.Walk(must(os.Getwd()).(string), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && shouldIgnore(path) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if shouldIgnore(path) {
			return filepath.SkipDir
		}

		lines, err := countLines(path)
		if err != nil {
			fmt.Println(helpers.Colorize("Error counting lines in file:", helpers.Red), path, err)
			return nil
		}
		totalLines += lines
		relativePath, err := filepath.Rel(must(os.Getwd()).(string), path)
		if err != nil {
			fmt.Println(helpers.Colorize("Error getting relative path:", helpers.Red), err)
			return nil
		}
		relativePath = filepath.ToSlash(relativePath)
		fmt.Println(
			helpers.Colorize("▪︎ ", helpers.Yellow),
			helpers.Colorize(relativePath, helpers.Blue),
			helpers.Colorize(fmt.Sprint(lines), helpers.Yellow),
		)
		return nil
	})
	if err != nil {
		fmt.Println(helpers.Colorize("Error walking the path:", helpers.Red), err)
		return
	}
	fmt.Println(helpers.Colorize(fmt.Sprintf(" TOTAL %d ", totalLines), helpers.YellowBg, helpers.Black))
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

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close() //nolint

	var count int
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func shouldIgnore(path string) bool {
	norm := filepath.ToSlash(path)

	for _, pattern := range IGNORE {
		pattern = filepath.ToSlash(pattern)
		matched, err := doublestar.Match(pattern, norm)
		if err != nil {
			fmt.Println("Invalid pattern:", pattern, err)
			continue
		}
		if matched {
			return true
		}
	}

	return false
}
