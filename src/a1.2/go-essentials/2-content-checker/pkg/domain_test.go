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
