package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		Action: func(cCtx *cli.Context) error {
			root := cCtx.Args().Get(0)
			noNum := cCtx.Bool("no-num")
			lowerAlphaOnly := cCtx.Bool("lower-alpha-only")
			wordFrom := cCtx.Int("word-from")
			wordTo := cCtx.Int("word-to")
			excludes := cCtx.StringSlice("exclude")
			maxCount := cCtx.Int("n")
			verbose := cCtx.Bool("verbose")

			p := tea.NewProgram(initialModel(root, noNum, lowerAlphaOnly, wordFrom, wordTo, excludes, maxCount, verbose))
			if _, err := p.Run(); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type gotOptionSuccessMsg byOccurrence

func (m model) getCategories() tea.Msg {
	files := getFiles(m.root, m.verbose)
	stats := getStats(files, m.noNum, m.lowerAlphaOnly, m.wordFrom, m.wordTo, m.excludes, m.verbose)
	sortedStats := getSortedStats(stats)
	// writeStats(sortedStats, maxCount, verbose)

	options := make([]string, 0, len(sortedStats))
	for _, wordStats := range sortedStats {
		options = append(options, wordStats.word)
	}

	return gotOptionSuccessMsg(sortedStats)
}

type model struct {
	textInput textinput.Model
	help      help.Model
	keymap    keymap

	root           string
	noNum          bool
	lowerAlphaOnly bool
	wordFrom       int
	wordTo         int
	excludes       []string
	maxCount       int
	verbose        bool

	options  []wordOccurrences
	selected *wordOccurrences
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.getCategories, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			for _, r := range m.options {
				if r.word == m.textInput.Value() {
					m.selected = &r

					break
				}
			}

			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case gotOptionSuccessMsg:
		m.options = msg
		var suggestions []string
		for _, r := range msg {
			suggestions = append(suggestions, r.word)
		}
		m.textInput.SetSuggestions(suggestions)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	sb := strings.Builder{}
	if m.selected != nil {
		sb.WriteString("Files:\n")
		for i, file := range m.selected.files {
			if i >= m.maxCount {
				break
			}

			sb.WriteString(fmt.Sprintf("%s\n", file))
		}
		sb.WriteString("\n")
	}

	return fmt.Sprintf(
		"Pick a category:\n\n  %s\n\n%s\n%s",
		m.textInput.View(),
		m.help.View(m.keymap),
		sb.String(),
	)
}

func initialModel(root string, noNum, lowerAlphaOnly bool, wordFrom, wordTo int, excludes []string, maxCount int, verbose bool) model {
	ti := textinput.New()
	ti.Placeholder = "category"
	// ti.Prompt = ""
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.ShowSuggestions = true

	h := help.New()

	km := keymap{}

	return model{
		textInput:      ti,
		help:           h,
		keymap:         km,
		root:           root,
		noNum:          noNum,
		lowerAlphaOnly: lowerAlphaOnly,
		wordFrom:       wordFrom,
		wordTo:         wordTo,
		excludes:       excludes,
		maxCount:       maxCount,
		verbose:        verbose,
	}
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
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
