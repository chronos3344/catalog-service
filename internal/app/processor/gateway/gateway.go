package pgateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/chronos3344/catalog-service/internal/app/processor"
	catalogv1 "github.com/chronos3344/catalog-service/internal/pkg/grpc/gen/catalog/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gatewayReadHeaderTimeout = 5 * time.Second
	gatewayReadTimeout       = 30 * time.Second
	gatewayWriteTimeout      = 30 * time.Second
	gatewayIdleTimeout       = 120 * time.Second
	gatewayShutdownTimeout   = 5 * time.Second
)

type gatewayProc struct {
	server   http.Server
	addr     string
	grpcAddr string
}

func NewGateway(
	cfgGateway section.ProcessorGateway,
	cfgGrpc section.ProcessorGrpc,
) processor.Processor {
	return &gatewayProc{
		server: http.Server{
			ReadHeaderTimeout: gatewayReadHeaderTimeout,
			ReadTimeout:       gatewayReadTimeout,
			WriteTimeout:      gatewayWriteTimeout,
			IdleTimeout:       gatewayIdleTimeout,
		},
		addr:     fmt.Sprintf(":%d", cfgGateway.ListenPort),
		grpcAddr: net.JoinHostPort("localhost", strconv.Itoa(int(cfgGrpc.ListenPort))),
	}
}

func (p *gatewayProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := catalogv1.RegisterCatalogServiceHandlerFromEndpoint(ctx, mux, p.grpcAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("Failed to register gRPC-Gateway handler")
		return
	}

	p.server.Addr = p.addr
	p.server.Handler = mux

	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create gateway listener")
		return
	}

	log.Info().Str("addr", p.addr).Str("grpc_addr", p.grpcAddr).Msg("gRPC-Gateway started successfully")

	go func() {
		if err := p.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("gRPC-Gateway serve error")
		}
	}()

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(listener.Close))
	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(
		p.server.Shutdown, context.Background(), gatewayShutdownTimeout,
	))
}
