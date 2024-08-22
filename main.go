package main

import (
	"flag"
	"fmt"
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

	var (
		tokenToggl    = os.Getenv("TOGGL_TOKEN")
		jiraToken     = os.Getenv("JIRA_TOKEN")
		jiraUser      = os.Getenv("JIRA_USER")
		jiraUrl       = os.Getenv("JIRA_URL")
		dateToProcess = flag.String("date", time.Now().Format(time.DateOnly), "date to process")
		dateTz        = flag.String("tz", "Europe/Prague", "date timezone")
	)
	flag.Parse()

	tz, err := time.LoadLocation(*dateTz)
	if err != nil {
		panic("cannot find tz: " + err.Error())
	}

	handleIssuesSince := time.Date(2024, 8, 20, 0, 0, 0, 0, tz)
	forDate, err := time.Parse(time.DateOnly, *dateToProcess)

	if err != nil {
		panic("cannot parse date: " + err.Error())
	}
	if forDate.Compare(handleIssuesSince) == -1 {
		panic("cannot go this far back")
	}

	forDate = forDate.In(tz)
	fmt.Printf("processing date %s\n", forDate)

	togglEntries, err := getTogglEntries(tokenToggl, forDate)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	fmt.Printf("will process %d toggl entries\n", len(togglEntries))

	for _, entry := range parseIssues(togglEntries) {

		if time.Time(*entry.Started).Compare(handleIssuesSince) == -1 {
			panic("really, updating something this far would be bad")
		}

		insertToJiraIfNotExists(entry, jiraUser, jiraToken, jiraUrl)
	}

	fmt.Printf("done\n")
}
