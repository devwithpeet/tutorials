package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	sectionMainVideo             = "main video"
	sectionSummary               = "summary"
	sectionTopics                = "topics"
	sectionRelatedVideos         = "related videos"
	sectionRelatedLinks          = "related links"
	sectionExercises             = "exercises"
	sectionEpisodes              = "episodes"
	sectionDescription           = "description"
	sectionRecommendedChallenges = "recommended challenges"
	sectionAdditionalChallenges  = "additional challenges"
)

const markdownHeaderLength = 3

func ParseMarkdown(rawContent string) (Content, error) {
	if len(rawContent) < markdownHeaderLength*2 {
		return Content{State: Unknown}, nil
	}

	// Convert DOS/Windows line endings (\r\n) into Linux/Unix line endings
	strContent := strings.Replace(rawContent, "\r\n", EOL, -1)

	header, body, err := splitMarkdown(strContent)
	if err != nil {
		return Content{}, fmt.Errorf("markdown header could not be extracted, err: %w", err)
	}

	sections := extractSection(body)

	_, hasEpisodes := sections[sectionEpisodes]
	if hasEpisodes {
		return NewIndexContent(header, sections), nil
	}

	_, hasDescription := sections[sectionDescription]
	if hasDescription {
		return NewPracticeContent(header, sections), nil
	}

	return NewDefaultContent(header, sections), nil
}

func splitMarkdown(in string) (string, string, error) {
	// Handle TOML front matter
	if in[:4] == "+++\n" {
		if idx := strings.Index(in[4:], "\n+++"); idx != -1 {
			return in[4 : idx+4], strings.Trim(in[idx+8:], "\n+"), nil
		}
	}

	return "", "", errors.New("could not split markdown")
}

func NewIndexContent(header string, sections map[string]string) Content {
	title := getHeaderValue(header, "title", "")
	body := sectionsToIndexBody(sections)
	state := State(getHeaderValue(header, "state", string(Unknown)))

	return Content{
		Title: title,
		State: state,
		Body:  body,
	}
}

func NewPracticeContent(header string, sections map[string]string) Content {
	title := getHeaderValue(header, "title", "")
	state := State(getHeaderValue(header, "state", string(Unknown)))
	slug := getHeaderValue(header, "slug", "")
	weight := getHeaderValue(header, "weight", "")
	body := sectionsToPracticeBody(sections)

	return Content{
		Title:  title,
		State:  state,
		Body:   body,
		Slug:   slug,
		Weight: weight,
	}
}

func NewDefaultContent(header string, sections map[string]string) Content {
	title := getHeaderValue(header, "title", "")
	state := NewState(getHeaderValue(header, "state", string(Unknown)))
	slug := getHeaderValue(header, "slug", "")
	weight := getHeaderValue(header, "weight", "")
	tags := getHeaderValues(header, "tags", nil)
	defaultBody := sectionsToDefaultBody(sections, tags)

	return Content{
		Title:  title,
		State:  state,
		Body:   defaultBody,
		Slug:   slug,
		Weight: weight,
	}
}

var regexHeader = regexp.MustCompile(`^(\S+)\s*=\s*(.*)$`)

func getHeaderValue(header, key, defaultValue string) string {
	for _, row := range strings.Split(header, "\n") {
		matches := regexHeader.FindStringSubmatch(row)

		if len(matches) != 3 {
			continue
		}

		if matches[1] == key {
			return strings.Trim(matches[2], `'"`)
		}
	}

	return defaultValue
}

func getHeaderValues(header, key string, defaultValue []string) []string {
	for _, row := range strings.Split(header, "\n") {
		matches := regexHeader.FindStringSubmatch(row)

		if len(matches) != 3 {
			continue
		}

		if matches[1] != key {
			continue
		}

		var tags []string
		for _, part := range strings.Split(strings.Trim(matches[2], `[]`), ",") {
			tags = append(tags, strings.Trim(part, ` '"`))
		}

		return tags
	}

	return defaultValue
}

func extractSection(body string) map[string]string {
	sections := make(map[string]string)

	currentSection := "root"
	sectionStart := 0

	rows := strings.Split(body, EOL)
	for i, row := range rows {
		if len(row) >= 3 && row[:3] == "## " {
			content := strings.Trim(strings.Join(rows[sectionStart:i], EOL), " \t\n")
			if len(content) > 0 {
				sections[currentSection] = content
			}

			sectionStart = i + 1

			currentSection = strings.ToLower(strings.Trim(row[3:], " \t"))

			continue
		}

		if i > sectionStart && len(row) >= 3 && row[:3] == "---" {
			content := strings.Trim(strings.Join(rows[sectionStart:i-1], EOL), " \t\n")
			if len(content) > 0 {
				sections[currentSection] = content
			}

			sectionStart = i + 1

			currentSection = strings.ToLower(strings.Trim(rows[i-1], " \t"))

			continue
		}
	}

	if len(rows) > sectionStart {
		content := strings.Trim(strings.Join(rows[sectionStart:], EOL), " \t\n")
		if len(content) > 0 {
			sections[currentSection] = content
		}
	}

	return sections
}

var regexMissing = regexp.MustCompile(`{{<\s*main-missing\s*>}}`)
var regexReallyMissing = regexp.MustCompile(`{{<\s*main-really-missing\s*>}}`)
var regexYoutube = regexp.MustCompile(`{{<\s*youtube\s+([^>]*)\s*>}}`)

func ExtractMainVideo(content string) MainVideo {
	matchCount := 0
	mainVideo := VideoProblem

	matches := regexMissing.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		matchCount += len(matches)
		mainVideo = VideoMissing
	}

	matches = regexReallyMissing.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		matchCount += len(matches)
		mainVideo = VideoReallyMissing
	}

	matches = regexYoutube.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		if matchCount == 0 {
			return VideoPresent
		}

		matchCount++
	}

	if matchCount != 1 {
		return VideoProblem
	}

	return mainVideo
}

