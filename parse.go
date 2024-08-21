package main

import (
	"github.com/jason0x43/go-toggl"
	"strings"
)

func parseIssues(togglEntries []toggl.TimeEntry) []jiraTimeEntry {
	scrumIssue := jiraIssue{ID: jiraScrumId}

	issuesToinsertTojira := []jiraTimeEntry{}

	for _, entry := range togglEntries {
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
