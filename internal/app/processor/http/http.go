package rprocessor

// Обратите внимание на название пакета
// Приставка r (m,p) будет добавляться во многих других пакетах.

import (
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"

	"net/http"

	"github.com/gorilla/mux"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	// создаем мультиплексор
	r := http.NewServeMux()

	// регистрируем HealthCheck
	r.HandleFunc("/health", hHealth.LastCheck).Methods(http.MethodGet)

	// здесь будет регистрация остальных хэндлеров
	// TODO: добавить регистрацию продуктов, категорий и т.д.

	// обходим маршруты для дебага через r.Walk

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("Registered route: %s %s", methods, pathTemplate)
		}
		return nil
	})
	// если не получится реализовать r.Walk() просто добавляем лог, когда регистрируем маршрут

	// создаем сервер и возвращаем его

	return &s
}
