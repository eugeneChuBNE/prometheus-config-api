package main

import (
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	tests := []struct {
		command []string
		want    string
		wantErr bool
	}{
		{[]string{"echo", "Hello, world!"}, "Hello, world!\n", false},
		{[]string{"false"}, "", true},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.command, " "), func(t *testing.T) {
			got, err := executeCommand(tt.command...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("executeCommand() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
