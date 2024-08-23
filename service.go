package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"net/http"
	"time"
)

type togglClient interface {
	GetTimeEntries(startDate, endDate time.Time) ([]toggl.TimeEntry, error)
}

type jiraClient interface {
	GetWorklogs(issueID string, options ...func(*http.Request) error) (*jira.Worklog, *jira.Response, error)
	AddWorklogRecord(issueID string, record *jira.WorklogRecord, options ...func(*http.Request) error) (*jira.WorklogRecord, *jira.Response, error)
}
type togglJiraService struct {
	togglClient togglClient
	jiraClient  jiraClient
}

func (s *togglJiraService) run(dateToProcess, dateTz *string) {
	tz, err := time.LoadLocation(*dateTz)
	if err != nil {
		panic("cannot find tz: " + err.Error())
	}

	sinceDate, _ := time.Parse(time.DateOnly, handleIssuesSince)
	forDate, err := time.Parse(time.DateOnly, *dateToProcess)
	if err != nil {
		panic("cannot parse date: " + err.Error())
	}

	sinceDate = sinceDate.In(tz)
	forDate = forDate.In(tz)
	if forDate.Compare(sinceDate) == -1 {
		panic("cannot go this far back")
	}

	fmt.Printf("processing date %s\n", forDate)

	togglEntries, err := s.getTogglEntries(forDate)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	fmt.Printf("will process %d toggl entries\n", len(togglEntries))

	for _, entry := range s.transformEntries(togglEntries) {

		if time.Time(*entry.Started).Compare(sinceDate) == -1 {
			panic("really, updating something this far would be bad")
		}

		s.insertToJiraIfNotExists(entry)
	}

	fmt.Printf("done\n")
}

func (s *togglJiraService) getTogglEntries(forDate time.Time) ([]toggl.TimeEntry, error) {
	start := forDate.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	return s.togglClient.GetTimeEntries(start, end)
}

func (s *togglJiraService) transformEntries(togglEntries []toggl.TimeEntry) []jira.WorklogRecord {
	return parseIssues(togglEntries)
}

func (s *togglJiraService) insertToJiraIfNotExists(record jira.WorklogRecord) {
	wl, _, err := s.jiraClient.GetWorklogs(record.IssueID)

	if err != nil {
		fmt.Printf("error getting worklogs: %v\n", err)
		return
	}

	for _, i := range wl.Worklogs {
		if i.Started.Equal(*record.Started) && i.TimeSpent == record.TimeSpent {
			fmt.Printf("is duplicate ID: %s, spent %s from %s\n", record.IssueID, record.TimeSpent, time.Time(*record.Started).Format(time.RFC3339))
			return
		}

	}

	wlAdded, _, errAdded := s.jiraClient.AddWorklogRecord(record.IssueID, &record)
	if errAdded != nil {
		fmt.Printf("error adding worklog record: %v\n", err)
		return
	}

	fmt.Printf("worklog record (%s, %s) added: https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s\n", record.IssueID, record.TimeSpent, record.IssueID, wlAdded.ID)
}
