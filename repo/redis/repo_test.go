package redis

import (
	"github.com/dzrry/dzurl/config"
	"github.com/dzrry/dzurl/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisRepo(t *testing.T) {
	rdct := &domain.Redirect{
		Key:       "avito-krsk",
		URL:       "www.avito.ru/krasnoyarsk",
		CreatedAt: time.Now().Unix(),
	}
	testCfg := &config.RedisConfig{
		Addr: "localhost",
		Port: "6379",
		Password: ""}
	rr, err := NewRepo(testCfg)
	assert.Nil(t, err)

	t.Run("Store", func(t *testing.T) {
		err := rr.Store(rdct)
		assert.Nil(t, err)
	})

	t.Run("Load", func(t *testing.T) {
		res, err := rr.Load(rdct.Key)
		assert.Nil(t, err)
		assert.Equal(t, res, rdct)
	})

	t.Run("Load non-existent redirect", func(t *testing.T) {
		_, err := rr.Load("non-existent key")
		assert.Equal(t, "repository.Redirect.Load: Redirect not found", err.Error())
	})
}
