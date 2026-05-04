package rprocessor

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server *http.Server
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	// Регистрируем health check
	if hHealth != nil {
		vGenericRegHealthCheck(r, hHealth)
	}

	// API version 1
	rV1 := r.PathPrefix("/v1").Subrouter()

	if hCategory != nil {
		v1RegCategoryHandler(rV1, hCategory)
	}
	if hProduct != nil {
		v1RegProductHandler(rV1, hProduct)
	}

	// Логируем все зарегистрированные маршруты
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Info().Str("path_template", pathTemplate).Strs("methods", methods).Msg("Registered route")
		return nil
	})

	addr := ":" + strconv.Itoa(cfg.ListenPort)
	log.Info().Int("listen_port", cfg.ListenPort).Msg("Listening on" + addr)

	return &httpProc{
		server: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

func (h *httpProc) Serve() error {
	log.Info().Msg("Starting HTTP server")
	return h.server.ListenAndServe()
}

func (h *httpProc) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server...")
	return h.server.Shutdown(ctx)
}
