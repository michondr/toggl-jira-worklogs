package main

import (
	"github.com/jason0x43/go-toggl"
	"time"
)

func getTogglEntries(token_toggl string) ([]toggl.TimeEntry, error) {
	toggl.EnableLog()
	session := toggl.OpenSession(token_toggl)

	now := time.Now()
	lastWorklogDateInJira := getLastWorklogDateInJira()

	return session.GetTimeEntries(lastWorklogDateInJira, now)
}
