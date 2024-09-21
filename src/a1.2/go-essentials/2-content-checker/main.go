package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/devwithpeet/tutorials/src/a1.2/go-essentials/2-content-checker/pkg"
)

type Command string

const (
	PrintCommand  Command = "print"
	ErrorsCommand Command = "errors"
	StatsCommand  Command = "stats"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	action := PrintCommand
	if len(os.Args) > 2 {
		action = Command(os.Args[2])
	}

	// collect markdown files
	files, err := findFiles(root)
	if err != nil {
		panic("cannot find files in root: " + root + ", error: " + err.Error())
	}

	// fetch markdown files
	courses, count := CrawlMarkdownFiles(files, maxErrors)

	Prepare(courses)

	switch action {
	case PrintCommand:
		Print(count, courses)

	case ErrorsCommand:
		Errors(count, courses)

	case StatsCommand:
		Stats(courses)

	default:
		panic("unknown command: " + string(action))
	}
}

func findFiles(root string) ([]string, error) {
	pattern := filepath.Join(root, "content") + "/**/**/*.md"

	return filepath.Glob(pattern)
}

const maxErrors = 3

func CrawlMarkdownFiles(matches []string, maxErrors int) (pkg.Courses, int) {
	if maxErrors < 0 {
		maxErrors = math.MaxInt
	}

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

		content, err := pkg.ParseMarkdown(string(rawContent))
		if err != nil {
			panic("cannot parse markdown: " + filePath + ", err: " + err.Error())
		}

		result = result.Add(filePath, course, chapter, page, content)

		if len(content.GetIssues(filePath)) > 0 {
			errCount++
		}

		count++
	}

	return result, count
}

func Prepare(courses pkg.Courses) {
	for _, course := range courses {
		course.Prepare()
	}
}

func Print(count int, courses pkg.Courses) {
	fmt.Println("Processed", count, "markdown files")

	for _, course := range courses {
		fmt.Print(course.String())
	}
}

func Errors(count int, courses pkg.Courses) {
	fmt.Println("Processed", count, "markdown files")

	errorsFound := false

	for _, course := range courses {
		errors := course.GetErrors()
		if len(errors) == 0 {
			continue
		}

		errorsFound = true

		fmt.Println(strings.Join(errors, "\n"))
	}

	if errorsFound {
		os.Exit(1)
	}
}

func Stats(courses pkg.Courses) {
	fmt.Println("Stats")

	var total, totalStub, totalIncomplete, totalComplete, totalErrors int

	for _, course := range courses {
		courseAll, courseStub, courseIncomplete, courseComplete, courseErrors := course.Stats()

		fmt.Println("Course:", course.Title)
		fmt.Println("  Total:", courseAll)
		fmt.Println("  Stub:", courseStub)
		fmt.Println("  Incomplete:", courseIncomplete)
		fmt.Println("  Complete:", courseComplete)
		fmt.Println("  Errors:", courseErrors)
		fmt.Println()

		total += courseAll
		totalStub += courseStub
		totalIncomplete += courseIncomplete
		totalComplete += courseComplete
		totalErrors += courseErrors
	}

	fmt.Println("Total:", total)
	fmt.Println("  Stub:", totalStub)
	fmt.Println("  Incomplete:", totalIncomplete)
	fmt.Println("  Complete:", totalComplete)
	fmt.Println("  Errors:", totalErrors)

}
