package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMarkdown(t *testing.T) {
	t.Run("panic on broken", func(t *testing.T) {
		filePath := "foo.md"
		rawContent := "+++\n???"

		assert.Panics(t, func() { ParseMarkdown(rawContent, filePath) })
	})

	type args struct {
		rawContent string
		filePath   string
	}
	tests := []struct {
		name string
		args args
		want Content
	}{
		{
			name: "empty",
			args: args{
				rawContent: ``,
				filePath:   "foo.md",
			},
			want: Content{
				State:           unknown,
				CalculatedState: stub,
				FilePath:        "foo.md",
			},
		},
		{
			name: "title-only",
			args: args{
				rawContent: `+++
title = "Hello"
+++`,
				filePath: "foo.md",
			},
			want: Content{
				Title:           "Hello",
				State:           unknown,
				CalculatedState: stub,
				FilePath:        "foo.md",
				Body:            DefaultBody{},
			},
		},
		{
			name: "state-only",
			args: args{
				rawContent: `+++
state = "incomplete"
+++`,
				filePath: "foo.md",
			},
			want: Content{
				Title:           "foo.md",
				State:           incomplete,
				CalculatedState: stub,
				FilePath:        "foo.md",
				Body:            DefaultBody{},
			},
		},
		{
			name: "empty-chapter",
			args: args{
				rawContent: ``,
				filePath:   "_index.md",
			},
			want: Content{
				State:           unknown,
				CalculatedState: stub,
				FilePath:        "_index.md",
			},
		},
		{
			name: "title-only-chapter",
			args: args{
				rawContent: `+++
title = "Hello"
+++`,
				filePath: "_index.md",
			},
			want: Content{
				Title:           "Hello",
				State:           unknown,
				CalculatedState: unknown,
				FilePath:        "_index.md",
				Body:            ChapterBody{},
			},
		},
		{
			name: "state-only-chapter",
			args: args{
				rawContent: `+++
state = "incomplete"
+++`,
				filePath: "_index.md",
			},
			want: Content{
				Title:           "_index.md",
				State:           incomplete,
				CalculatedState: stub,
				FilePath:        "_index.md",
				Body:            ChapterBody{},
			},
		},
		{
			name: "complete-chapter-without-state",
			args: args{
				rawContent: `+++
title = "Hello"
+++
Episodes
--------

- bar
`,
				filePath: "_index.md",
			},
			want: Content{
				Title:           "Hello",
				State:           unknown,
				CalculatedState: unknown,
				FilePath:        "_index.md",
				Body: ChapterBody{
					HasEpisodes: true,
				},
			},
		},
		{
			name: "complete-chapter-with-state",
			args: args{
				rawContent: `+++
title = "Hello"
state = "complete"
+++
Episodes
--------

- bar
`,
				filePath: "_index.md",
			},
			want: Content{
				Title:           "Hello",
				State:           complete,
				CalculatedState: complete,
				FilePath:        "_index.md",
				Body: ChapterBody{
					HasEpisodes: true,
				},
			},
		},
		{
			name: "almost-complete-page",
			args: args{
				rawContent: `+++
title = "Hello"
state = "complete"
+++
Summary
-------

- bar

Main Video
----------

Topics
------

- bar

Related Videos
--------------

- bar

Related Links
-------------

- bar

Practice
--------

- bar
`,
				filePath: "foo.md",
			},
			want: Content{
				Title:           "Hello",
				State:           complete,
				CalculatedState: incomplete,
				FilePath:        "foo.md",
				Body: DefaultBody{
					HasMainVideo:     false,
					HasSummary:       true,
					HasTopics:        true,
					HasPractice:      true,
					HasRelatedLinks:  true,
					HasRelatedVideos: true,
				},
			},
		},
		{
			name: "complete-page",
			args: args{
				rawContent: `+++
title = "Hello"
state = "complete"
+++
Summary
-------

- bar

Main Video
----------

- bar

Topics
------

- bar

Related Videos
--------------

- bar

Related Links
-------------

- bar

Practice
--------

- bar
`,
				filePath: "foo.md",
			},
			want: Content{
				Title:           "Hello",
				State:           complete,
				CalculatedState: complete,
				FilePath:        "foo.md",
				Body: DefaultBody{
					HasMainVideo:     true,
					HasSummary:       true,
					HasTopics:        true,
					HasPractice:      true,
					HasRelatedLinks:  true,
					HasRelatedVideos: true,
				},
			},
		},
		{
			name: "incomplete-if-practice-is-missing",
			args: args{
				rawContent: `+++
title = "Hello"
state = "complete"
+++
## Summary

- bar

## Main Video

- bar

## Topics

- bar

## Related Videos

- bar

## Related Links

- bar
`,
				filePath: "foo.md",
			},
			want: Content{
				Title:           "Hello",
				State:           complete,
				CalculatedState: incomplete,
				FilePath:        "foo.md",
				Body: DefaultBody{
					HasMainVideo:     true,
					HasSummary:       true,
					HasTopics:        true,
					HasPractice:      false,
					HasRelatedLinks:  true,
					HasRelatedVideos: true,
				},
			},
		},
		{
			name: "complete-page-with-hashmark-headers",
			args: args{
				rawContent: `+++
title = "Hello"
state = "complete"
+++
## Summary

- bar

## Main Video

- bar

## Topics

- bar

## Related Videos

- bar

## Related Links

- bar

## Practice

- bar
`,
				filePath: "foo.md",
			},
			want: Content{
				Title:           "Hello",
				State:           complete,
				CalculatedState: complete,
				FilePath:        "foo.md",
				Body: DefaultBody{
					HasMainVideo:     true,
					HasSummary:       true,
					HasTopics:        true,
					HasPractice:      true,
					HasRelatedLinks:  true,
					HasRelatedVideos: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseMarkdown(tt.args.rawContent, tt.args.filePath)

			assert.Equal(t, tt.want, got)
		})
	}
}
