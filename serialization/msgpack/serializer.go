package msgpack

import (
	"fmt"
	"github.com/dzrry/dzurl/domain"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(d []byte) (*domain.Redirect, error) {
	rdct := &domain.Redirect{}
	if err := msgpack.Unmarshal(d, rdct); err != nil {
		return nil, fmt.Errorf("serializer.Msgpack.Decode: %w", err)
	}
	return rdct, nil
}

func (r *Redirect) Encode(v *domain.Redirect) ([]byte, error) {
	raw, err := msgpack.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("serializer.Msgpack.Encode: %w", err)
	}
	return raw, nil
}
