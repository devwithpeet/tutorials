package main

import (
	"reflect"
	"testing"
)

func Test_stats(t *testing.T) {
	type args struct {
		data [][]string
	}
	tests := []struct {
		name string
		args args
		want []Stats
	}{
		{
			name: "default",
			args: args{
				data: [][]string{
					{"", "Software Engineer", "", "", "50000"},
					{"", "Software Engineer", "", "", "60000"},
					{"", "Software Engineer", "", "", "45000"},
					{"", "QA", "", "", "45000"},
					{"", "Software Engineer", "", "", "45000"},
					{"", "QA", "", "", "42000"},
				},
			},
			want: []Stats{
				{
					Position: "Software Engineer",
					Average:  50000,
					Count:    4,
				},
				{
					Position: "QA",
					Average:  43500,
					Count:    2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := parse(tt.args.data)

			if got := stats(parsed); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
