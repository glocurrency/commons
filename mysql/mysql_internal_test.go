package mysql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvAsInt(t *testing.T) {
	t.Run("returns default when env var does not exist", func(t *testing.T) {
		key := "TEST_ENV_VAR_NOT_EXIST"
		os.Unsetenv(key)
		val := getEnvAsInt(key, 123)
		assert.Equal(t, 123, val)
	})

	t.Run("returns parsed int when env var exists and is valid", func(t *testing.T) {
		key := "TEST_ENV_VAR_VALID"
		os.Setenv(key, "456")
		defer os.Unsetenv(key)
		val := getEnvAsInt(key, 123)
		assert.Equal(t, 456, val)
	})

	t.Run("returns default when env var exists but is invalid", func(t *testing.T) {
		key := "TEST_ENV_VAR_INVALID"
		os.Setenv(key, "not-an-int")
		defer os.Unsetenv(key)
		val := getEnvAsInt(key, 123)
		assert.Equal(t, 123, val)
	})
}

func TestConnectionLimitEnvVars(t *testing.T) {
	t.Run("resolves MYSQL_DB_MAX_OPEN_CONNS first", func(t *testing.T) {
		os.Setenv("MYSQL_DB_MAX_OPEN_CONNS", "40")
		os.Setenv("DB_MAX_OPEN_CONNS", "30")
		defer os.Unsetenv("MYSQL_DB_MAX_OPEN_CONNS")
		defer os.Unsetenv("DB_MAX_OPEN_CONNS")

		maxOpen := getEnvAsInt("MYSQL_DB_MAX_OPEN_CONNS", getEnvAsInt("DB_MAX_OPEN_CONNS", 20))
		assert.Equal(t, 40, maxOpen)
	})

	t.Run("falls back to DB_MAX_OPEN_CONNS if MYSQL_ prefix is missing", func(t *testing.T) {
		os.Setenv("DB_MAX_OPEN_CONNS", "30")
		defer os.Unsetenv("DB_MAX_OPEN_CONNS")

		maxOpen := getEnvAsInt("MYSQL_DB_MAX_OPEN_CONNS", getEnvAsInt("DB_MAX_OPEN_CONNS", 20))
		assert.Equal(t, 30, maxOpen)
	})

	t.Run("falls back to default if both are missing", func(t *testing.T) {
		maxOpen := getEnvAsInt("MYSQL_DB_MAX_OPEN_CONNS", getEnvAsInt("DB_MAX_OPEN_CONNS", 20))
		assert.Equal(t, 20, maxOpen)
	})

	t.Run("resolves MYSQL_DB_MAX_IDLE_CONNS first", func(t *testing.T) {
		os.Setenv("MYSQL_DB_MAX_IDLE_CONNS", "15")
		os.Setenv("DB_MAX_IDLE_CONNS", "10")
		defer os.Unsetenv("MYSQL_DB_MAX_IDLE_CONNS")
		defer os.Unsetenv("DB_MAX_IDLE_CONNS")

		maxIdle := getEnvAsInt("MYSQL_DB_MAX_IDLE_CONNS", getEnvAsInt("DB_MAX_IDLE_CONNS", 5))
		assert.Equal(t, 15, maxIdle)
	})

	t.Run("falls back to DB_MAX_IDLE_CONNS if MYSQL_ prefix is missing", func(t *testing.T) {
		os.Setenv("DB_MAX_IDLE_CONNS", "10")
		defer os.Unsetenv("DB_MAX_IDLE_CONNS")

		maxIdle := getEnvAsInt("MYSQL_DB_MAX_IDLE_CONNS", getEnvAsInt("DB_MAX_IDLE_CONNS", 5))
		assert.Equal(t, 10, maxIdle)
	})

	t.Run("falls back to default if both are missing for idle conns", func(t *testing.T) {
		maxIdle := getEnvAsInt("MYSQL_DB_MAX_IDLE_CONNS", getEnvAsInt("DB_MAX_IDLE_CONNS", 5))
		assert.Equal(t, 5, maxIdle)
	})
}
