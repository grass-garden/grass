package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.grass.garden/router"
)

func BenchmarkRouter(b *testing.B) {
	r := grassRouter()
	request, _ := http.NewRequest(http.MethodGet, "/bench", nil)
	response := httptest.NewRecorder()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(response, request)
	}
}

func BenchmarkMux(b *testing.B) {
	mux := serveMux()
	request, _ := http.NewRequest(http.MethodGet, "/bench", nil)
	response := httptest.NewRecorder()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(response, request)
	}
}

func grassRouter() *router.Router {
	r := router.New()
	router.Get(r, "/bench", func(*router.ContextAny) (any, error) {
		return nil, nil
	})
	return r
}

func serveMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /bench", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(err)
			}
		}()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nil)
	})
	return mux
}
