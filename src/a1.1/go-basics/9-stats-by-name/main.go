package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type wordOccurrences struct {
	word  string
	files []string
}

type byOccurrence []wordOccurrences

func (b byOccurrence) Len() int           { return len(b) }
func (b byOccurrence) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byOccurrence) Less(i, j int) bool { return len(b[i].files) > len(b[j].files) }

func main() {
	files := getFiles()
	stats := getStats(files)
	sortedStats := getSortedStats(stats)
	writeStats(sortedStats)
}

func getFiles() map[string]os.FileInfo {
	root := getRoot()
	verbose := getVerbose()

	if verbose {
		fmt.Println("root:", root)
	}

	pathInfos := make(map[string]os.FileInfo)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		pathInfos[path] = info

		return nil
	})
	if err != nil {
		panic(err)
	}

	return pathInfos
}

func getStats(files map[string]os.FileInfo) map[string]wordOccurrences {
	noNum := getNoNum()
	lowerAlphaOnly := getLowerAlphaOnly()
	wordFrom := getWordFrom()
	wordTo := getWordTo()
	skips := getSkips()
	verbose := getVerbose()

	if verbose {
		fmt.Println("noNum:", noNum)
		fmt.Println("lowerAlphaOnly:", lowerAlphaOnly)
		fmt.Println("wordFrom:", wordFrom)
		fmt.Println("wordTo:", wordTo)
		fmt.Println("skip:", strings.Join(skips, ", "))
	}

	stats := make(map[string]wordOccurrences)
	for path, pathInfo := range files {
		words := strings.Split(pathInfo.Name(), "-")

		for i, word := range words {
			if word == "" {
				continue
			}

			if noNum && word[0] >= '0' && word[0] <= '9' {
				continue
			}

			if lowerAlphaOnly {
				nonLowerAlphaFound := false
				for _, c := range word {
					if c < 'a' || c > 'z' {
						nonLowerAlphaFound = true
						break
					}
				}

				if nonLowerAlphaFound {
					continue
				}
			}

			if i < wordFrom || i > wordTo {
				continue
			}

			skipFound := false
			for _, skip := range skips {
				if skip == word {
					skipFound = true
					break
				}
			}
			if skipFound {
				continue
			}

			wordStats, ok := stats[word]
			if !ok {
				wordStats = wordOccurrences{
					word:  word,
					files: []string{},
				}
			}

			wordStats.files = append(wordStats.files, path)

			stats[word] = wordStats
		}
	}

	return stats
}

func getSortedStats(stats map[string]wordOccurrences) byOccurrence {
	sortedStats := byOccurrence{}
	for _, wordStats := range stats {
		sortedStats = append(sortedStats, wordStats)
	}

	sort.Sort(sortedStats)

	return sortedStats
}

func writeStats(sortedStats byOccurrence) {
	maxCount := getCount()
	verbose := getVerbose()

	if verbose {
		fmt.Println("n:", maxCount)
	}

	for i, wordStats := range sortedStats {
		if i >= maxCount {
			break
		}

		fmt.Printf("%s (%d)\n", wordStats.word, len(wordStats.files))
	}
}

func getRoot() string {
	if len(os.Args) < 2 {
		return ".*"
	}

	return os.Args[1]
}

func getLowerAlphaOnly() bool {
	for i, arg := range os.Args {
		if i < 2 {
			continue
		}

		if arg == "--lower-alpha-only" {
			return true
		}
	}

	return false
}

func getNoNum() bool {
	for i, arg := range os.Args {
		if i < 2 {
			continue
		}

		if arg == "--no-num" {
			return true
		}
	}

	return false
}

func getWordTo() int {
	for i, arg := range os.Args {
		if i < 2 || len(arg) < len("--word-to=") {
			continue
		}

		if arg[:10] == "--word-to=" {
			num, err := strconv.Atoi(arg[10:])
			if err != nil {
				panic(err)
			}

			return num
		}
	}

	return 100
}

func getWordFrom() int {
	for i, arg := range os.Args {
		if i < 2 || len(arg) < len("--word-from") {
			continue
		}

		if arg[:12] == "--word-from=" {
			num, err := strconv.Atoi(arg[12:])
			if err != nil {
				panic(err)
			}

			return num
		}
	}

	return 0
}

const defaultCount = 10

func getCount() int {
	for i, arg := range os.Args {
		if i < 2 {
			continue
		}

		if arg[:4] == "--n=" {
			num, err := strconv.Atoi(arg[4:])
			if err != nil {
				panic(err)
			}

			return num
		}

	}

	return defaultCount
}

func getVerbose() bool {
	for i, arg := range os.Args {
		if i < 2 {
			continue
		}

		if arg == "--verbose" {
			return true
		}
	}

	return false
}

func getSkips() []string {
	skip := []string{}
	for i, arg := range os.Args {
		if i < 2 || len(arg) < len("--skip=") {
			continue
		}

		if arg[:7] == "--skip=" {
			skip = append(skip, arg[7:])
		}
	}

	return skip
}
