package main

import (
	"errors"
	"time"
)

type RedirectService interface {
	Load(key string) (*Redirect, error)
	Store(redirect *Redirect) error
}
var (
	ErrRedirectNotFound = errors.New("Redirect not found")
	ErrRedirectInvalid = errors.New("Redirect invalid")
)

type redirectService struct {
	redirectRepo RedirectRepo
}

func NewRedirectService(redirectRepo RedirectRepo) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

func (r *redirectService) Load(key string) (*Redirect, error) {
	return r.redirectRepo.Load(key)
}

func (r *redirectService) Store(redirect *Redirect) error {
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}