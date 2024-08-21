package main

import (
	"github.com/jason0x43/go-toggl"
	"time"
)

func getTogglIssuesInJiraFormat(token_toggl string, err error) []jiraTimeEntry {
	toggl.EnableLog()
	session := toggl.OpenSession(token_toggl)

	now := time.Now()
	lastWorklogDateInJira := getLastWorklogDateInJira()

	relevantEntries, err := session.GetTimeEntries(lastWorklogDateInJira, now)
	if err != nil {
		panic("cannot get time entries: " + err.Error())
	}

	return parseIssues(relevantEntries)
}
