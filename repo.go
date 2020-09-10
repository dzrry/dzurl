package main

type RedirectRepo interface {
	Load(key string) (*Redirect, error)
	Store(redirect *Redirect) error
}
