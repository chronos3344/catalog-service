package rprocessor

import (
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write([]byte("not found")); err != nil {
		log.Info().Msg("Failed to write response")
	}
}
