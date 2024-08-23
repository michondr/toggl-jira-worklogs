package main

import (
	"flag"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	jiraScrumId       = "REC-3123"
	handleIssuesSince = "2024-08-20"
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

	service := togglJiraService{
		togglClient: loginToToggl(tokenToggl),
		jiraClient:  loginToJira(jiraUser, jiraToken, jiraUrl).Issue,
	}

	service.run(dateToProcess, dateTz)
}

func loginToJira(jiraUser, jiraToken, jiraUrl string) *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: jiraUser,
		Password: jiraToken,
	}

	client, err := jira.NewClient(tp.Client(), jiraUrl)
	if err != nil {
		fmt.Printf("\nerror client: %v\n", err)
		return nil
	}

	return client
}

func loginToToggl(tokenToggl string) *toggl.Session {
	ses := toggl.OpenSession(tokenToggl)

	return &ses
}
