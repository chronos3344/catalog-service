package rprocessor

import (
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	// Используем нашу вспомогательную функцию функцию reg().
	// Еще раз внимательно посмотрите, что мы в нее передаем.
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	// Передаем в заголовок http.StatusNotFound
}
