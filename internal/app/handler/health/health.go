package rhealth // Обратите внимание на именование пакета!

import (
	"net/http"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	// Добавьте в хэадер http.StatusOK
	w.WriteHeader(http.StatusOK)
	// Добавьте в body ответа на запрос - "ok"
	fmt.Fprint(w, "OK")
}
