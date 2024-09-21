package pkg

import (
	"fmt"
	"strings"
)

type Color string

const (
	cliRed    Color = "\x1b[31m"
	cliYellow Color = "\x1b[33m"
	cliGreen  Color = "\x1b[32m"
	cliBlue   Color = "\x1b[34m"
	cliReset  Color = "\x1b[0m"
)

type State string

const (
	Unknown    State = "unknown"
	Incomplete State = "incomplete"
	Complete   State = "complete"
	Stub       State = "stub"
)

func NewState(content string) State {
	switch content {
	case "complete":
		return Complete
	case "incomplete":
		return Incomplete
	case "stub":
		return Stub
	}

	return Unknown
}

type Badge string

const (
	Alternative Badge = "alternative"
	Extra       Badge = "extra"
	Fun         Badge = "fun"
	Hint        Badge = "hint"
	MustSee     Badge = "must-see"
	Summary     Badge = "summary"
	Unchecked   Badge = "unchecked"
	NoEmbed     Badge = "no-embed"
)

type MainVideo string

const (
	VideoMissing       MainVideo = "missing"
	VideoReallyMissing MainVideo = "really missing"
	VideoPresent       MainVideo = "present"
	VideoProblem       MainVideo = "problem"
)

type RelatedVideo struct {
	Badge   Badge
	Issues  []string
	Minutes int
}

type RelatedVideos []RelatedVideo

func (rv RelatedVideos) GetIssues() []string {
	var issues []string

	for _, item := range rv {
		issues = append(issues, item.Issues...)
	}

	return issues
}

func (rv RelatedVideos) Has(badge Badge) bool {
	for _, item := range rv {
		if item.Badge == badge {
			return true
		}
	}

	return false
}

type DefaultBody struct {
	MainVideo       MainVideo
	HasSummary      bool
	HasTopics       bool
	HasPractice     bool
	RelatedVideos   RelatedVideos
	HasRelatedLinks bool
}

func (db DefaultBody) GetIssues(state State) []string {
	issues := db.RelatedVideos.GetIssues()

	switch db.MainVideo {
	case VideoReallyMissing:
		if db.RelatedVideos.Has(Alternative) {
			issues = append(issues, "main video is NOT REALLY missing")
		}
	case VideoMissing:
		if !db.RelatedVideos.Has(Alternative) {
			issues = append(issues, "main video is REALLY missing")
		}
	}

	if state != db.CalculateState() {
		issues = append(issues, fmt.Sprintf("state mismatch. got: %s, want: %s", state, db.CalculateState()))
	}

	return issues
}

func (db DefaultBody) CalculateState() State {
	if db.MainVideo == VideoPresent && db.HasSummary {
		return Complete
	}

	if db.MainVideo == VideoPresent || db.RelatedVideos.Has(Alternative) {
		return Incomplete
	}

	return Stub
}

func (db DefaultBody) IsIndex() bool {
	return false
}

type ChapterBody struct {
	HasEpisodes bool
}

func (cb ChapterBody) GetIssues(_ State) []string {
	return nil
}

func (cb ChapterBody) CalculateState() State {
	if cb.HasEpisodes {
		return Complete
	}

	return Stub
}

func (cb ChapterBody) IsIndex() bool {
	return true
}

type Body interface {
	GetIssues(state State) []string
	CalculateState() State
	IsIndex() bool
}

type Content struct {
	Title    string
	State    State
	Body     Body
	FilePath string
	Slug     string
	Weight   string
}

func (c Content) GetIssues() []string {
	issues := c.Body.GetIssues(c.State)

	if c.Body.IsIndex() {
		if !strings.HasPrefix(c.FilePath, c.Weight) {
			issues = append(issues, "title is not prefixed with weight: '"+c.FilePath+"', prefix: '"+c.Weight+"'")
		}
	}

	return issues
}

type Page struct {
	Filename string
	Content  Content
}

func (p Page) GetIssues() []string {
	issues := p.Content.GetIssues()

	return issues
}

func (p Page) String() string {
	color := cliBlue

	switch p.Content.State {
	case Complete:
		color = cliGreen
	case Incomplete:
		color = cliYellow
	case Unknown:
		color = cliRed
	}

	issues := p.GetIssues()
	if len(issues) > 0 {
		color = cliRed
	}

	result := fmt.Sprintln("    ", color, p.Filename, "-", p.Content.State, cliReset)

	for _, issue := range issues {
		result += fmt.Sprintln("        - ", issue)
	}

	return result
}

type Pages []Page

func (p Pages) Add(pageFN string, content Content) Pages {
	return append(p, Page{Filename: pageFN, Content: content})
}

type Chapter struct {
	Title string
	Pages Pages
}

func (c Chapter) String() string {
	result := fmt.Sprintln("  ", c.Title)

	for _, page := range c.Pages {
		result += page.String()
	}

	return result
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

func (c Course) String() string {
	result := fmt.Sprintln(c.Title)

	for _, chapter := range c.Chapters {
		result += chapter.String()
	}

	return result
}

type Courses []Course

func (c Courses) Add(courseFN, chapterFN, pageFN string, content Content) Courses {
	for i, course := range c {
		if course.Title == courseFN {
			c[i].Chapters = c[i].Chapters.Add(chapterFN, pageFN, content)
			return c
		}
	}

	return append(c, Course{Title: courseFN, Chapters: Chapters{{Title: chapterFN, Pages: Pages{{Filename: pageFN, Content: content}}}}})
}