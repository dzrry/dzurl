package transport

import (
	"bytes"
	"context"
	"github.com/dzrry/dzurl/domain"
	"github.com/dzrry/dzurl/mocks"
	jsonn "github.com/dzrry/dzurl/serialization/json"
	msgpackk "github.com/dzrry/dzurl/serialization/msgpack"
	"github.com/dzrry/dzurl/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	rdct := &domain.Redirect{
		Key:       "avito-tech",
		URL:       "https://start.avito.ru/tech",
		CreatedAt: time.Now().Unix(),
	}

	t.Run("LoadRedirect with invalid key", func(t *testing.T) {
		srvc := mocks.RedirectService{}
		srvc.On("Load", "invalid-avito-key").Return(nil, service.ErrRedirectNotFound)
		handler := NewHandler(&srvc)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/invalid-avito-key", nil)
		cctx := chi.NewRouteContext()
		cctx.URLParams.Add("key", "invalid-avito-key")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, cctx))
		handler.LoadRedirect(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("LoadRedirect with valid key", func(t *testing.T) {
		srvc := mocks.RedirectService{}
		srvc.On("Load", "avito-tech").Return(rdct, nil)
		handler := NewHandler(&srvc)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/avito-tech", nil)
		cctx := chi.NewRouteContext()
		cctx.URLParams.Add("key", "avito-tech")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, cctx))
		handler.LoadRedirect(rec, req)
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
		assert.Equal(t, rdct.URL, rec.Header().Get("Location"))
	})

	t.Run("StoreRedirect with json", func(t *testing.T) {
		serializer := jsonn.Redirect{}
		body, _ := serializer.Encode(rdct)
		srvc := mocks.RedirectService{}
		srvc.On("Store", rdct).Return(nil)
		handler := NewHandler(&srvc)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		cctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, cctx))
		req.Header.Set("Content-Type", "application/json")
		handler.StoreRedirect(rec, req)
		resp, _ := ioutil.ReadAll(rec.Body)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, body, resp)
	})

	t.Run("StoreRedirect with msgpack", func(t *testing.T) {
		serializer := msgpackk.Redirect{}
		body, _ := serializer.Encode(rdct)
		srvc := mocks.RedirectService{}
		srvc.On("Store", rdct).Return(nil)
		handler := NewHandler(&srvc)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		cctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, cctx))
		req.Header.Set("Content-Type", "application/x-msgpack")
		handler.StoreRedirect(rec, req)
		resp, _ := ioutil.ReadAll(rec.Body)
		assert.Equal(t, "application/x-msgpack", rec.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, body, resp)
	})

	t.Run("StoreRedirect with wrong content type", func(t *testing.T) {
		serializer := msgpackk.Redirect{}
		body, _ := serializer.Encode(rdct)
		srvc := mocks.RedirectService{}
		srvc.On("Store", rdct).Return(nil)
		handler := NewHandler(&srvc)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		cctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, cctx))
		req.Header.Set("Content-Type", "application/json")
		handler.StoreRedirect(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
