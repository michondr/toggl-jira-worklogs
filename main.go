package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/getsentry/sentry-go"
	"github.com/jason0x43/go-toggl"
	"github.com/joho/godotenv"
)

const (
	jiraScrumId       = "REC-3123"
	handleIssuesSince = "2024-08-20"
	sentryCronSlug    = "insert-to-jira"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
		return
	}

	defaultToday := time.Now()
	defaultTomorrow := defaultToday.AddDate(0, 0, 1)

	var (
		tokenToggl           = os.Getenv("TOGGL_TOKEN")
		jiraToken            = os.Getenv("JIRA_TOKEN")
		jiraUser             = os.Getenv("JIRA_USER")
		jiraUrl              = os.Getenv("JIRA_URL")
		sentryDsn            = os.Getenv("SENTRY_DSN")
		expectedCronSchedule = os.Getenv("EXPECTED_CRON_SCHEDULE")
		dateFrom             = flag.String("from", defaultToday.Format(time.DateOnly), "from date, default to today")
		dateTo               = flag.String("to", defaultTomorrow.Format(time.DateOnly), "to date, default to tomorrow")
		dateTz               = flag.String("tz", "Europe/Prague", "date timezone")
	)
	flag.Parse()

	monitorConfig, cronId, err := initSentry(sentryDsn, expectedCronSchedule)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

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

	sentry.CaptureCheckIn(
		&sentry.CheckIn{
			ID:          *cronId,
			MonitorSlug: sentryCronSlug,
			Status:      sentry.CheckInStatusOK,
		},
		monitorConfig,
	)
}

func initSentry(sentryDsn, expectedCronSchedule string) (*sentry.MonitorConfig, *sentry.EventID, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: sentryDsn,
		//EnableLogs: true,
	})
	if err != nil {
		return nil, nil, err
	}

	monitorSchedule := sentry.CrontabSchedule(expectedCronSchedule)

	// Create a monitor config object
	monitorConfig := &sentry.MonitorConfig{
		Schedule:      monitorSchedule,
		MaxRuntime:    2,
		CheckInMargin: 1,
	}

	// 🟡 Notify Sentry your job is running:
	checkinId := sentry.CaptureCheckIn(
		&sentry.CheckIn{
			MonitorSlug: sentryCronSlug,
			Status:      sentry.CheckInStatusInProgress,
		},
		monitorConfig,
	)
	//
	//// The SentryLogger requires context, to link logs with the appropriate traces. You can either create a new logger
	//// by providing the context, or use WithCtx() to pass the context inline.
	//ctx := context.Background()
	//logger := sentry.NewLogger(ctx)
	//
	//// You can use the logger like [fmt.Print]
	//logger.Info().Emit("Hello ", "world!")
	//// Or like [fmt.Printf]
	//logger.Info().Emitf("Hello %v!", "world")

	return monitorConfig, checkinId, nil
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
