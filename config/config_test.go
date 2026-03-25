package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/glocurrency/commons/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQL_DSN(t *testing.T) {
	t.Parallel()

	mysql := config.MySQL{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "secret",
		Database: "clearfx",
		Location: "Europe/London",
	}

	assert.Equal(t, "root:secret@tcp(localhost:3306)/clearfx?charset=utf8mb4&parseTime=true&loc=Europe%2FLondon", mysql.DSN())
}

func TestLoadWithPrefix_Success(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	tempDir := t.TempDir()
	require.NoError(t, os.Chdir(tempDir))

	defer viper.Reset()

	// ADDED: [httpclient] section so Viper knows this key path exists.
	// The env var will successfully override the "1s" to "5s".
	configContent := `
[httpserver]
addr = ":8080"
readtimeout = "10s"

[httpclient]
timeout = "1s"

[mysql]
host = "127.0.0.1"
port = 3306
`
	err = os.WriteFile("config.toml", []byte(configContent), 0644)
	require.NoError(t, err)

	os.Setenv("TESTAPP_HTTPCLIENT_TIMEOUT", "5s")
	os.Setenv("TESTAPP_MYSQL_PORT", "3307")

	defer os.Unsetenv("TESTAPP_HTTPCLIENT_TIMEOUT")
	defer os.Unsetenv("TESTAPP_MYSQL_PORT")

	type AppConfig struct {
		HTTPServer config.HTTPServer
		HTTPClient config.HTTPClient
		MySQL      config.MySQL
	}

	var cfg AppConfig
	err = config.LoadWithPrefix("TESTAPP", &cfg)
	require.NoError(t, err)

	assert.Equal(t, ":8080", cfg.HTTPServer.Addr)
	assert.Equal(t, 10*time.Second, cfg.HTTPServer.ReadTimeout)
	assert.Equal(t, "127.0.0.1", cfg.MySQL.Host)

	// This will now pass!
	assert.Equal(t, 5*time.Second, cfg.HTTPClient.Timeout)
	assert.Equal(t, 3307, cfg.MySQL.Port)
}

func TestLoadWithPrefix_FileNotFound(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	// Move to an empty temp dir and DO NOT write a config.toml
	tempDir := t.TempDir()
	require.NoError(t, os.Chdir(tempDir))

	defer viper.Reset()

	type AppConfig struct {
		HTTPServer config.HTTPServer
	}

	var cfg AppConfig
	err = config.LoadWithPrefix("TESTAPP", &cfg)

	// viper.ReadInConfig() should fail because the file doesn't exist
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}
