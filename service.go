package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"net/http"
	"sync"
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

func (s *togglJiraService) run(dateToProcess, dateTz *string) error {
	tz, err := time.LoadLocation(*dateTz)
	if err != nil {
		return fmt.Errorf("cannot find tz: %w", err)
	}

	sinceDate, _ := time.ParseInLocation(time.DateOnly, handleIssuesSince, tz)
	forDate, err := time.ParseInLocation(time.DateOnly, *dateToProcess, tz)
	if err != nil {
		return fmt.Errorf("cannot parse date: %w", err)
	}
	if forDate.Compare(sinceDate) == -1 {
		return fmt.Errorf("cannot go this far back")
	}

	start := forDate
	end := forDate.AddDate(0, 0, 1)

	fmt.Printf("from  %s\n", start.Format(time.RFC3339))
	fmt.Printf("until %s\n", end.Format(time.RFC3339))

	togglEntries, err := s.getTogglEntries(start, end)
	if err != nil {
		return fmt.Errorf("cannot get time entries: %w", err)
	}

	fmt.Printf("will process %d toggl entries\n\n", len(togglEntries))

	insertInfo := make(chan string)
	var wg sync.WaitGroup

	for _, entry := range s.transformEntries(togglEntries) {

		if time.Time(*entry.Started).Compare(sinceDate) == -1 {
			return fmt.Errorf("really, updating something this far would be bad. entry %s", entry.ID)
		}

		wg.Add(1)
		go s.insertToJiraIfNotExists(entry, &wg, insertInfo)
	}

	go func() {
		wg.Wait()
		close(insertInfo)
	}()

	for msg := range insertInfo {
		fmt.Println(msg)
	}
	fmt.Println("done")

	return nil
}

func (s *togglJiraService) getTogglEntries(start, end time.Time) ([]toggl.TimeEntry, error) {
	return s.togglClient.GetTimeEntries(start, end)
}

func (s *togglJiraService) transformEntries(togglEntries []toggl.TimeEntry) []jira.WorklogRecord {
	return parseIssues(togglEntries)
}

func (s *togglJiraService) insertToJiraIfNotExists(record jira.WorklogRecord, wg *sync.WaitGroup, insertInfo chan<- string) {
	defer wg.Done()

	wl, _, err := s.jiraClient.GetWorklogs(record.IssueID)

	if err != nil {
		insertInfo <- fmt.Sprintf("error getting worklogs: %v", err)
		return
	}

	for _, i := range wl.Worklogs {
		if i.Started.Equal(*record.Started) && i.TimeSpent == record.TimeSpent {
			insertInfo <- fmt.Sprintf(
				"is duplicate ID: %s, spent %s from %s (of %s)",
				record.IssueID,
				record.TimeSpent,
				time.Time(*record.Started).Format(time.RFC3339),
				fmt.Sprintf("https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s", record.IssueID, i.ID))
			return
		}

	}

	wlAdded, _, errAdded := s.jiraClient.AddWorklogRecord(record.IssueID, &record)
	if errAdded != nil {
		insertInfo <- fmt.Sprintf("error adding worklog record: %v\n", err)
		return
	}

	insertInfo <- fmt.Sprintf("worklog record (%s, %s) added: https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s", record.IssueID, record.TimeSpent, record.IssueID, wlAdded.ID)
}
