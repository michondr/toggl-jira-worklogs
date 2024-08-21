package main

import (
	"github.com/jason0x43/go-toggl"
	"reflect"
	"testing"
	"time"
)

func Test_parseIssues(t *testing.T) {
	suStart := time.Date(2024, 8, 20, 8, 30, 0, 0, time.UTC)
	suEnd := time.Date(2024, 8, 20, 9, 0, 0, 0, time.UTC)

	su1Start := time.Date(2024, 8, 21, 8, 30, 0, 0, time.UTC)
	su1End := time.Date(2024, 8, 21, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          []toggl.TimeEntry
		expectedOutput []jiraTimeEntry
	}{
		{
			"standup record",
			[]toggl.TimeEntry{
				{
					Description: "SU",
					Start:       &suStart,
					Stop:        &suEnd,
				}, {
					Description: "SU",
					Start:       &su1Start,
					Stop:        &su1End,
				},
			},
			[]jiraTimeEntry{
				{
					jiraIssue{ID: jiraScrumId},
					"SU",
					&suStart,
					&suEnd,
				}, {
					jiraIssue{ID: jiraScrumId},
					"SU",
					&su1Start,
					&su1End,
				},
			},
		},
		{
			"feature ticket",
			[]toggl.TimeEntry{
				{
					Description: "REC-1234 - some feature name",
					Start:       &suStart,
					Stop:        &suEnd,
				},
			},
			[]jiraTimeEntry{
				{
					jiraIssue{ID: "REC-1234"},
					"some feature name",
					&suStart,
					&suEnd,
				},
			},
		},
		{
			"cr ticket",
			[]toggl.TimeEntry{
				{
					Description: "REC-1234",
					Start:       &suStart,
					Stop:        &suEnd,
				},
			},
			[]jiraTimeEntry{
				{
					jiraIssue{ID: "REC-1234"},
					"code review",
					&suStart,
					&suEnd,
				},
			},
		}, {
			"other",
			[]toggl.TimeEntry{
				{
					Description: "foo bar",
					Start:       &suStart,
					Stop:        &suEnd,
				},
			},
			[]jiraTimeEntry{
				{
					jiraIssue{ID: jiraScrumId},
					"foo bar",
					&suStart,
					&suEnd,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("parse toggle entries to jira worklog records", func(t *testing.T) {
			if got := parseIssues(tt.input); !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("parseIssues() = %v, expectedOutput %v", got, tt.expectedOutput)
			}
		})
	}
}
