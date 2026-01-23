package rprocessor

import (
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"

	"net/http"

	"github.com/gorilla/mux"
)

func vGenericRegHealthCheck(r *mux.Router, h handler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
