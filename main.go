package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	jiraScrumId = "REC-3123"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token_toggl := os.Getenv("TOKEN_TOGGL")
	if token_toggl == "" {
		panic("missing TOKEN_TOGGL")
	}

	togglEntries, err := getTogglEntries(token_toggl)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	jiraWorklogRecords := parseIssues(togglEntries)

	for _, issue := range jiraWorklogRecords {
		fmt.Println(issue)
	}
}
