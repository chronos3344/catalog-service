package rprocessor

// Обратите внимание на название пакета
// Приставка r (m,p) будет добавляться во многих других пакетах.

import (
	"fmt"
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
	r.HandleFunc("/health", hHealth.LastCheck).Methods("GET")

	// здесь будет регистрация остальных хэндлеров
	// TODO: добавить регистрацию продуктов, категорий и т.д.
	//r.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET")
	//r.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET")

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
	addr := fmt.Sprintf(":%d", cfg.ListenPort)

	s := &httpProc{
		server: http.Server{
			Addr:    addr,
			Handler: r,
		},
		addr: addr,
		//router: r,
	}

	log.Printf("HTTP server configured on %s", addr)
	return &s
}

// Start запускает HTTP сервер
func (h *httpProc) Start() error {
	log.Printf("Starting HTTP server on %s", h.addr)
	return h.server.ListenAndServe()
}
