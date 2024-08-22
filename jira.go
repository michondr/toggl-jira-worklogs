package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"time"
)

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

func insertToJiraIfNotExists(record jira.WorklogRecord, jiraUser, jiraToken, jiraUrl string) {
	client := loginToJira(jiraUser, jiraToken, jiraUrl)

	wl, _, err := client.Issue.GetWorklogs(record.IssueID)

	if err != nil {
		fmt.Printf("error getting worklogs: %v\n", err)
		return
	}

	for _, i := range wl.Worklogs {
		if i.Started.Equal(*record.Started) && i.TimeSpent == record.TimeSpent {
			fmt.Printf("is duplicate ID: %s, spent %s from %s\n", record.IssueID, record.TimeSpent, time.Time(*record.Started).Format(time.RFC3339))
			return
		}

	}

	wlAdded, _, errAdded := client.Issue.AddWorklogRecord(record.IssueID, &record)
	if errAdded != nil {
		fmt.Printf("error adding worklog record: %v\n", err)
		return
	}

	fmt.Printf("worklog record (%s, %s) added: https://recruitis.atlassian.net/browse/%s?focusedWorklogId=%s\n", record.IssueID, record.TimeSpent, record.IssueID, wlAdded.ID)
}
