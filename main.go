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

	defaultToday := time.Now()
	defaultTomorrow := defaultToday.AddDate(0, 0, 1)

	var (
		tokenToggl = os.Getenv("TOGGL_TOKEN")
		jiraToken  = os.Getenv("JIRA_TOKEN")
		jiraUser   = os.Getenv("JIRA_USER")
		jiraUrl    = os.Getenv("JIRA_URL")
		dateFrom   = flag.String("from", defaultToday.Format(time.DateOnly), "from date, default to today")
		dateTo     = flag.String("to", defaultTomorrow.Format(time.DateOnly), "to date, default to tomorrow")
		dateTz     = flag.String("tz", "Europe/Prague", "date timezone")
	)
	flag.Parse()

	service := togglJiraService{
		togglClient: loginToToggl(tokenToggl),
		jiraClient:  loginToJira(jiraUser, jiraToken, jiraUrl).Issue,
		jiraUser:    jiraUser,
	}

	tz, err := time.LoadLocation(*dateTz)
	if err != nil {
		panic("cannot find tz")
	}
	start, err := time.ParseInLocation(time.DateOnly, *dateFrom, tz)
	if err != nil {
		panic("cannot parse from date")
	}
	end, err := time.ParseInLocation(time.DateOnly, *dateTo, tz)
	if err != nil {
		panic("cannot parse to date")
	}
	sinceDate, _ := time.ParseInLocation(time.DateOnly, handleIssuesSince, tz)

	if start.Compare(sinceDate) == -1 {
		panic("cannot go this far back")
	}

	if err := service.run(start, end, sinceDate); err != nil {
		log.Fatal(err)
	}
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
	toggl.DisableLog()

	return &ses
}
