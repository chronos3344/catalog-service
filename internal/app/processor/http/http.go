package rprocessor

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/chronos3344/catalog-service/internal/app/processor"
	"github.com/chronos3344/catalog-service/internal/app/util"
	"github.com/chronos3344/catalog-service/internal/pkg/http/httph"
	"github.com/chronos3344/catalog-service/internal/pkg/http/mzerolog"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server *http.Server
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product, cfg section.ProcessorWebServer) processor.Processor {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	r.Use(
		httph.NewErrorMiddleware(),
		mzerolog.NewMiddleware(
			mzerolog.WithSkipper(util.IsFilteredHttpRoute),
		),
	)

	if hHealth != nil {
		vGenericRegHealthCheck(r, hHealth)
	}

	rV1 := r.PathPrefix("/v1").Subrouter()

	if hCategory != nil {
		v1RegCategoryHandler(rV1, hCategory)
	}
	if hProduct != nil {
		v1RegProductHandler(rV1, hProduct)
	}

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

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", p.server.Addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create listener")
	}

	log.Info().Str("addr", p.server.Addr).Msg("HTTP server started successfully")

	go func() {
		if err := p.serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("HTTP server serve error")
		}
	}()

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(listener.Close))
	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(
		p.server.Shutdown, context.Background(), 5*time.Second,
	))
}

func (p *httpProc) serve(l net.Listener) error {
	log.Info().Msg("Starting HTTP server")
	return p.server.Serve(l)
}