var regexSubHeader = regexp.MustCompile(`\n####?#? .*\n`)

func ExtractRelatedVideos(content string) RelatedVideos {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	sections := regexSubHeader.Split("\n"+content, -1)

	relatedVideos := make(RelatedVideos, 0, len(sections))
	for _, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		relatedVideo := extractRelatedVideo(section)

		if relatedVideo.Valid {
			relatedVideos = append(relatedVideos, relatedVideo)
		}
	}

	return relatedVideos
}

var regexTime = regexp.MustCompile(`{{<\s*time\s+(\d+)\s*>}}`)

func extractTime(content string) (int, []string) {
	var (
		issues  []string
		minutes int
		err     error
	)

	timeMatches := regexTime.FindAllStringSubmatch(content, -1)
	if len(timeMatches) == 0 {
		issues = append(issues, "missing time shortcode")
	} else {
		minutes, err = strconv.Atoi(timeMatches[0][1])
		if err != nil {
			issues = append(issues, fmt.Sprintf("failed to parse duration: %s", timeMatches[0][1]))
		}
	}
	if len(timeMatches) > 1 {
		issues = append(issues, "multiple time shortcodes found")
	}

	return minutes, issues
}

var regexBadge = regexp.MustCompile(`{{<\s*badge-(\S*)\s*>}}`)

func extractBadges(content string) (Badge, bool, []string) {
	var (
		badges []Badge
		issues []string
	)

	noEmbed := false
	badgeMatches := regexBadge.FindAllStringSubmatch(content, -1)

	for _, match := range badgeMatches {
		switch badge := Badge(match[1]); badge {
		case Unchecked, Alternative, Extra, Fun, Hint, MustSee, Summary:
			badges = append(badges, badge)
		case NoEmbed:
			noEmbed = true
		default:
			issues = append(issues, fmt.Sprintf("Unknown badge: '%s'", badge))
		}
	}

	if len(badges) == 0 {
		issues = append(issues, "missing badge shortcode")

		return "", noEmbed, issues
	} else if len(badges) > 1 {
		for _, badge := range badges[1:] {
			if badge == NoEmbed {
				continue
			}

			issues = append(issues, "unexpected badge shortcode found: "+string(badge))
		}
	}

	return badges[0], noEmbed, issues
}

func extractYoutube(content string, noEmbed bool) (int, []string) {
	var issues []string

	youtubeMatches := regexYoutube.FindAllStringSubmatch(content, -1)

	switch len(youtubeMatches) {
	case 0:
		if !noEmbed {
			issues = append(issues, "missing youtube shortcode")
		}
	case 1:
		if noEmbed {
			issues = append(issues, "unexpected youtube shortcode together with no-embed badge")
		}
	default:
		issues = append(issues, "multiple youtube shortcodes found")
	}

	return len(youtubeMatches), issues
}

func extractRelatedVideo(content string) RelatedVideo {
	var (
		badge   Badge
		issues  []string
		minutes int
	)

	minutes, timeIssues := extractTime(content)
	issues = append(issues, timeIssues...)

	badge, noEmbed, badgeIssues := extractBadges(content)
	issues = append(issues, badgeIssues...)

	ytCount, ytIssues := extractYoutube(content, noEmbed)
	issues = append(issues, ytIssues...)

	if ytCount == 0 && !noEmbed && badge == "" && minutes == 0 {
		return RelatedVideo{}
	}

	return RelatedVideo{
		Badge:   badge,
		Issues:  issues,
		Minutes: minutes,
		Valid:   true,
	}
}

const (
	tagUsefulWithoutVideo = "useful-without-video"
	tagNoExercise         = "no-exercise"
)

func sectionsToDefaultBody(sections map[string]string, tags []string) DefaultBody {
	_, hasSummary := sections[sectionSummary]
	_, hasTopics := sections[sectionTopics]
	_, hasRelatedLinks := sections[sectionRelatedLinks]
	exercises, hasExercises := sections[sectionExercises]

	mainVideo := ExtractMainVideo(sections[sectionMainVideo])
	relatedVideos := ExtractRelatedVideos(sections[sectionRelatedVideos])

	if hasExercises && strings.TrimSpace(exercises) == "" {
		hasExercises = false
	}

	usefulWithoutVideo := false
	for _, tag := range tags {
		if tag == tagUsefulWithoutVideo {
			usefulWithoutVideo = true
		}
		if tag == tagNoExercise {
			hasExercises = true
		}
	}

	return DefaultBody{
		MainVideo:          mainVideo,
		HasSummary:         hasSummary,
		HasTopics:          hasTopics,
		RelatedVideos:      relatedVideos,
		HasRelatedLinks:    hasRelatedLinks,
		HasExercises:       hasExercises,
		UsefulWithoutVideo: usefulWithoutVideo,
	}
}

func sectionsToIndexBody(sections map[string]string) *IndexBody {
	_, hasEpisodes := sections[sectionEpisodes]

	return &IndexBody{
		HasEpisodes:   hasEpisodes,
		CompleteState: Incomplete,
	}
}

func sectionsToPracticeBody(sections map[string]string) *PracticeBody {
	_, hasDescription := sections[sectionDescription]
	_, hasRecommendedChallenges := sections[sectionRecommendedChallenges]
	_, hasAdditionalChallenges := sections[sectionAdditionalChallenges]

	return &PracticeBody{
		HasDescription:           hasDescription,
		HasRecommendedChallenges: hasRecommendedChallenges,
		HasAdditionalChallenges:  hasAdditionalChallenges,
	}
}
