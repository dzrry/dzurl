package serialization

import "github.com/dzrry/dzurl/domain"

type RedirectSerializer interface {
	Decode(data []byte) (*domain.Redirect, error)
	Encode(value *domain.Redirect) ([]byte, error)
}
