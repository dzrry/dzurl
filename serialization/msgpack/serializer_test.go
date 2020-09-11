package msgpack

import (
	"github.com/dzrry/dzurl/domain"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
	"testing"
	"time"
)

func TestMsgpackSerializer(t *testing.T) {
	serializer := Redirect{}

	t.Run("Decode invalid data", func(t *testing.T) {
		invalidRaw := []byte("raw-invalid-message-for-avito-tech")
		_, err := serializer.Decode(invalidRaw)
		assert.Equal(t, "serializer.Msgpack.Decode: msgpack: invalid code=72 decoding map length", err.Error())
	})

	t.Run("Decode valid data", func(t *testing.T) {
		rdct := &domain.Redirect{
			Key:       "avito-tech",
			URL:       "https://start.avito.ru/tech",
			CreatedAt: time.Now().Unix(),
		}
		raw, err := msgpack.Marshal(rdct)
		assert.Nil(t, err)
		res, err := serializer.Decode(raw)
		assert.Nil(t, err)
		assert.Equal(t, rdct, res)
	})

	t.Run("Encode valid value", func(t *testing.T) {
		rdct := &domain.Redirect{
			Key:       "avito-tech",
			URL:       "https://start.avito.ru/tech",
			CreatedAt: time.Now().Unix(),
		}
		raw, err := msgpack.Marshal(rdct)
		assert.Nil(t, err)
		res, err := serializer.Encode(rdct)
		assert.Equal(t, res, raw)
	})
}

