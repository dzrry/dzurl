package service

import (
	"errors"
	"github.com/dzrry/dzurl/domain"
	"github.com/dzrry/dzurl/repo"
	"github.com/rs/xid"
	"time"
)

type RedirectService interface {
	Load(key string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}

var (
	ErrRedirectNotFound = errors.New("Redirect not found")
	ErrRedirectInvalid  = errors.New("Redirect invalid")
)

type redirectService struct {
	redirectRepo repo.RedirectRepo
}

func NewRedirectService(redirectRepo repo.RedirectRepo) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

func (r *redirectService) Load(key string) (*domain.Redirect, error) {
	return r.redirectRepo.Load(key)
}

func (r *redirectService) Store(redirect *domain.Redirect) error {
	redirect.Key = xid.New().String()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}
