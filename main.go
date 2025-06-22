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

var IGNORE = []string{".git", "**/node_modules/**", "**/dist/**", "**/coverage/**"}
var BINARY_EXTENSIONS = []string{"**/*.exe", "**/*.dll", "**/*.so", "**/*.dylib", "**/*.o", "**/*.a", "**/*.lib", "**/*.class", "**/*.jar", "**/*.pyc", "**/*.pyo", "**/*.jpg", "**/*.png", "**/*.gif", "**/*.svg", "**/*.ico", "**/*.webp", "**/*.mp4", "**/*.mp3", "**/*.wav", "**/*.flac", "**/*.avi", "**/*.mkv", "**/*.mov", "**/*.wmv", "**/*.zip", "**/*.tar.gz", "**/*.tar.bz2", "**/*.rar"}

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
	fmt.Println(helpers.Colorize("Ignoring patterns:", helpers.Yellow))
	for _, pattern := range IGNORE {
		fmt.Print("  ", helpers.Colorize(pattern, helpers.Yellow))
	}
	fmt.Println()
	var totalLines int
	err := filepath.Walk(must(os.Getwd()).(string), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if shouldIgnore(path) {
			return nil
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

	rel, err := filepath.Rel(must(os.Getwd()).(string), path)
	if err != nil {
		rel = norm
	} else {
		rel = filepath.ToSlash(rel)
	}

	for _, pattern := range IGNORE {
		pattern = filepath.ToSlash(pattern)

		// Match against:
		// - base name (e.g. go.mod)
		// - relative path (e.g. .git/hooks/xyz)
		// - full normalized path
		candidates := []string{
			filepath.Base(norm),
			rel,
			norm,
		}

		for _, candidate := range candidates {
			matched, err := doublestar.Match(pattern, candidate)
			if err != nil {
				fmt.Println("Bad pattern:", pattern, err)
				continue
			}
			if matched {
				//fmt.Println(helpers.Colorize("Ignoring:", helpers.Yellow), helpers.Colorize(path, helpers.Blue), "due to pattern", helpers.Colorize(pattern, helpers.Yellow))
				return true
			}
		}
	}
	return false
}
