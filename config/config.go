package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type HTTPServer struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HTTPClient struct {
	Timeout time.Duration
}

type MySQL struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Location string
}

// DSN returns the MySQL DSN.
func (m MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s", m.Username, m.Password, m.Host, m.Port, m.Database, url.QueryEscape(m.Location))
}

type Redis struct {
	Addr     string
	Password string
}

func LoadWithPrefix(prefix string, v interface{}) error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(v); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
