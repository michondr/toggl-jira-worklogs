package main

import (
	"github.com/jason0x43/go-toggl"
	"time"
)

func getTogglEntries(token_toggl string, date_to_run time.Time) ([]toggl.TimeEntry, error) {
	toggl.DisableLog()
	session := toggl.OpenSession(token_toggl)

	start := date_to_run.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	return session.GetTimeEntries(start, end)
}
