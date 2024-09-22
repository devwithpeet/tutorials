package pkg

import (
	"fmt"
	"path/filepath"
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

	return ""
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

type Audience string

const (
	All              Audience = "all"
	AllProfessionals Audience = "all professionals"
	LinuxUsers       Audience = "Linux users"
	WindowsUsers     Audience = "Windows users"
	MacUsers         Audience = "Mac users"
	AllDevelopers    Audience = "all developers"
	WebDevelopers    Audience = "web developers"
	MobileDevelopers Audience = "mobile developers"
	GameDevelopers   Audience = "game developers"
	SysAdmins        Audience = "sysadmins"
)

var validAudiences = map[Audience]struct{}{
	All:              {},
	AllProfessionals: {},
	LinuxUsers:       {},
	WindowsUsers:     {},
	MacUsers:         {},
	AllDevelopers:    {},
	WebDevelopers:    {},
	MobileDevelopers: {},
	GameDevelopers:   {},
	SysAdmins:        {},
}

type Importance string

const (
	Critical   Importance = "critical"
	Essential  Importance = "essential"
	Important  Importance = "important"
	Relevant   Importance = "relevant"
	Optional   Importance = "optional"
	Irrelevant Importance = "irrelevant"
)

func (i Importance) Level() int {
	switch i {
	case Critical:
		return 5
	case Essential:
		return 4
	case Important:
		return 3
	case Relevant:
		return 2
	case Optional:
		return 1
	case Irrelevant:
		return 0
	}

	return -1
}

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
	Valid   bool
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
	MainVideo          MainVideo
	HasSummary         bool
	HasTopics          bool
	HasExercises       bool
	RelatedVideos      RelatedVideos
	HasRelatedLinks    bool
	UsefulWithoutVideo bool
	Titles             []string
}

var defaultBodySectionMap = map[string]int{
	sectionMainVideo:       0,
	sectionSummary:         1,
	sectionTopics:          2,
	sectionRelatedVideos:   3,
	sectionRelatedArticles: 4,
	sectionRelatedLinks:    5,
	sectionExercises:       6,
}

func isOrderedCorrectly(goldenMap map[string]int, givenSlice []string) (string, bool) {
	lastIndex := -1
	for _, item := range givenSlice {
		if index, exists := goldenMap[item]; exists {
			if index < lastIndex {
				return item, false
			}

			lastIndex = index
		}
	}

	return "", true
}

func (db DefaultBody) GetIssues(state State) []string {
	issues := db.RelatedVideos.GetIssues()

	switch db.MainVideo {
	case VideoReallyMissing:
		if db.RelatedVideos.Has(Alternative) || db.UsefulWithoutVideo {
			issues = append(issues, "main video is NOT REALLY missing")
		}
	case VideoMissing:
		if !db.RelatedVideos.Has(Alternative) && !db.UsefulWithoutVideo {
			issues = append(issues, "main video is REALLY missing")
		}
	}

	if state != db.CalculateState() {
		issues = append(issues, fmt.Sprintf("state mismatch. got: %s, want: %s", state, db.CalculateState()))
	}

	if item, ok := isOrderedCorrectly(defaultBodySectionMap, db.Titles); !ok {
		issues = append(issues, "sections are not in the correct order, first out of order: "+item)
	}

	return issues
}

func (db DefaultBody) CalculateState() State {
	if db.MainVideo == VideoPresent && db.HasSummary && db.HasExercises {
		return Complete
	}

	if db.MainVideo == VideoPresent || db.RelatedVideos.Has(Alternative) || db.UsefulWithoutVideo {
		return Incomplete
	}

	return Stub
}

func (db DefaultBody) IsIndex() bool {
	return false
}

type IndexBody struct {
	HasEpisodes   bool
	CompleteState State
}

func (ib *IndexBody) GetIssues(_ State) []string {
	return nil
}

func (ib *IndexBody) CalculateState() State {
	if ib.HasEpisodes {
		return ib.CompleteState
	}

	return Stub
}

func (ib *IndexBody) SetCompleteState(state State) {
	ib.CompleteState = state
}

type PracticeBody struct {
	HasDescription           bool
	HasRecommendedChallenges bool
	HasAdditionalChallenges  bool
}

func (pb PracticeBody) GetIssues(_ State) []string {
	return nil
}

func (pb PracticeBody) CalculateState() State {
	if !pb.HasDescription {
		return Stub
	}

	if pb.HasRecommendedChallenges && pb.HasAdditionalChallenges {
		return Complete
	}

	return Incomplete
}

type Body interface {
	GetIssues(state State) []string
	CalculateState() State
}

type Content struct {
	Title             string
	State             State
	Body              Body
	Slug              string
	Weight            string
	Audience          Audience
	Importance        Importance
	OutsideImportance Importance
	Tags              []string
}

