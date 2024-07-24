package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/devwithpeet/tutorials/src/a1.2/1-go-essentials/2-content-checker/pkg"
)

func main() {
	// editor to use when opening a markdown file
	editor := "open"
	if len(os.Args) > 1 {
		editor = os.Args[1]
	}

	// collect markdown files
	matches, err := filepath.Glob("./content/**/**/*.md")
	if err != nil {
		fmt.Println("Error fetching markdown files:", err)
		os.Exit(1)
	}

	// fetch markdown files
	courses, count := CrawlMarkdownFiles(matches)

	// convert the fetched raw data into rows for table display
	rows := courses.ToRows(count)

	// start CLI UI
	m := pkg.NewModel(rows, editor)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func CrawlMarkdownFiles(matches []string) (pkg.Courses, int) {
	result := make(pkg.Courses, 0, len(matches))

	for _, filePath := range matches {
		parts := strings.Split(filePath, "/")

		if len(parts) < 4 {
			fmt.Println("Skipping:", filePath)
			continue
		}

		course := parts[1]
		chapter := parts[2]
		page := parts[3]

		rawContent, err := os.ReadFile(filePath)
		if err != nil {
			panic("cannot open file: " + filePath)
		}

		content := pkg.ParseMarkdown(string(rawContent), filePath)

		result = result.Add(course, chapter, page, content)
	}

	return result, len(matches)
}
