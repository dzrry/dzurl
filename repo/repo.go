package repo

import (
	"github.com/dzrry/dzurl/domain"
)

type RedirectRepo interface {
	Load(key string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}
