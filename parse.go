package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"strings"
	"time"
)

func parseIssues(togglEntries []toggl.TimeEntry) []jira.WorklogRecord {
	worklogRecords := make([]jira.WorklogRecord, len(togglEntries))

	for key, entry := range togglEntries {
		started := jira.Time(*entry.Start)

		if isREC := strings.HasPrefix(entry.Description, "REC-"); isREC {
			descr := "code review"
			if hasDescription := len(entry.Description) > 8; hasDescription {
				descr = entry.Description[11:]
			}

			worklogRecords[key] = jira.WorklogRecord{
				IssueID:   entry.Description[:8],
				Comment:   descr,
				Started:   &started,
				TimeSpent: timeToTimeSpent(entry.Stop.Sub(*entry.Start)),
			}

			continue
		}

		worklogRecords[key] = jira.WorklogRecord{
			IssueID:   jiraScrumId,
			Comment:   entry.Description,
			Started:   &started,
			TimeSpent: timeToTimeSpent(entry.Stop.Sub(*entry.Start)),
		}
	}
	return worklogRecords
}

func timeToTimeSpent(d time.Duration) string {
	if d < 0 {
		panic("time diff is negative")
	}

	hour := int(d.Hours())
	minute := int(d.Minutes()) % 60

	return fmt.Sprintf("%dh %dm", hour, minute)
}
