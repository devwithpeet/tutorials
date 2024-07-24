package pkg

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os/exec"
)

func NewModel(rows []table.Row, editor string) model {
	columns := []table.Column{
		{Title: "", Width: 0},
		{Title: "#", Width: 3},
		{Title: "Course", Width: 6},
		{Title: "Chapter", Width: 20},
		{Title: "Page", Width: 20},
		{Title: "Main", Width: 4},
		{Title: "Summary", Width: 7},
		{Title: "Topics", Width: 6},
		{Title: "Videos", Width: 6},
		{Title: "Links", Width: 5},
		{Title: "Practice", Width: 8},
		{Title: "State", Width: 10},
		{Title: "Calculated", Width: 10},
		{Title: "OK", Width: 2},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{table: t, editor: editor}
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table  table.Model
	editor string
}

func (m model) Init() tea.Cmd {
	return nil
}

type execFinishedMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			ex := exec.Command(m.editor, m.table.SelectedRow()[0])
			return m, tea.Batch(
				tea.ExecProcess(ex, func(err error) tea.Msg {
					return execFinishedMsg{err}
				}),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + EOL
}
