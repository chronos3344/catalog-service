package rprocessor

import (
	"log"
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"

	"github.com/gorilla/mux"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", h.LastCheck)
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write([]byte("not found")); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
