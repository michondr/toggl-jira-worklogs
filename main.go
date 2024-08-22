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

	tokenToggl := os.Getenv("TOGGL_TOKEN")
	if tokenToggl == "" {
		panic("missing TOGGL_TOKEN")
	}

	jiraToken := os.Getenv("JIRA_TOKEN")
	if jiraToken == "" {
		panic("missing JIRA_TOKEN")
	}
	jiraUser := os.Getenv("JIRA_USER")
	if jiraUser == "" {
		panic("missing JIRA_USER")
	}
	jiraUrl := os.Getenv("JIRA_URL")
	if jiraUrl == "" {
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

	togglEntries, err := getTogglEntries(tokenToggl, forDate)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	for _, entry := range parseIssues(togglEntries) {

		if time.Time(*entry.Started).Compare(handleIssuesSince) == -1 {
			panic("really, updating something this far would be bad")
		}

		insertToJiraIfNotExists(entry, jiraUser, jiraToken, jiraUrl)
	}

}
