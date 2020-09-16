package transport

import (
	"errors"
	"github.com/dzrry/dzurl/serialization"
	jsonn "github.com/dzrry/dzurl/serialization/json"
	msgpackk "github.com/dzrry/dzurl/serialization/msgpack"
	"github.com/dzrry/dzurl/service"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
)

type RedirectHandler interface {
	LoadRedirect(http.ResponseWriter, *http.Request)
	StoreRedirect(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService service.RedirectService
}

func (h *handler) serializer(contentType string) serialization.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpackk.Redirect{}
	}
	return &jsonn.Redirect{}
}

func NewHandler(redirectService service.RedirectService) *handler {
	return &handler{
		redirectService: redirectService,
	}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) LoadRedirect(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	rdct, err := h.redirectService.Load(key)
	if err != nil {
		if errors.Is(err, service.ErrRedirectNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, rdct.URL, http.StatusMovedPermanently)
}

func (h *handler) StoreRedirect(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rdct, err := h.serializer(ct).Decode(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = h.redirectService.Store(rdct); err != nil {
		if errors.Is(err, service.ErrRedirectInvalid) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := h.serializer(ct).Encode(rdct)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, ct, resp, http.StatusCreated)
}
