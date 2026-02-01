package rprocessor

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

type httpProc struct {
	server *http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	// создаем мультиплексор
	r := mux.NewRouter()

	// Настраиваем NotFound handler
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	// Регистрируем healthcheck хэндлер
	vGenericRegHealthCheck(r, hHealth)

	//обходим маршруты для дебага через r.Walk
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Registered route: %s %s", methods, pathTemplate)
		return nil
	})
	// если не получится реализовать r.Walk() просто добавляем лог, когда регистрируем маршрут
	//log.Printf("HTTP server routes registered")

	// создаем сервер и возвращаем его
	addr := fmt.Sprintf(":%d", cfg.ListenPort)
	s := &httpProc{
		server: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
		},
		addr: addr,
		//router: r,
	}

	log.Printf("HTTP server configured on %s", addr)
	return s
}

// Start запускает HTTP сервер
func (h *httpProc) Serve() error {
	log.Printf("Starting HTTP server on %s", h.addr)
	return h.server.ListenAndServe()
}
