package config_test

import (
	"testing"

	"github.com/glocurrency/commons/config"
	"github.com/stretchr/testify/assert"
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
