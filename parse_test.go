package main

import (
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"strconv"
	"testing"
	"time"
)

func Test_parseIssues(t *testing.T) {
	t1Start := time.Date(2024, 8, 20, 8, 30, 0, 0, time.UTC)
	t1End := time.Date(2024, 8, 20, 9, 0, 0, 0, time.UTC)
	t1Jira := jira.Time(t1Start)

	t2Start := time.Date(2024, 8, 21, 8, 30, 0, 0, time.UTC)
	t2End := time.Date(2024, 8, 21, 10, 28, 12, 0, time.UTC)
	t2Jira := jira.Time(t2Start)

	tests := []struct {
		name           string
		input          []toggl.TimeEntry
		expectedOutput []jira.WorklogRecord
	}{
		{
			"standup record",
			[]toggl.TimeEntry{
				{
					Description: "SU",
					Start:       &t1Start,
					Stop:        &t1End,
				}, {
					Description: "SU",
					Start:       &t2Start,
					Stop:        &t2End,
				},
			},
			[]jira.WorklogRecord{
				{
					IssueID:   jiraScrumId,
					Comment:   "SU",
					Started:   &t1Jira,
					TimeSpent: "30m",
				}, {
					IssueID:   jiraScrumId,
					Comment:   "SU",
					Started:   &t2Jira,
					TimeSpent: "1h 58m",
				},
			},
		},
		{
			"feature ticket",
			[]toggl.TimeEntry{
				{
					Description: "REC-1234 - some feature name",
					Start:       &t1Start,
					Stop:        &t1End,
				},
			},
			[]jira.WorklogRecord{
				{
					IssueID:   "REC-1234",
					Comment:   "some feature name",
					Started:   &t1Jira,
					TimeSpent: "30m",
				},
			},
		},
		{
			"cr ticket",
			[]toggl.TimeEntry{
				{
					Description: "REC-1234",
					Start:       &t1Start,
					Stop:        &t1End,
				},
			},
			[]jira.WorklogRecord{
				{
					IssueID:   "REC-1234",
					Comment:   "code review",
					Started:   &t1Jira,
					TimeSpent: "30m",
				},
			},
		}, {
			"other",
			[]toggl.TimeEntry{
				{
					Description: "foo bar",
					Start:       &t1Start,
					Stop:        &t1End,
				},
			},
			[]jira.WorklogRecord{
				{
					IssueID:   jiraScrumId,
					Comment:   "foo bar",
					Started:   &t1Jira,
					TimeSpent: "30m",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("check length", func(t *testing.T) {
				if len(tt.input) != len(tt.expectedOutput) {
					t.Errorf("input and expected output do not match in length")
				}
			})

			actualAll := parseIssues(tt.input)
			expectedAll := tt.expectedOutput

			for i := 0; i < len(tt.input); i++ {
				t.Run("index"+strconv.Itoa(i), func(t *testing.T) {
					actual := actualAll[i]
					expected := expectedAll[i]

					t.Run("issue ID", func(t *testing.T) {
						if actual.IssueID != expected.IssueID {
							t.Errorf("issue ID %s does not match with expected %s", actual.IssueID, expected.IssueID)
						}
					})
					t.Run("comment", func(t *testing.T) {
						if actual.Comment != expected.Comment {
							t.Errorf("comment %s does not match with expected %s", actual.Comment, expected.Comment)
						}
					})
					t.Run("started", func(t *testing.T) {
						if actual.Started == nil {
							t.Errorf("actual is nil")
							return
						}
						if actual.Started == nil {
							t.Errorf("expected is nil")
							return
						}
						if !actual.Started.Equal(*expected.Started) {
							t.Errorf("started does not match with expected ")
							return
						}
					})
					t.Run("timeSpent", func(t *testing.T) {
						if actual.TimeSpent != expected.TimeSpent {
							t.Errorf("timeSpent %s does not match with expected %s", actual.TimeSpent, expected.TimeSpent)
						}
					})
				})
			}
		})
	}
}

func Test_timeToTimeSpent(t *testing.T) {
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1 second",
			args: args{
				from: time.Date(2024, 8, 20, 8, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 8, 20, 8, 0, 1, 0, time.UTC),
			},
			want: "0m",
		},
		{
			name: "1 minute 0 second",
			args: args{
				from: time.Date(2024, 8, 20, 8, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 8, 20, 8, 1, 0, 0, time.UTC),
			},
			want: "1m",
		},
		{
			name: "1 minute 1 second",
			args: args{
				from: time.Date(2024, 8, 20, 8, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 8, 20, 8, 1, 1, 0, time.UTC),
			},
			want: "1m",
		},
		{
			name: "1 hour 1 minute 1 second",
			args: args{
				from: time.Date(2024, 8, 20, 8, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 8, 20, 9, 1, 1, 0, time.UTC),
			},
			want: "1h 1m",
		},
		{
			name: "1 day 1 hour 1 minute 1 second",
			args: args{
				from: time.Date(2024, 8, 20, 8, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 8, 21, 9, 1, 1, 0, time.UTC),
			},
			want: "25h 1m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeToTimeSpent(tt.args.to.Sub(tt.args.from)); got != tt.want {
				t.Errorf("timeToTimeSpent() = %v, want %v", got, tt.want)
			}
		})
	}
}
