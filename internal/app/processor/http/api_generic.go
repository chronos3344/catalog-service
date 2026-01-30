package rprocessor

import (
	"log"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	// Используем нашу вспомогательную функцию функцию reg().
	// Еще раз внимательно посмотрите, что мы в нее передаем.
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	// Передаем в заголовок http.StatusNotFound
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write([]byte("not found")); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
