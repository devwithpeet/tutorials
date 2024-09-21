package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCourses_Add(t *testing.T) {
	stubContent := Content{
		Title: "Baz",
		State: Complete,
		Body: DefaultBody{
			MainVideo:   VideoProblem,
			HasPractice: true,
		},
	}

	type args struct {
		filePath  string
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
				filePath:  "foo/bar/baz.md",
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
									FilePath: "baz.md",
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
				filePath:  "foo/bar/baz.md",
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
									FilePath: "baz0.md",
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
									FilePath: "baz0.md",
									Content:  stubContent,
								},
								{
									FilePath: "foo/bar/baz.md",
									Title:    "baz.md",
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
				filePath:  "foo/bar/baz.md",
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
									FilePath: "baz.md",
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
									FilePath: "baz.md",
									Content:  stubContent,
								},
							},
						},
						{
							Title: "bar",
							Pages: Pages{
								{
									FilePath: "baz.md",
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
				filePath:  "foo/bar/baz.md",
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
									FilePath: "baz.md",
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
									FilePath: "baz.md",
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
									FilePath: "baz.md",
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
			// execute
			got := tt.c.Add(tt.args.filePath, tt.args.courseFN, tt.args.chapterFN, tt.args.pageFN, tt.args.content)

			// verify
			assert.Equal(t, tt.want, got)
		})
	}
}
