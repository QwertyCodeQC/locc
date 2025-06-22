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



func must(v any, err error) any {
	if err != nil {
		fmt.Println(helpers.Colorize("Error:", helpers.Red), err)
		os.Exit(1)
	}
	return v
}

func main() {
	fmt.Println(helpers.Colorize("locc", helpers.Cyan, helpers.Bold), "v"+VERSION)
	IGNORE = append(IGNORE, helpers.BINARY_EXTENSIONS...)
	loadIgnore()

	var totalLines int
	var totalCommentCount int
	var totalBlankCount int
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

		commentMark := helpers.CommentMarkers[filepath.Ext(path)]
		if commentMark == "" {
			commentMark = "//"
		}

		lines, commentCount, blankCount, err := countLines(path, commentMark)
		if err != nil {
			fmt.Println(helpers.Colorize("Error counting lines in file:", helpers.Red), path, err)
			return nil
		}
		totalLines += lines
		totalCommentCount += commentCount
		totalBlankCount += blankCount
		relativePath, err := filepath.Rel(must(os.Getwd()).(string), path)
		if err != nil {
			fmt.Println(helpers.Colorize("Error getting relative path:", helpers.Red), err)
			return nil
		}
		relativePath = filepath.ToSlash(relativePath)

		lang := helpers.LanguageMeta{
			Lang:  "Unknown",
			Color: helpers.Gray,
		}
		ext := filepath.Ext(relativePath)
		for _, l := range helpers.Languages {
			if l.Ext == ext {
				lang = l
				break
			}
		}

		fmt.Println(
			helpers.Colorize("▪︎ ", helpers.Yellow),
			helpers.Colorize(lang.Lang, lang.Color),
			helpers.Colorize(relativePath, helpers.Blue),
			helpers.Colorize(fmt.Sprint(lines), helpers.Yellow),
			helpers.Colorize(fmt.Sprint(commentCount), helpers.Green),
			helpers.Colorize(fmt.Sprint(blankCount), helpers.Gray),
		)
		return nil
	})
	if err != nil {
		fmt.Println(helpers.Colorize("Error walking the path:", helpers.Red), err)
		return
	}
	codeLines := totalLines - totalCommentCount - totalBlankCount
	var codePercent float64
	if totalLines > 0 {
		codePercent = float64(codeLines) / float64(totalLines) * 100.0
	}
	commentsString := fmt.Sprintf(" COMMENTS %d ", totalCommentCount)
	totalString := fmt.Sprintf(" TOTAL %d ", totalLines)
	blanksString := fmt.Sprintf(" BLANKS %d ", totalBlankCount)
	codeLinesString := fmt.Sprintf(" CODE %d (%.2f%%) ", codeLines, codePercent)

	longestLength := max(len(totalString), len(commentsString), len(blanksString), len(codeLinesString))
	totalString = fmt.Sprintf("%-*s", longestLength, totalString)
	commentsString = fmt.Sprintf("%-*s", longestLength, commentsString)
	blanksString = fmt.Sprintf("%-*s", longestLength, blanksString)
	codeLinesString = fmt.Sprintf("%-*s", longestLength, codeLinesString)

	fmt.Println()
	fmt.Println(helpers.Colorize(totalString, helpers.YellowBg, helpers.Black))
	fmt.Println(helpers.Colorize(commentsString, helpers.GreenBg, helpers.Black))
	fmt.Println(helpers.Colorize(blanksString, helpers.GrayBg, helpers.Black))
	fmt.Println(helpers.Colorize(codeLinesString, helpers.CyanBg, helpers.Black))
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

// Returns count, commentCount, blankCount, and error.
func countLines(path string, comment string) (int, int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close() //nolint

	var count int
	var commentCount int
	var blankCount int
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			blankCount++
		}
		if strings.HasPrefix(line, comment) {
			commentCount++
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, 0, err
	}

	return count, commentCount, blankCount, nil
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
