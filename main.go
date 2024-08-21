package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	jiraScrumId = "REC-3123"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token_toggl := os.Getenv("TOGGL_TOKEN")
	if token_toggl == "" {
		panic("missing TOGGL_TOKEN")
	}

	jira_token := os.Getenv("JIRA_TOKEN")
	if jira_token == "" {
		panic("missing JIRA_TOKEN")
	}
	jira_user := os.Getenv("JIRA_USER")
	if jira_user == "" {
		panic("missing JIRA_USER")
	}
	jira_url := os.Getenv("JIRA_URL")
	if jira_url == "" {
		panic("missing JIRA_URL")
	}

	tz, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		panic("cannot find tz: " + err.Error())
	}
	handleIssuesSince := time.Date(2024, 8, 20, 0, 0, 0, 0, tz)

	forDate := time.Date(2024, 8, 21, 12, 14, 15, 0, tz)

	if forDate.Compare(handleIssuesSince) == -1 {
		panic("cannot go this far back")
	}

	togglEntries, err := getTogglEntries(token_toggl, forDate)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	for _, entry := range parseIssues(togglEntries) {

		if time.Time(*entry.Started).Compare(handleIssuesSince) == -1 {
			panic("really, updating something this far would be bad")
		}

		insertToJiraIfNotExists(entry, jira_user, jira_token, jira_url)
	}

}
