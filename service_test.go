package main

import (
	"errors"
	"testing"
)

func Test_togglJiraService_run(t *testing.T) {
	type fields struct {
		togglClient togglClient
		jiraClient  jiraClient
	}
	type args struct {
		dateToProcess *string
		dateTz        *string
	}

	tz := "Europe/Prague"
	invalidDate := "2024-08-19"

	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedErr error
	}{
		{
			"run before automation should fail",
			fields{
				nil,
				nil,
			},
			args{
				&invalidDate,
				&tz,
			},
			errors.New("cannot go this far back"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &togglJiraService{
				togglClient: tt.fields.togglClient,
				jiraClient:  tt.fields.jiraClient,
			}
			err := s.run(tt.args.dateToProcess, tt.args.dateTz)

			if tt.expectedErr == nil && err != nil {
				t.Errorf("togglJiraService.run() should have returned nil")
				return
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("togglJiraService.run() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
