package pkg

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCourses_Add(t *testing.T) {
	stubContent := Content{
		Title:           "Baz",
		State:           complete,
		CalculatedState: complete,
		Body: DefaultBody{
			HasMainVideo: true,
			HasPractice:  true,
		},
	}

	type args struct {
		courseFN  string
		chapterFN string
		pageFN    string
		content   Content
	}
	tests := []struct {
		name string
		args args
		c    Courses
		want Courses
	}{
		{
			name: "add to empty",
			args: args{
				courseFN:  "foo",
				chapterFN: "bar",
				pageFN:    "baz.md",
				content:   stubContent,
			},
			c: Courses{},
			want: Courses{
				{
					Title: "foo",
					Chapters: Chapters{
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "add to pages",
			args: args{
				courseFN:  "foo",
				chapterFN: "bar",
				pageFN:    "baz.md",
				content:   stubContent,
			},
			c: Courses{
				{
					Title: "foo",
					Chapters: Chapters{
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz0.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
			want: Courses{
				{
					Title: "foo",
					Chapters: Chapters{
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz0.md",
									Content:  stubContent,
								},
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "add to chapters",
			args: args{
				courseFN:  "foo",
				chapterFN: "bar",
				pageFN:    "baz.md",
				content:   stubContent,
			},
			c: Courses{
				{
					Title: "foo",
					Chapters: Chapters{
						{
							Title: "bar0",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
			want: Courses{
				{
					Title: "foo",
					Chapters: Chapters{
						{
							Title: "bar0",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "add to courses",
			args: args{
				courseFN:  "foo",
				chapterFN: "bar",
				pageFN:    "baz.md",
				content:   stubContent,
			},
			c: Courses{
				{
					Title: "foo0",
					Chapters: Chapters{
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
			want: Courses{
				{
					Title: "foo0",
					Chapters: Chapters{
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
				{
					Title: "foo",
					Chapters: Chapters{
						{
							Title: "bar",
							Pages: Pages{
								{
									Filename: "baz.md",
									Content:  stubContent,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.Add(tt.args.courseFN, tt.args.chapterFN, tt.args.pageFN, tt.args.content)
			assert.Equalf(t, tt.want, got, "Add(%v, %v, %v, %v)", tt.args.courseFN, tt.args.chapterFN, tt.args.pageFN, tt.args.content)
		})
	}
}

func TestCourses_ToRows(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		c    Courses
		args args
		want []table.Row
	}{
		{
			name: "default",
			c: Courses{
				{
					Title: "foo #0",
					Chapters: Chapters{
						{
							Title: "bar #0",
							Pages: Pages{
								{
									Filename: "baz-0.md",
									Content: Content{
										Title:           "Baz 0",
										FilePath:        "foo-0/bar-0/baz-0.md",
										State:           complete,
										CalculatedState: incomplete,
										Body: DefaultBody{
											HasMainVideo:     true,
											HasSummary:       true,
											HasTopics:        true,
											HasPractice:      false,
											HasRelatedVideos: true,
											HasRelatedLinks:  true,
										},
									},
								},
								{
									Filename: "baz-1.md",
									Content: Content{
										Title:           "Baz 1",
										FilePath:        "foo-0/bar-0/baz-1.md",
										State:           complete,
										CalculatedState: stub,
										Body: DefaultBody{
											HasMainVideo:     false,
											HasSummary:       false,
											HasTopics:        false,
											HasPractice:      false,
											HasRelatedVideos: false,
											HasRelatedLinks:  false,
										},
									},
								},
							},
						},
						{
							Title: "bar #1",
							Pages: Pages{
								{
									Filename: "baz-2.md",
									Content: Content{
										Title:           "Baz 2",
										FilePath:        "foo-0/bar-1/baz-2.md",
										State:           complete,
										CalculatedState: incomplete,
										Body: DefaultBody{
											HasMainVideo:     false,
											HasSummary:       true,
											HasTopics:        true,
											HasPractice:      true,
											HasRelatedVideos: true,
											HasRelatedLinks:  true,
										},
									},
								},
							},
						},
					},
				},
				{
					Title: "foo #1",
					Chapters: Chapters{
						{
							Title: "bar #2",
							Pages: Pages{
								{
									Filename: "baz-3.md",
									Content: Content{
										Title:           "Baz 3",
										FilePath:        "foo-1/bar-2/baz-3.md",
										State:           complete,
										CalculatedState: stub,
										Body: DefaultBody{
											HasMainVideo:     false,
											HasSummary:       false,
											HasTopics:        false,
											HasPractice:      false,
											HasRelatedVideos: false,
											HasRelatedLinks:  false,
										},
									},
								},
							},
						},
					},
				},
			},
			args: args{
				count: 4,
			},
			want: []table.Row{
				{
					"foo-0/bar-0/baz-0.md",
					"1",
					"foo #0",
					"bar #0",
					"Baz 0",
					"Y",
					"Y",
					"Y",
					"Y",
					"Y",
					"N",
					"complete",
					"incomplete",
					"N",
				},
				{
					"foo-0/bar-0/baz-1.md",
					"2",
					"foo #0",
					"bar #0",
					"Baz 1",
					"N",
					"N",
					"N",
					"N",
					"N",
					"N",
					"complete",
					"stub",
					"N",
				},
				{
					"foo-0/bar-1/baz-2.md",
					"3",
					"foo #0",
					"bar #1",
					"Baz 2",
					"N",
					"Y",
					"Y",
					"Y",
					"Y",
					"Y",
					"complete",
					"incomplete",
					"N",
				},
				{
					"foo-1/bar-2/baz-3.md",
					"4",
					"foo #1",
					"bar #2",
					"Baz 3",
					"N",
					"N",
					"N",
					"N",
					"N",
					"N",
					"complete",
					"stub",
					"N",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.c.ToRows(tt.args.count), "ToRows(%v)", tt.args.count)
		})
	}
}
