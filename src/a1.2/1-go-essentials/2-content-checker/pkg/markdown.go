package pkg

import (
	"errors"
	"fmt"
	"strings"
)

const (
	sectionMainVideo     = "main video"
	sectionSummary       = "summary"
	sectionTopics        = "topics"
	sectionRelatedVideos = "related videos"
	sectionRelatedLinks  = "related links"
	sectionPractice      = "practice"
	sectionEpisodes      = "episodes"
)

const markdownHeaderLength = 3

func ParseMarkdown(rawContent, filePath string) Content {
	if len(rawContent) < markdownHeaderLength*2 {
		return Content{
			FilePath:        filePath,
			State:           unknown,
			CalculatedState: stub,
		}
	}

	// Convert DOS/Windows line endings (\r\n) into Linux/Unix line endings
	strContent := strings.Replace(rawContent, "\r\n", EOL, -1)

	header, body, err := splitMarkdown(strContent)
	if err != nil {
		panic(fmt.Errorf("markdown header could not be extracted, file: %s, err: %w", filePath, err))
	}

	sections := extractSection(body)

	if strings.Contains(filePath, "_index.md") {
		return chapterContent(header, sections, filePath)
	}

	return defaultContent(header, sections, filePath)
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

func chapterContent(header string, sections map[string]string, filePath string) Content {
	body := chapterToContentBody(sections)
	calculatedState := unknown

	state := State(getHeaderValue(header, "state", string(unknown)))
	if state != unknown {
		calculatedState = stub
		if _, ok := sections[sectionEpisodes]; ok {
			calculatedState = complete
		}
	}

	return Content{
		Title:           getHeaderValue(header, "title", filePath),
		State:           state,
		CalculatedState: calculatedState,
		Body:            body,
		FilePath:        filePath,
	}
}

func defaultContent(header string, sections map[string]string, filePath string) Content {
	body := defaultToContentBody(sections)

	calculatedState := stub
	if _, ok := sections[sectionMainVideo]; ok {
		calculatedState = complete
		if _, ok := sections[sectionPractice]; !ok {
			calculatedState = incomplete
		}
	} else if _, ok := sections[sectionRelatedVideos]; ok {
		calculatedState = incomplete
	} else if _, ok := sections[sectionRelatedLinks]; ok {
		calculatedState = incomplete
	}

	return Content{
		Title:           getHeaderValue(header, "title", filePath),
		State:           State(getHeaderValue(header, "state", string(unknown))),
		CalculatedState: calculatedState,
		Body:            body,
		FilePath:        filePath,
	}
}

func getHeaderValue(header, key, defaultValue string) string {
	prefixLength := len(key) + 3
	for _, row := range strings.Split(header, "\n") {
		if len(row) > prefixLength && row[:prefixLength] == key+" = " {
			char := row[prefixLength]
			if char == '\'' || char == '"' {
				return strings.Trim(row[prefixLength:], string(char))
			}
		}
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

func defaultToContentBody(sections map[string]string) DefaultBody {
	_, hasSummary := sections[sectionSummary]
	_, hasVideo := sections[sectionMainVideo]
	_, hasTopics := sections[sectionTopics]
	_, hasRelatedVideos := sections[sectionRelatedVideos]
	_, hasRelatedLinks := sections[sectionRelatedLinks]
	_, hasPractice := sections[sectionPractice]

	return DefaultBody{
		HasMainVideo:     hasVideo,
		HasSummary:       hasSummary,
		HasTopics:        hasTopics,
		HasRelatedVideos: hasRelatedVideos,
		HasRelatedLinks:  hasRelatedLinks,
		HasPractice:      hasPractice,
	}
}

func chapterToContentBody(sections map[string]string) ChapterBody {
	_, hasEpisodes := sections[sectionEpisodes]

	return ChapterBody{
		HasEpisodes: hasEpisodes,
	}
}
