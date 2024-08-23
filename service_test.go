package main

import (
	"errors"
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"net/http"
	"testing"
	"time"
)

func Test_togglJiraService_run(t *testing.T) {
	type fields struct {
		recordsTogglReturns []toggl.TimeEntry
		recordsJiraExpects  []jira.WorklogRecord
	}
	type args struct {
		dateToProcess *string
		dateTz        *string
	}

	tz := "Europe/Prague"
	invalidDate := "2024-08-19"
	validDate := "2024-08-21"
	t1Start, _ := time.Parse(time.RFC3339, "2024-08-21T06:15:00+00:00")
	t1StartJira := jira.Time(t1Start)
	t1End, _ := time.Parse(time.RFC3339, "2024-08-21T07:15:00+00:00")
	t2Start, _ := time.Parse(time.RFC3339, "2024-08-21T04:30:00+00:00")
	t2StartJira := jira.Time(t2Start)
	t2End, _ := time.Parse(time.RFC3339, "2024-08-21T05:45:00+00:00")
	pid := 212216954

	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedErr error
	}{
		{
			"run before automation should fail",
			fields{
				recordsTogglReturns: []toggl.TimeEntry{},
				recordsJiraExpects:  []jira.WorklogRecord{},
			},
			args{
				&invalidDate,
				&tz,
			},
			errors.New("cannot go this far back"),
		},
		{
			"feature and scrum",
			fields{
				recordsTogglReturns: []toggl.TimeEntry{
					{
						ID:          3574476483,
						Wid:         2111248,
						Pid:         &pid,
						Tid:         nil,
						Description: "REC-5085 - Nové UI kalendáře - migrace dat",
						Start:       &t1Start,
						Stop:        &t1End,
						Tags:        nil,
						Duration:    3600,
						DurOnly:     true,
						Billable:    false,
					},
					{
						ID:          3574476484,
						Wid:         2111248,
						Pid:         &pid,
						Tid:         nil,
						Description: "test issue",
						Start:       &t2Start,
						Stop:        &t2End,
						Tags:        nil,
						Duration:    4500,
						DurOnly:     true,
						Billable:    false,
					},
				},
				recordsJiraExpects: []jira.WorklogRecord{
					{
						Started:   &t1StartJira,
						TimeSpent: "1h 0m",
						IssueID:   "REC-5085",
						Comment:   "Nové UI kalendáře - migrace dat",
					},
					{
						Started:   &t2StartJira,
						TimeSpent: "1h 15m",
						IssueID:   "REC-3123",
						Comment:   "test issue",
					},
				},
			},
			args{
				&validDate,
				&tz,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &togglJiraService{
				togglClient: testTogglClient{
					tt.fields.recordsTogglReturns,
				},
				jiraClient: testJiraClient{
					t,
					tt.fields.recordsJiraExpects,
				},
			}
			actualErr := s.run(tt.args.dateToProcess, tt.args.dateTz)

			if tt.expectedErr == nil && actualErr != nil {
				t.Errorf("togglJiraService.run() error = %v, wantErr nil", actualErr)
			} else if tt.expectedErr != nil && tt.expectedErr.Error() != actualErr.Error() {
				t.Errorf("togglJiraService.run() error = %v, wantErr %v", actualErr, tt.expectedErr)
			}
		})
	}
}

type testTogglClient struct {
	recordsToReturn []toggl.TimeEntry
}

func (t testTogglClient) GetTimeEntries(startDate, endDate time.Time) ([]toggl.TimeEntry, error) {
	return t.recordsToReturn, nil
}

type testJiraClient struct {
	t               *testing.T
	expectedRecords []jira.WorklogRecord
}

func (t testJiraClient) GetWorklogs(issueID string, options ...func(*http.Request) error) (*jira.Worklog, *jira.Response, error) {
	emptyWorklog := jira.Worklog{Worklogs: []jira.WorklogRecord{}}

	return &emptyWorklog, nil, nil
}
func (tc testJiraClient) AddWorklogRecord(issueID string, record *jira.WorklogRecord, options ...func(*http.Request) error) (*jira.WorklogRecord, *jira.Response, error) {
	tc.t.Run(issueID+record.Comment, func(t *testing.T) {
		var expectedRecord jira.WorklogRecord

		for _, r := range tc.expectedRecords {
			if r.IssueID == issueID {
				expectedRecord = r
			}
		}

		if expectedRecord.IssueID == "" {
			panic("expected record missing issueId, cannot assert")
		}

		t.Run("TimeSpent", func(t *testing.T) {
			if record.TimeSpent != expectedRecord.TimeSpent {
				t.Errorf("expected.TimeSpent = %v, actual = %v", expectedRecord.TimeSpent, record.TimeSpent)
			}
		})
		t.Run("Started", func(t *testing.T) {
			if *record.Started != *expectedRecord.Started {
				t.Errorf("expected.Started = %v, actual = %v", expectedRecord.Started, record.Started)
			}
		})
		t.Run("Comment", func(t *testing.T) {
			if record.Comment != expectedRecord.Comment {
				t.Errorf("expected.Comment = %v, actual = %v", expectedRecord.Comment, record.Comment)
			}
		})

	})
	return &jira.WorklogRecord{ID: "record-for-" + issueID}, nil, nil
}
