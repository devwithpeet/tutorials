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
		want [][2]string
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
			want: [][2]string{
				{
					"Software Engineer",
					"50000",
				},
				{
					"QA",
					"43500",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stats(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stats() = %v, want %v", got, tt.want)
			}
		})
	}
}
