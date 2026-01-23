package rprocessor

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/chronos3344/catalog-service/internal/app/handler"

	"github.com/gorilla/mux"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth handler.Health, cfg section.ProcessorWebServer) *httpProc {
	// Создаем мультиплексор
	r := mux.NewRouter()

	// Настраиваем NotFound handler
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	// Регистрируем health-check хэндлер
	vGenericRegHealthCheck(r, hHealth)

	// TODO: добавить регистрацию остальных хэндлеров

	// Обходим маршруты для дебага
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route registered: %s %s", methods, path)
		return nil
	})
	if err != nil {
		log.Printf("Error walking routes: %v", err)
	}

	// Создаем адрес для сервера
	addr := fmt.Sprintf(":%d", cfg.Port)

	// Инициализируем структуру httpProc
	s := &httpProc{
		server: http.Server{
			Addr:    addr,
			Handler: r,
		},
		addr: addr,
	}

	log.Printf("HTTP server routes registered")

	return s
}

func (h *httpProc) Serve() error {
	log.Printf("Starting HTTP server on %s", h.addr)
	return h.server.ListenAndServe()
}
