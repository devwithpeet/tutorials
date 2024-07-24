package pkg

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
)

type State string

const (
	unknown    State = "unknown"
	incomplete State = "incomplete"
	complete   State = "complete"
	stub       State = "stub"
)

type DefaultBody struct {
	HasMainVideo     bool
	HasSummary       bool
	HasTopics        bool
	HasPractice      bool
	HasRelatedVideos bool
	HasRelatedLinks  bool
}

func (db DefaultBody) ToRow() [6]string {
	return [6]string{
		formatBoolean(db.HasMainVideo),
		formatBoolean(db.HasSummary),
		formatBoolean(db.HasTopics),
		formatBoolean(db.HasRelatedVideos),
		formatBoolean(db.HasRelatedLinks),
		formatBoolean(db.HasPractice),
	}
}

type ChapterBody struct {
	HasEpisodes bool
}

func (cb ChapterBody) ToRow() [6]string {
	return [6]string{
		formatBoolean(cb.HasEpisodes),
		"",
		"",
		"",
		"",
		"",
	}
}

type Body interface {
	ToRow() [6]string
}

type Content struct {
	Title           string
	State           State
	CalculatedState State
	Body            Body
	FilePath        string
}

func (c Content) ToRow() [11]string {
	bodyRow := c.Body.ToRow()

	if _, ok := c.Body.(ChapterBody); ok {
		return [11]string{
			c.FilePath,
			c.Title,
			bodyRow[0],
			bodyRow[1],
			bodyRow[2],
			bodyRow[3],
			bodyRow[4],
			bodyRow[5],
			"",
			"",
			bodyRow[0],
		}

	}

	return [11]string{
		c.FilePath,
		c.Title,
		bodyRow[0],
		bodyRow[1],
		bodyRow[2],
		bodyRow[3],
		bodyRow[4],
		bodyRow[5],
		string(c.State),
		string(c.CalculatedState),
		formatBoolean(c.CalculatedState == c.State),
	}
}

type Page struct {
	Filename string
	Content  Content
}

type Pages []Page

func (p Pages) Add(pageFN string, content Content) Pages {
	return append(p, Page{Filename: pageFN, Content: content})
}

type Chapter struct {
	Title string
	Pages Pages
}

type Chapters []Chapter

func (c Chapters) Add(chapterFN, pageFN string, content Content) Chapters {
	for i, chapter := range c {
		if chapter.Title == chapterFN {
			c[i].Pages = c[i].Pages.Add(pageFN, content)
			return c
		}
	}

	return append(c, Chapter{Title: chapterFN, Pages: Pages{{Filename: pageFN, Content: content}}})
}

type Course struct {
	Title    string
	Chapters Chapters
}

type Courses []Course

func (c Courses) ToRows(count int) []table.Row {
	result := make([]table.Row, 0, count)

	row := 1
	for _, course := range c {
		for _, chapter := range course.Chapters {
			for _, page := range chapter.Pages {
				content := page.Content

				bodyRow := content.ToRow()

				result = append(result, table.Row{
					bodyRow[0],
					fmt.Sprintf("%d", row),
					course.Title,
					chapter.Title,
					bodyRow[1],
					bodyRow[2],
					bodyRow[3],
					bodyRow[4],
					bodyRow[5],
					bodyRow[6],
					bodyRow[7],
					bodyRow[8],
					bodyRow[9],
					bodyRow[10],
				})

				row++
			}
		}
	}

	return result
}

func formatBoolean(b bool) string {
	if b {
		return "Y"
	}

	return "N"
}

func (c Courses) Add(courseFN, chapterFN, pageFN string, content Content) Courses {
	for i, course := range c {
		if course.Title == courseFN {
			c[i].Chapters = c[i].Chapters.Add(chapterFN, pageFN, content)
			return c
		}
	}

	return append(c, Course{Title: courseFN, Chapters: Chapters{{Title: chapterFN, Pages: Pages{{Filename: pageFN, Content: content}}}}})
}
