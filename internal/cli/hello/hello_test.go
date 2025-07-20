package hello

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelloCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "basic hello",
			args:       []string{},
			wantOutput: "Hello, World!",
		},
		{
			name:       "hello with emoji",
			args:       []string{"--emoji"},
			wantOutput: "👋 Hello, World!",
		},
		{
			name:       "hello with json",
			args:       []string{"--json"},
			wantOutput: `"message"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			output := strings.TrimSpace(buf.String())
			if !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Execute() output = %q, want substring %q", output, tt.wantOutput)
			}
		})
	}
}
