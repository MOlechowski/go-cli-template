package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		wantOut   string
		notWantOut string
	}{
		{
			name:    "basic run",
			args:    []string{"example"},
			wantErr: false,
			wantOut: "Running example command",
		},
		{
			name:    "with flag",
			args:    []string{"example", "--flag=test"},
			wantErr: false,
			wantOut: "Flag value: test",
		},
		{
			name:    "with boolean flag",
			args:    []string{"example", "--enable"},
			wantErr: false,
			wantOut: "Boolean flag is enabled",
		},
		{
			name:    "with argument",
			args:    []string{"example", "myarg"},
			wantErr: false,
			wantOut: "Argument provided: myarg",
		},
		{
			name:    "too many arguments",
			args:    []string{"example", "arg1", "arg2"},
			wantErr: true,
		},
		{
			name:       "without boolean flag",
			args:       []string{"example"},
			wantErr:    false,
			notWantOut: "Boolean flag is enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			output := buf.String()
			if tt.wantOut != "" {
				assert.Contains(t, output, tt.wantOut)
			}
			if tt.notWantOut != "" {
				assert.NotContains(t, output, tt.notWantOut)
			}
		})
	}
}

func TestExampleCommandFlags(t *testing.T) {
	// Reset flags for testing
	exampleFlag = ""
	exampleBool = false

	cmd := exampleCmd
	
	// Test flag parsing
	cmd.ParseFlags([]string{"--flag", "value", "--enable"})
	
	assert.Equal(t, "value", exampleFlag)
	assert.True(t, exampleBool)
}