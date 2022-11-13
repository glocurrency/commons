package config_test

import (
	"testing"

	"github.com/glocurrency/commons/config"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer_NewServer(t *testing.T) {
	t.Parallel()

	cfg := config.HTTPServer{
		Addr:         ":8080",
		ReadTimeout:  10,
		WriteTimeout: 10,
		IdleTimeout:  10,
	}

	srv := cfg.NewServer(nil)

	assert.Equal(t, cfg.Addr, srv.Addr)
	assert.Equal(t, cfg.ReadTimeout, srv.ReadTimeout)
	assert.Equal(t, cfg.WriteTimeout, srv.WriteTimeout)
	assert.Equal(t, cfg.IdleTimeout, srv.IdleTimeout)
}

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
