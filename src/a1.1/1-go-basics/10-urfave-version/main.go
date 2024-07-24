package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

const defaultCount = 10

type wordOccurrences struct {
	word  string
	files []string
}

type byOccurrence []wordOccurrences

func (b byOccurrence) Len() int           { return len(b) }
func (b byOccurrence) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byOccurrence) Less(i, j int) bool { return len(b[i].files) > len(b[j].files) }

func main() {
	app := &cli.App{
		Name:  "filecat",
		Usage: "display file categories",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-num",
				Usage: "ignore words starting with a number",
			},
			&cli.BoolFlag{
				Name:  "lower-alpha-only",
				Usage: "ignore words with non-lower alpha characters",
			},
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"n"},
				Value:   defaultCount,
				Usage:   "number of words to display (top occurrences)",
			},
			&cli.IntFlag{
				Name:    "word-from",
				Aliases: []string{"f"},
				Value:   0,
				Usage:   "ignore words before a certain count",
			},
			&cli.IntFlag{
				Name:    "word-to",
				Aliases: []string{"t"},
				Value:   100,
				Usage:   "ignore words after a certain count",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "exclude specific words from being counted",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "verbose output",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "categories",
				Aliases: []string{"c"},
				Usage:   "list categories",
				Action: func(cCtx *cli.Context) error {
					root := cCtx.Args().Get(0)

					noNum := cCtx.Bool("no-num")
					loweAlphaOnly := cCtx.Bool("lower-alpha-only")
					wordFrom := cCtx.Int("word-from")
					wordTo := cCtx.Int("word-to")
					excludes := cCtx.StringSlice("exclude")
					maxCount := cCtx.Int("n")
					verbose := cCtx.Bool("verbose")

					files := getFiles(root, verbose)
					stats := getStats(files, noNum, loweAlphaOnly, wordFrom, wordTo, excludes, verbose)
					sortedStats := getSortedStats(stats)
					writeStats(sortedStats, maxCount, verbose)

					return nil
				},
			},
			{
				Name:    "files",
				Aliases: []string{"f"},
				Usage:   "list files",
				Action: func(cCtx *cli.Context) error {
					root := cCtx.Args().First()
					categories := cCtx.Args().Tail()

					noNum := cCtx.Bool("no-num")
					loweAlphaOnly := cCtx.Bool("lower-alpha-only")
					wordFrom := cCtx.Int("word-from")
					wordTo := cCtx.Int("word-to")
					excludes := cCtx.StringSlice("exclude")
					maxCount := cCtx.Int("n")
					verbose := cCtx.Bool("verbose")

					files := getFiles(root, verbose)
					stats := getStats(files, noNum, loweAlphaOnly, wordFrom, wordTo, excludes, verbose)
					sortedStats := getSortedStats(stats)
					writeFiles(sortedStats, categories, maxCount, verbose)

					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "display version",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("0.0.1")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getFiles(root string, verbose bool) map[string]os.FileInfo {
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

var fileExtension = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)

func getStats(files map[string]os.FileInfo, noNum, lowerAlphaOnly bool, wordFrom, wordTo int, excludes []string, verbose bool) map[string]wordOccurrences {
	if verbose {
		fmt.Println("noNum:", noNum)
		fmt.Println("lowerAlphaOnly:", lowerAlphaOnly)
		fmt.Println("wordFrom:", wordFrom)
		fmt.Println("wordTo:", wordTo)
		fmt.Println("skip:", strings.Join(excludes, ", "))
	}

	stats := make(map[string]wordOccurrences)
	for path, pathInfo := range files {
		fileName := fileExtension.ReplaceAllString(pathInfo.Name(), "")

		words := strings.Split(fileName, "-")
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
			for _, exclude := range excludes {
				if exclude == word {
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

func writeStats(sortedStats byOccurrence, maxCount int, verbose bool) {
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

func writeFiles(sortedStats byOccurrence, categories []string, maxCount int, verbose bool) {
	if verbose {
		fmt.Println("n:", maxCount)
	}

	for _, wordStats := range sortedStats {
		found := true
		for _, category := range categories {
			if !strings.Contains(wordStats.word, category) {
				found = false
				break
			}
		}
		if !found {
			continue
		}

		fmt.Println(wordStats.word)
		fmt.Println(strings.Repeat("-", len(wordStats.word)))
		for i, path := range wordStats.files {
			if i >= maxCount {
				break
			}

			fmt.Println(path)
		}
	}
}
