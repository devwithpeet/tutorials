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
				State:    Unknown,
				FilePath: "foo.md",
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
				Title:    "Hello",
				State:    Unknown,
				FilePath: "foo.md",
				Body: DefaultBody{
					MainVideo: VideoProblem,
				},
			},
		},
		{
			name: "state-only",
			args: args{
				rawContent: `+++
state = "incomplete"
+++`,
				filePath: "15-foo.md",
			},
			want: Content{
				Title:    "",
				State:    Incomplete,
				FilePath: "15-foo.md",
				Weight:   "",
				Slug:     "",
				Body: DefaultBody{
					MainVideo: VideoProblem,
				},
			},
		},
		{
			name: "empty-chapter",
			args: args{
				rawContent: ``,
				filePath:   "_index.md",
			},
			want: Content{
				State:    Unknown,
				FilePath: "_index.md",
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
				Title:    "Hello",
				State:    Unknown,
				FilePath: "_index.md",
				Body:     ChapterBody{},
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
				Title:    "_index.md",
				State:    Incomplete,
				FilePath: "_index.md",
				Body:     ChapterBody{},
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
				Title:    "Hello",
				State:    Unknown,
				FilePath: "_index.md",
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
				Title:    "Hello",
				State:    Complete,
				FilePath: "_index.md",
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
				filePath: "10-foo.md",
			},
			want: Content{
				Title:    "Hello",
				State:    Complete,
				FilePath: "10-foo.md",
				Weight:   "",
				Slug:     "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     true,
					HasRelatedLinks: true,
					RelatedVideos: RelatedVideos{
						{
							Badge: "",
							Issues: []string{
								"missing time shortcode",
								"missing badge shortcode",
								"missing youtube shortcode",
							},
							Minutes: 0,
						},
					},
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
				filePath: "10-foo.md",
			},
			want: Content{
				Title:    "Hello",
				State:    Complete,
				FilePath: "10-foo.md",
				Weight:   "",
				Slug:     "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     true,
					HasRelatedLinks: true,
					RelatedVideos: RelatedVideos{
						{
							Badge: "",
							Issues: []string{
								"missing time shortcode",
								"missing badge shortcode",
								"missing youtube shortcode",
							},
							Minutes: 0,
						},
					},
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
				filePath: "10-foo.md",
			},
			want: Content{
				Title:    "Hello",
				State:    Complete,
				FilePath: "10-foo.md",
				Weight:   "",
				Slug:     "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     false,
					HasRelatedLinks: true,
					RelatedVideos: RelatedVideos{
						{
							Badge: "",
							Issues: []string{
								"missing time shortcode",
								"missing badge shortcode",
								"missing youtube shortcode",
							},
							Minutes: 0,
						},
					},
				},
			},
		},
		{
			name: "complete-page-with-hashmark-headers",
			args: args{
				rawContent: `+++
title = "Hello"
state = "complete"
weight = 9
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
				filePath: "9-foo.md",
			},
			want: Content{
				Title:    "Hello",
				State:    Complete,
				FilePath: "9-foo.md",
				Weight:   "9",
				Slug:     "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     true,
					HasRelatedLinks: true,
					RelatedVideos: RelatedVideos{
						{
							Badge: "",
							Issues: []string{
								"missing time shortcode",
								"missing badge shortcode",
								"missing youtube shortcode",
							},
							Minutes: 0,
						},
					},
				},
			},
		},
		{
			name: "bug unclear",
			args: args{
				rawContent: `+++
title = 'What Your Text Editor Says About You'
date = 2024-07-21T12:31:33+02:00
weight = 60
state = 'complete'
draft = false
slug = 'what-your-text-editor-says-about-you'
tags = ["no-practice", "fun", "vim", "vscode", "goland", "jetbrains"]
disableMermaid = true
disableOpenapi = true
audience = 'all'
audienceImportance = 'irrelevant'
+++

Main Video
----------

{{< time 5 >}}

This is just a fun video, don't take it too seriously. But also it's good to know what others will think about you based
on your choice of text editor. :D

{{< youtube sbdFwFDTDqU >}}
`,
				filePath: "9-foo.md",
			},
			want: Content{
				Title:    "What Your Text Editor Says About You",
				State:    Complete,
				FilePath: "9-foo.md",
				Weight:   "60",
				Slug:     "what-your-text-editor-says-about-you",
				Body: DefaultBody{
					MainVideo:       VideoPresent,
					HasSummary:      false,
					HasTopics:       false,
					HasPractice:     false,
					HasRelatedLinks: false,
					RelatedVideos:   nil,
				},
			},
		},
		{
			name: "bug unclear 2",
			args: args{
				rawContent: `+++
title = 'Installing Linux for OSX Users'
date = 2024-07-07T12:31:24+02:00
weight = 30
state = 'complete'
draft = false
slug = 'installing-linux-for-osx-users'
tags = ["linux", "no practice"]
disableMermaid = true
disableOpenapi = true
audience = "Mac users"
audienceImportance = "important"
outsideImportance = "optional"
+++

Main Video
----------

{{< main-really-missing >}}

Summary
-------

Okay, so the truth is, while I think Linux is a superior desktop in general, OSX is okay too. If you've already
invested your money into buying a Mac, you might be able to just get away using it. As it's a Unix based OS, most of the
command line tools work on OSX as well, and it's quite okay to switch back and forth between a Mac dev machine and a
Linux server if you need to. Learning the Linux commands would still make a lot of sense for you.

Related Videos
--------------

### Running Linux in a Virtual Machine

#### How to install Ubuntu in MacOS Apple Silicon Chip M1/M2/M3 - TopNotch Programmer

{{< time 2 >}} {{<badge-extra>}}

{{< youtube 3pj2Uck_Afc >}}


### Dual-booting Linux on your ARM Mac

OK, so this is not for the faint-hearted. **I do not recommend anyone to do this.** But I would totally do this on my own
MacBook, if I had one. :)

#### I use Arch on an M1 MacBook, btw - Fireship

{{< time 3 >}} {{<badge-fun>}}

{{< youtube j_I9nkpovCQ >}}
`,
				filePath: "9-foo.md",
			},
			want: Content{
				Title:    "Installing Linux for OSX Users",
				State:    Complete,
				FilePath: "9-foo.md",
				Weight:   "30",
				Slug:     "installing-linux-for-osx-users",
				Body: DefaultBody{
					MainVideo:       VideoReallyMissing,
					HasSummary:      true,
					HasTopics:       false,
					HasPractice:     false,
					HasRelatedLinks: false,
					RelatedVideos: RelatedVideos{
						{
							Badge:   "extra",
							Minutes: 2,
						},
						{
							Badge:   "fun",
							Minutes: 3,
						},
					},
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

func TestExtractRelatedVideos(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want RelatedVideos
	}{
		{
			name: "empty content is skipped",
			args: args{
				content: "",
			},
			want: nil,
		},
		{
			name: "no badges lead to issues",
			args: args{
				content: "### This is a title\n\nfoo\n",
			},
			want: RelatedVideos{
				{
					Badge: "",
					Issues: []string{
						"missing time shortcode",
						"missing badge shortcode",
						"missing youtube shortcode",
					},
					Minutes: 0,
				},
			},
		},
		{
			name: "too many badges lead to issues",
			args: args{
				content: `### This is it
{{< time 5 >}} {{<time  12>}}

{{<  badge-extra   >}} {{<badge-extra>}} 

{{< youtube abc >}} {{<youtube def>}}
`,
			},
			want: RelatedVideos{
				{
					Badge: "extra",
					Issues: []string{
						"multiple time shortcodes found",
						"unexpected badge shortcode found: extra",
						"multiple youtube shortcodes found",
					},
					Minutes: 5,
				},
			},
		},
		{
			name: "multiple badges, multiple issues",
			args: args{
				content: `### Almost empty

foo

### Missing badge

{{< time 5 >}}

{{< youtube abc >}}

### Multiple badge

{{< time 123 >}} {{<badge-alternative>}} {{<badge-extra>}}

{{< youtube foo >}}

### Multiple youtube videos

{{< time 17 >}} {{<badge-extra>}}

{{< youtube bar >}}
{{< youtube foo >}}
`,
			},
			want: RelatedVideos{
				{
					Badge: "",
					Issues: []string{
						"missing time shortcode",
						"missing badge shortcode",
						"missing youtube shortcode",
					},
					Minutes: 0,
				},
				{
					Badge: "",
					Issues: []string{
						"missing badge shortcode",
					},
					Minutes: 5,
				},
				{
					Badge: "alternative",
					Issues: []string{
						"unexpected badge shortcode found: extra",
					},
					Minutes: 123,
				},
				{
					Badge: "extra",
					Issues: []string{
						"multiple youtube shortcodes found",
					},
					Minutes: 17,
				},
			},
		},
		{
			name: "success",
			args: args{
				content: `### Multiple youtube videos

{{< time 17 >}} {{<badge-extra>}}

{{< youtube bar >}}
`,
			},
			want: RelatedVideos{
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 17,
				},
			},
		},
		{
			name: "success - empty skipped",
			args: args{
				content: `### Skipped
### Multiple youtube videos

{{< time 17 >}} {{<badge-extra>}}

{{< youtube bar >}}
`,
			},
			want: RelatedVideos{
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 17,
				},
			},
		},
		{
			name: "success - no-embed badge does not count as badge",
			args: args{
				content: `### Multiple youtube videos

{{< time 17 >}} {{<badge-extra>}} {{<badge-no-embed>}}
`,
			},
			want: RelatedVideos{
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 17,
				},
			},
		},
		{
			name: "failure - no-embed badge does count as a video",
			args: args{
				content: `### Multiple youtube videos

{{< time 17 >}} {{<badge-extra>}} {{<badge-no-embed>}}

{{< youtube bar >}}
`,
			},
			want: RelatedVideos{
				{
					Badge:   "extra",
					Issues:  []string{"unexpected youtube shortcode together with no-embed badge"},
					Minutes: 17,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute
			got := ExtractRelatedVideos(tt.args.content)

			// verify
			assert.Equal(t, tt.want, got)
		})
	}
}
