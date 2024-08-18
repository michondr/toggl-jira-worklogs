package main

import (
	"encoding/base64"
	"fmt"
	"github.com/jason0x43/go-toggl"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token_toggl := os.Getenv("TOKEN_TOGGL")
	if token_toggl == "" {
		panic("missing TOKEN_TOGGL")
	}

	issuesToinsertTojira := getTogglIssues(token_toggl, err)

	for _, issue := range issuesToinsertTojira {
		fmt.Println(issue)
	}
}

func getTogglIssues(token_toggl string, err error) []jiraTimeEntry {
	authToken := base64.StdEncoding.Strict().EncodeToString([]byte("michondr:" + token_toggl))
	fmt.Println(authToken)

	toggl.EnableLog()
	session := toggl.OpenSession(token_toggl)

	lastWorklogDateInJira := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	relevantEntries, err := session.GetTimeEntries(lastWorklogDateInJira, now)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	const (
		PROJECT_MEETING     = 202187965
		PROJECT_CODE_REVIEW = 202268075
		PROJECT_OPERATIVE   = 204586275
		PROJECT_FEATURE     = 202217954
	)

	projectIdToProject := map[int]toggl.Project{
		PROJECT_CODE_REVIEW: {Name: "cr"},
		PROJECT_MEETING:     {Name: "meetings"},
		PROJECT_OPERATIVE:   {Name: "někdo po mě něco chce"},
		PROJECT_FEATURE:     {Name: "feature"},
	}

	scrumIssue := jiraIssue{ID: "REC-3333"} //TODO: what's scrum ID

	issuesToinsertTojira := []jiraTimeEntry{}

	for _, entry := range relevantEntries {
		projectId := *entry.Pid

		_, exists := projectIdToProject[projectId]
		if !exists {
			panic("project " + string(rune(projectId)) + " not found in array: " + err.Error())
		}

		isREC := strings.HasPrefix(entry.Description, "REC-")
		hasDescription := len(entry.Description) > 8

		if isREC {
			descr := "code review"
			if hasDescription {
				descr = entry.Description[11:]
			}

			issuesToinsertTojira = append(
				issuesToinsertTojira,
				jiraTimeEntry{
					Issue:       jiraIssue{entry.Description[:8]},
					Description: descr,
					From:        entry.Start,
					To:          entry.Stop,
				},
			)

			continue
		}

		issuesToinsertTojira = append(
			issuesToinsertTojira,
			jiraTimeEntry{
				Issue:       scrumIssue,
				Description: entry.Description,
				From:        entry.Start,
				To:          entry.Stop,
			},
		)
	}
	return issuesToinsertTojira
}
