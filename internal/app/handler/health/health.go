package rhealth

import (
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"

	"net/http"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
