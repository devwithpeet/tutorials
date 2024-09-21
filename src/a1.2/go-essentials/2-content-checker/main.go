package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devwithpeet/tutorials/src/a1.2/go-essentials/2-content-checker/pkg"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// collect markdown files
	pattern := filepath.Join(root, "content") + "/**/**/*.md"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error fetching markdown files:", err)
		os.Exit(1)
	}

	// fetch markdown files
	courses, count := CrawlMarkdownFiles(matches)

	Print(count, courses)
}

const maxErrors = 3

func CrawlMarkdownFiles(matches []string) (pkg.Courses, int) {
	result := make(pkg.Courses, 0, len(matches))

	var count, errCount int

	for _, filePath := range matches {
		if errCount >= maxErrors {
			break
		}

		parts := strings.Split(filePath, "/")

		if len(parts) < 3 {
			fmt.Println("Skipping:", filePath)
			continue
		}

		course := parts[len(parts)-3]
		chapter := parts[len(parts)-2]
		page := parts[len(parts)-1]

		rawContent, err := os.ReadFile(filePath)
		if err != nil {
			panic("cannot open file: " + filePath)
		}

		content := pkg.ParseMarkdown(string(rawContent), filePath)
		result = result.Add(course, chapter, page, content)

		if len(content.GetIssues()) > 0 {
			errCount++
		}

		count++
	}

	return result, count
}

func Print(count int, courses pkg.Courses) {
	fmt.Println("Processed", count, "markdown files")

	for _, course := range courses {
		fmt.Print(course.String())
	}
}
