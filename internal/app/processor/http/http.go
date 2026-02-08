package rprocessor

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
)

type httpProc struct {
	server *http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product,
	cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	vGenericRegHealthCheck(r, hHealth)

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Registered route: %s %s", methods, pathTemplate)
		return nil
	})

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

func (h *httpProc) Serve() error {
	log.Printf("Starting HTTP server on %s", h.addr)
	return h.server.ListenAndServe()
}