func (c Content) GetIssues(filePath string) []string {
	issues := c.Body.GetIssues(c.State)

	_, isDefaultBody := c.Body.(*DefaultBody)
	if isDefaultBody {
		filename := filepath.Base(filePath)
		if !strings.HasPrefix(filename, c.Weight) {
			issues = append(issues, "file name is not prefixed with weight: '"+filename+"', prefix: '"+c.Weight+"'")
		}
	}

	if _, exists := validAudiences[c.Audience]; !exists {
		issues = append(issues, "invalid audience: "+string(c.Audience))
	}

	if c.Importance.Level() < c.OutsideImportance.Level() {
		issues = append(issues, "importance is lower than outside importance")
	}

	if c.OutsideImportance == "" && c.Audience != All {
		issues = append(issues, "outside importance is invalid")
	}

	if c.Audience == All && c.OutsideImportance != "" {
		issues = append(issues, "audience is 'all', outside importance must be empty")
	}

	for _, tag := range c.Tags {
		if tag == "unsorted" {
			issues = append(issues, "tag is 'unsorted'")
		}
		if strings.ToLower(tag) != tag {
			issues = append(issues, "tag is not lowercase: "+tag)
		}
		if strings.Replace(tag, " ", "", 1) != tag {
			issues = append(issues, "tag contains spaces: "+tag)
		}
	}

	return issues
}

type Page struct {
	FilePath string
	Title    string
	Content  Content
}

func (p Page) GetIssues() []string {
	issues := p.Content.GetIssues(p.FilePath)

	return issues
}

func (p Page) GetErrors() []string {
	var errors []string

	for _, issue := range p.GetIssues() {
		errors = append(errors, fmt.Sprintf("%s - %s", p.FilePath, issue))
	}

	return errors
}

func (p Page) GetState() State {
	return p.Content.State
}

func (p Page) String() string {
	color := cliBlue

	switch p.GetState() {
	case Complete:
		color = cliGreen
	case Incomplete:
		color = cliYellow
	default:
		color = cliRed
	}

	issues := p.GetIssues()
	if len(issues) > 0 {
		color = cliRed
	}

	result := fmt.Sprintln("    ", color, p.FilePath, "-", p.Content.State, cliReset)

	for _, issue := range issues {
		result += fmt.Sprintln("        - ", issue)
	}

	return result
}

type Pages []Page

func (p Pages) Add(filePath, pageFN string, content Content) Pages {
	return append(p, Page{FilePath: filePath, Title: pageFN, Content: content})
}

type Chapter struct {
	Title    string
	Pages    Pages
	prepared bool
}

func (c *Chapter) Prepare() {
	if c.prepared {
		return
	}

	c.prepared = true

	var (
		indexPage  *IndexBody
		pagesExist = false
		incomplete = false
	)

	for _, page := range c.Pages {
		chapter, ok := page.Content.Body.(*IndexBody)
		if ok {
			indexPage = chapter

			continue
		}

		pagesExist = true
		if page.GetState() != Complete {
			incomplete = true
		}

		if incomplete && indexPage != nil {
			break
		}
	}

	if indexPage == nil || !pagesExist || incomplete {
		return
	}

	indexPage.SetCompleteState(Complete)
}

func (c *Chapter) String() string {
	result := fmt.Sprintln("  ", c.Title)

	c.Prepare()

	for _, page := range c.Pages {
		result += page.String()
	}

	return result
}

func (c *Chapter) GetErrors() []string {
	var errors []string

	for _, page := range c.Pages {
		errors = append(errors, page.GetErrors()...)
	}

	return errors
}

type Chapters []*Chapter

func (c Chapters) Add(filePath, chapterFN, pageFN string, content Content) Chapters {
	for i, chapter := range c {
		if chapter.Title == chapterFN {
			c[i].Pages = c[i].Pages.Add(filePath, pageFN, content)
			return c
		}
	}

	return append(c, &Chapter{Title: chapterFN, Pages: Pages{{FilePath: filePath, Title: pageFN, Content: content}}})
}

type Course struct {
	Title    string
	Chapters Chapters
}

func (c Course) Prepare() {
	for _, chapter := range c.Chapters {
		chapter.Prepare()
	}
}

func (c Course) String() string {
	result := fmt.Sprintln(c.Title)

	for _, chapter := range c.Chapters {
		result += chapter.String()
	}

	return result
}

func (c Course) GetErrors() []string {
	var issues []string

	for _, chapter := range c.Chapters {
		issues = append(issues, chapter.GetErrors()...)
	}

	return issues
}

func (c Course) Stats() (int, int, int, int, int) {
	var (
		total, stub, incomplete, complete, errors int
	)

	for _, chapter := range c.Chapters {
		for _, page := range chapter.Pages {
			switch page.GetState() {
			case Stub:
				stub++
			case Incomplete:
				incomplete++
			case Complete:
				complete++
			}

			if len(page.GetIssues()) > 0 {
				errors++
			}
		}

		total += len(chapter.Pages)
	}

	return total, stub, incomplete, complete, errors
}

type Courses []Course

func (c Courses) Add(filePath, courseFN, chapterFN, pageFN string, content Content) Courses {
	for i, course := range c {
		if course.Title == courseFN {
			c[i].Chapters = c[i].Chapters.Add(filePath, chapterFN, pageFN, content)
			return c
		}
	}

	return append(c, Course{Title: courseFN, Chapters: Chapters{{Title: chapterFN, Pages: Pages{{FilePath: filePath, Title: pageFN, Content: content}}}}})
}
