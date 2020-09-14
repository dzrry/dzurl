package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	file := "config_test.yml"
	cfg, err := Read(file)
	assert.Nil(t, err)

	assert.Equal(t,"redis-addr", cfg.Redis.Addr)
	assert.Equal(t,"6379", cfg.Redis.Port)
	assert.Equal(t,"", cfg.Redis.Password)

	assert.Equal(t, "server-addr", cfg.Server.Addr)
	assert.Equal(t, "8080", cfg.Server.Port)
}
