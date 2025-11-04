package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/getsentry/sentry-go"
	"github.com/jason0x43/go-toggl"
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
	jiraUser    string
}

func (s *togglJiraService) run(start, end, sinceDate time.Time) error {
	fmt.Printf("_______________________________________________\n")
	fmt.Printf("now\t%s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("from\t%s\n", start.Format(time.RFC3339))
	fmt.Printf("until\t%s\n", end.Format(time.RFC3339))

	togglEntries, err := s.getTogglEntries(start, end)
	if err != nil {
		return fmt.Errorf("cannot get time entries: %w", err)
	}

	fmt.Printf("\nwill process %d toggl entries\n\n", len(togglEntries))
	fmt.Printf("Issue ID\tTime\t\n")

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
	msg := fmt.Sprintf("%s\t%s\t", record.IssueID, record.TimeSpent)

	wl, _, err := s.jiraClient.GetWorklogs(record.IssueID)

	if err != nil {
		sentry.CaptureException(err)
		insertInfo <- fmt.Sprintf("%s error getting worklogs: %v", msg, err)
		return
	}

	for _, i := range wl.Worklogs {
		if i.Author.EmailAddress != s.jiraUser {
			continue
		}

		if i.Started.Equal(*record.Started) {
			if i.TimeSpent == record.TimeSpent {
				existing := fmt.Sprintf("https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s", record.IssueID, i.ID)

				insertInfo <- fmt.Sprintf("%s duplicate of %s", msg, existing)
				return
			}

			message := fmt.Sprintf("%s started at the same time %s, but time in toggl: %s and time in jira: %s", msg, time.Time(*record.Started).Format(time.RFC3339), record.TimeSpent, i.TimeSpent)
			sentry.CaptureMessage(message)

			insertInfo <- message
			return
		}

	}

	wlAdded, _, errAdded := s.jiraClient.AddWorklogRecord(record.IssueID, &record)
	if errAdded != nil {
		sentry.CaptureException(errAdded)

		insertInfo <- fmt.Sprintf("error adding worklog record: %v\n", err)
		return
	}

	insertInfo <- fmt.Sprintf("%s added: https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s", msg, record.IssueID, wlAdded.ID)
}
