package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type jiraIssue struct {
	ID string
}
type jiraTimeEntry struct {
	Issue       jiraIssue
	Description string
	From        *time.Time
	To          *time.Time
}

const (
	jiraScrumId = "REC-3333" //TODO: check it matches
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

	issuesToinsertTojira := getTogglIssuesInJiraFormat(token_toggl, err)

	for _, issue := range issuesToinsertTojira {
		fmt.Println(issue)
	}
}
