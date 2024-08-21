package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"time"
)

func loginToJira(jira_user, jira_token, jira_url string) *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: jira_user,
		Password: jira_token,
	}

	client, err := jira.NewClient(tp.Client(), jira_url)
	if err != nil {
		fmt.Printf("\nerror client: %v\n", err)
		return nil
	}

	return client
}

func insertToJiraIfNotExists(record jira.WorklogRecord, jira_user, jira_token, jira_url string) {
	client := loginToJira(jira_user, jira_token, jira_url)

	wl, _, err := client.Issue.GetWorklogs(record.IssueID)

	if err != nil {
		fmt.Printf("\nerror getting worklogs: %v\n", err)
		return
	}

	for _, i := range wl.Worklogs {
		if i.Started.Equal(*record.Started) && i.TimeSpent == record.TimeSpent {
			fmt.Printf("\nis duplicate ID: %s, spent %s from %s\n", record.IssueID, record.TimeSpent, time.Time(*record.Started).Format(time.RFC3339))
			return
		}

	}

	wlAdded, _, errAdded := client.Issue.AddWorklogRecord(record.IssueID, &record)
	if errAdded != nil {
		fmt.Printf("\nerror adding worklog record: %v\n", err)
		return
	}

	fmt.Printf("\nworklog record (%s, %s) added: https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s\n", record.IssueID, record.TimeSpent, record.IssueID, wlAdded.ID)
}
