package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jason0x43/go-toggl"
	"regexp"
	"strings"
	"time"
)

var recIssueRegex = regexp.MustCompile(`^REC-\d+`)

func parseIssues(togglEntries []toggl.TimeEntry) []jira.WorklogRecord {
	worklogRecords := make([]jira.WorklogRecord, len(togglEntries))

	for key, entry := range togglEntries {
		started := jira.Time(*entry.Start)

		if issueID := recIssueRegex.FindString(entry.Description); issueID != "" {
			descr := "code review"
			if rest := strings.TrimPrefix(entry.Description, issueID); strings.HasPrefix(rest, " - ") {
				descr = rest[3:]
			}

			worklogRecords[key] = jira.WorklogRecord{
				IssueID:   issueID,
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

	if hour > 8 {
		days := hour / 8
		hours := hour % 8

		if minute == 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd %dh %dm", days, hours, minute)
	}

	if hour > 0 {
		if minute == 0 {
			return fmt.Sprintf("%dh", hour)
		}
		return fmt.Sprintf("%dh %dm", hour, minute)
	}

	return fmt.Sprintf("%dm", minute)
}
