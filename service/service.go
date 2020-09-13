package service

import (
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/dzrry/dzurl/domain"
	"github.com/dzrry/dzurl/repo"
	"github.com/rs/xid"
	"time"
)

type RedirectService interface {
	Load(key string) (*domain.Redirect, error)
	Store(redirect *domain.Redirect) error
}

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
	_, err := valid.ValidateStruct(redirect)
	if err != nil {
		return fmt.Errorf("service.Redirect.Store: %w", err)
	}
	if redirect.Key == "" {
		redirect.Key = xid.New().String()
	}
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}
