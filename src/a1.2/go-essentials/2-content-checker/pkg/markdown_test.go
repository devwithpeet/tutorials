package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMarkdown(t *testing.T) {
	t.Run("panic on broken", func(t *testing.T) {
		rawContent := "+++\n???"

		// execute
		content, err := ParseMarkdown(rawContent)
		require.Error(t, err)

		// verify
		assert.Empty(t, content)
	})

	type args struct {
		rawContent string
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
			},
			want: Content{
				State: Unknown,
			},
		},
		{
			name: "title-only",
			args: args{
				rawContent: `+++
title = "Prepare"
+++`,
			},
			want: Content{
				Title: "Prepare",
				State: Unknown,
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
			},
			want: Content{
				Title:  "",
				State:  Incomplete,
				Weight: "",
				Slug:   "",
				Body: DefaultBody{
					MainVideo: VideoProblem,
				},
			},
		},
		{
			name: "empty-chapter",
			args: args{
				rawContent: ``,
			},
			want: Content{
				State: Unknown,
			},
		},
		{
			name: "title-only-chapter",
			args: args{
				rawContent: `+++
title = "Prepare"
+++`,
			},
			want: Content{
				Title: "Prepare",
				State: Unknown,
				Body: DefaultBody{
					MainVideo: VideoProblem,
				},
			},
		},
		{
			name: "state-only-chapter",
			args: args{
				rawContent: `+++
state = "incomplete"
+++`,
			},
			want: Content{
				State: Incomplete,
				Body: DefaultBody{
					MainVideo: VideoProblem,
				},
			},
		},
		{
			name: "complete-chapter-without-state",
			args: args{
				rawContent: `+++
title = "Prepare"
+++
Episodes
--------

- bar
`,
			},
			want: Content{
				Title: "Prepare",
				State: Unknown,
				Body: &IndexBody{
					HasEpisodes:   true,
					CompleteState: Incomplete,
				},
			},
		},
		{
			name: "complete-chapter-with-state",
			args: args{
				rawContent: `+++
title = "Prepare"
state = "complete"
+++
Episodes
--------

- bar
`,
			},
			want: Content{
				Title: "Prepare",
				State: Complete,
				Body: &IndexBody{
					HasEpisodes:   true,
					CompleteState: Incomplete,
				},
			},
		},
		{
			name: "almost-complete-page",
			args: args{
				rawContent: `+++
title = "Prepare"
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
			},
			want: Content{
				Title:  "Prepare",
				State:  Complete,
				Weight: "",
				Slug:   "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     true,
					HasRelatedLinks: true,
					RelatedVideos:   RelatedVideos{},
				},
			},
		},
		{
			name: "complete-page",
			args: args{
				rawContent: `+++
title = "Prepare"
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
			},
			want: Content{
				Title:  "Prepare",
				State:  Complete,
				Weight: "",
				Slug:   "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     true,
					HasRelatedLinks: true,
					RelatedVideos:   RelatedVideos{},
				},
			},
		},
		{
			name: "incomplete-if-practice-is-missing",
			args: args{
				rawContent: `+++
title = "Prepare"
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
			},
			want: Content{
				Title:  "Prepare",
				State:  Complete,
				Weight: "",
				Slug:   "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     false,
					HasRelatedLinks: true,
					RelatedVideos:   RelatedVideos{},
				},
			},
		},
		{
			name: "complete-page-with-hashmark-headers",
			args: args{
				rawContent: `+++
title = "Prepare"
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
			},
			want: Content{
				Title:  "Prepare",
				State:  Complete,
				Weight: "9",
				Slug:   "",
				Body: DefaultBody{
					MainVideo:       VideoProblem,
					HasSummary:      true,
					HasTopics:       true,
					HasPractice:     true,
					HasRelatedLinks: true,
					RelatedVideos:   RelatedVideos{},
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
			},
			want: Content{
				Title:  "What Your Text Editor Says About You",
				State:  Complete,
				Weight: "60",
				Slug:   "what-your-text-editor-says-about-you",
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
			name: "practice",
			args: args{
				rawContent: `+++
title = 'Data Cleanup'
date = 2024-07-09T19:26:57+02:00
weight = 20
state = 'complete'
draft = false
slug = 'data-cleanup'
tags = ["vim", "practice"]
disableMermaid = true
disableOpenapi = true
audience = "all"
audienceImportance = "important"
+++

Description
-----------

Download [this SQL File](/a1.1/practice-data-cleanup.sql).

At this point, you don't really need to understand what this file is about, all you need to know is that we want to
turn it into a JSON file, using the values found in the second parentheses.

So basically your task is not turn that file into something that looks like this:

You don't need to worry about the white spaces the following two examples are also acceptable solutions:

### Examples

#### Example 1

#### Example 2

#### Example 3

### Hints

**Hint:** Arguably the fastest solution is using an editor with Vim motions and Vim macros, but other solutions are fine
as well. If you're familiar with tools like grep, sed or awk, those can be quite efficient for tasks like this too.

Recommended challenges
----------------------

### Display overall stats

Write an app to display the coordinates (x, y) for the largest, and smallest values for the whole dataset.

Example output:


### Display stats for each chart

Write an app to display the coordinates (x, y) for the largest, and smallest values for each chart.

So an example output could be the following:

Note that the order of the stat blocks does not matter, lemmy could come before lemmy.


Additional challenges
---------------------

{{<badge-extra>}}

### Sorting

This one is only different from the "Display chart stats" challenge is that here the order of the stats matter, they
should be ordered by the chart name, ordered Z to A, plus we should display all the coordinates where the value is the
maximum or minimum and make sure that they're ordered in incremental order.

Example output:

### Find the size of chart maps

Find a program that is able to tell the size of the map we have complete coordinates, meaning that no points are missing.

Example output:


### Find the size of intended chart maps and errors

People make mistakes, the provided SQL file is also imperfect. Modify your application so that it can give reviewers a
hint on missing coordinates or duplicates.

Example output:
`,
			},
			want: Content{
				Title:  "Data Cleanup",
				State:  Complete,
				Weight: "20",
				Slug:   "data-cleanup",
				Body: &PracticeBody{
					HasDescription:           true,
					HasRecommendedChallenges: true,
					HasAdditionalChallenges:  true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute
			got, err := ParseMarkdown(tt.args.rawContent)
			require.NoError(t, err)

			// verify
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
				content: `### This is a title\n\nfoo\n
{{< time 5 >}}

{{< youtube abc >}}`,
			},
			want: RelatedVideos{
				{
					Badge: "",
					Issues: []string{
						"missing badge shortcode",
					},
					Minutes: 5,
					Valid:   true,
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
					Valid:   true,
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
						"missing badge shortcode",
					},
					Minutes: 5,
					Valid:   true,
				},
				{
					Badge: "alternative",
					Issues: []string{
						"unexpected badge shortcode found: extra",
					},
					Minutes: 123,
					Valid:   true,
				},
				{
					Badge: "extra",
					Issues: []string{
						"multiple youtube shortcodes found",
					},
					Minutes: 17,
					Valid:   true,
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
					Valid:   true,
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
					Valid:   true,
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
					Valid:   true,
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
					Valid:   true,
				},
			},
		},
		{
			name: "complex video section",
			args: args{
				content: `### The Analytical Engine (Charles Babbage, Ada Lovelace)

#### The greatest machine that never was - John Graham-Cumming - TED-Ed

{{< time 12 >}} {{<badge-extra>}}

{{< youtube FlfChYGv3Z4 >}}

#### Babbage's Analytical Engine - Computerphile

{{< time 14 >}} {{<badge-extra>}}

{{< youtube 5rtKoKFGFSM >}}

#### Ada Lovelace: The First Computer Programmer - Biographics

{{< time 21 >}} {{<badge-extra>}}

{{< youtube id=IZptxisyVqQ start=60 >}}

---

### Harvard Mark I

#### Supercomputer Where It All Started - Harvard Mark 1 - Major Hardware

{{< time 6 >}} {{<badge-extra>}}

{{< youtube cd2DV-AoCk4 >}}

#### Harvard Mark I, 2022 - CS50

{{< time 3 >}} {{<badge-extra>}}

{{< youtube 7l8W96I7_ew >}}

---

### Enigma, Bombe (Alan Turing)

#### How did the Enigma Machine work? - Jared Owen

{{< time 20 >}} {{<badge-extra>}}

{{< youtube ybkkiGtJmkM >}}


### Lorenz and Colossus (Tommy Flowers, Bill Tutte)

#### Why the Toughest Code to Break in WW2 WASN'T Enigma - The Story of the Lorenz Cipher

{{< time 11 >}} {{<badge-extra>}}

{{< youtube RCWgOaDOzpY >}}

#### Colossus & Bletchley Park - Computerphile

{{< time 9 >}} {{<badge-extra>}}

{{< youtube 9HH-asvLAj4 >}}

#### Colossus - The Greatest Secret in the History of Computing - The Centre for Computing History

{{< time 60 >}} {{<badge-extra>}}

This is not only about Colossus, but provides a lot of context, including basic cryptographic problems of the time. It's
probably my favorite video recommended on this page.

{{< youtube g2tMcMQqSbA >}}

### Why Build Colossus? (Bill Tutte) - Computerphile

{{< time 8 >}} {{<badge-extra>}}

{{< youtube 1f82-aTYNb8 >}}


---

### Transistors and ENIAC (John Mauchly, J. Presper Eckert)

#### Transistors - The Invention That Changed The World - Real Engineering

{{< time 8 >}} {{<badge-extra>}}

{{< youtube OwS9aTE2Go4 >}}
`,
			},
			want: RelatedVideos{
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 12,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 14,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 21,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 6,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 3,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 20,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 11,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 9,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 60,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 8,
					Valid:   true,
				},
				{
					Badge:   "extra",
					Issues:  nil,
					Minutes: 8,
					Valid:   true,
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
