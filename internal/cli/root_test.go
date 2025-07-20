package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// Test that Execute returns no error with help flag
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	os.Args = []string{"cmd", "--help"}
	err := Execute()
	assert.NoError(t, err)
}

func TestInitConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, ".testconfig.yaml")
	
	// Write test config
	configContent := `verbose: true
test_value: "hello"`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)
	
	// Reset viper for test
	viper.Reset()
	
	// Set config file
	cfgFile = configFile
	
	// Run initConfig
	initConfig()
	
	// Check that config was loaded
	assert.True(t, viper.GetBool("verbose"))
	assert.Equal(t, "hello", viper.GetString("test_value"))
	
	// Reset
	cfgFile = ""
	viper.Reset()
}

func TestRootCommandFlags(t *testing.T) {
	// Test that flags are properly defined
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("config"))
	assert.NotNil(t, rootCmd.PersistentFlags().Lookup("verbose"))
	
	// Test short flag
	flag := rootCmd.PersistentFlags().ShorthandLookup("v")
	assert.NotNil(t, flag)
	assert.Equal(t, "verbose", flag.Name)
}