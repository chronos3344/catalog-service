package rhealth

import (
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	// Добавьте в хэадер http.StatusOK
	w.WriteHeader(http.StatusOK)
	// Добавьте в body ответа на запрос - "ok"
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
