package json

import (
	"encoding/json"
	"fmt"
	"github.com/dzrry/dzurl/domain"
)

type Redirect struct {}

func (r *Redirect) Decode(d []byte) (*domain.Redirect, error) {
	rdct := &domain.Redirect{}
	if err := json.Unmarshal(d, rdct); err != nil {
		return rdct, fmt.Errorf("serializer.Json.Decode: %w", err)
	}
	return rdct, nil
}

func (r *Redirect) Encode(v *domain.Redirect) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return raw, fmt.Errorf("serializer.Json.Encode: %w", err)
	}
	return raw, nil
}
