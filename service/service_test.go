package service

import (
	"errors"
	"fmt"
	"github.com/dzrry/dzurl/domain"
	"testing"
	"time"

	"github.com/dzrry/dzurl/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	t.Run("Store", func(t *testing.T) {
		rdct := &domain.Redirect{URL: "https://start.avito.ru/tech"}
		rdctRepo := mocks.RedirectRepo{}
		rdctRepo.On("Store", rdct).Return(nil)
		srvc := NewRedirectService(&rdctRepo)
		err := srvc.Store(rdct)
		assert.Nil(t, err)
		assert.NotEmpty(t, rdct.Key)
		assert.NotEmpty(t, rdct.CreatedAt)
	})

	t.Run("Store with validation error", func(t *testing.T) {
		rdct := &domain.Redirect{URL: "invalid-avito-url"}
		rdctRepo := mocks.RedirectRepo{}
		service := NewRedirectService(&rdctRepo)
		err := service.Store(rdct)
		assert.Equal(t, "service.Redirect.Store: url: invalid-avito-url does not validate as requrl", err.Error())
	})

	t.Run("Store with repository error", func(t *testing.T) {
		rdct := &domain.Redirect{URL: "https://start.avito.ru/tech"}
		rdctRepo := mocks.RedirectRepo{}
		rdctRepo.On("Store", rdct).Return(errors.New("Repository error"))
		srvc := NewRedirectService(&rdctRepo)
		err := srvc.Store(rdct)
		assert.Equal(t, "Repository error", err.Error())
	})

	t.Run("Load", func(t *testing.T) {
		rdct := &domain.Redirect{
			Key:       "avito-tech",
			URL:       "https://start.avito.ru/tech",
			CreatedAt: time.Now().Unix(),
		}
		rdctRepo := mocks.RedirectRepo{}
		rdctRepo.On("Load", "avito-tech").Return(rdct, nil)
		srvc := NewRedirectService(&rdctRepo)
		res, err := srvc.Load("avito-tech")
		assert.Nil(t, err)
		assert.Equal(t, rdct, res)
	})

	t.Run("Find with invalid code", func(t *testing.T) {
		rdctRepo := mocks.RedirectRepo{}
		rdctRepo.On("Load", "invalid-avito-key").Return(nil,
			fmt.Errorf("repository.MockedRedirectRepository.Load: %w", ErrRedirectNotFound))
		srvc := NewRedirectService(&rdctRepo)
		_, err := srvc.Load("invalid-avito-key")
		assert.Equal(t, "repository.MockedRedirectRepository.Load: Redirect not found", err.Error())
	})
}
