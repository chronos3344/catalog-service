package ghcatalogv1

import (
	"context"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/mapper"
	mcatv1 "github.com/chronos3344/catalog-service/internal/app/mapper/catalog/v1"
	"github.com/chronos3344/catalog-service/internal/app/service"
	catalogv1 "github.com/chronos3344/catalog-service/internal/pkg/grpc/gen/catalog/v1"
	"github.com/google/uuid"
)

type handler struct {
	catalogv1.UnimplementedCatalogServiceServer
	srv service.Product
}

func NewHandler(srv service.Product) catalogv1.CatalogServiceServer {
	return &handler{srv: srv}
}

func (h *handler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
	guid, err := uuid.Parse(req.GetGuid())
	if err != nil {
		return nil, mapper.ErrorToGRPC(entity.ErrIncorrectParameters)
	}

	products, err := h.srv.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return nil, mapper.ErrorToGRPC(err)
	}
	if len(products) == 0 {
		return nil, mapper.ErrorToGRPC(entity.ErrNotFound)
	}

	return &catalogv1.GetProductResponse{
		Product: mcatv1.ProductToProto(products[0]),
	}, nil
}

func (h *handler) GetProducts(ctx context.Context, req *catalogv1.GetProductsRequest) (*catalogv1.GetProductsResponse, error) {
	rawGuids := req.GetGuids()
	if len(rawGuids) == 0 {
		return &catalogv1.GetProductsResponse{}, nil
	}

	guids := make([]uuid.UUID, 0, len(rawGuids))
	for _, raw := range rawGuids {
		guid, err := uuid.Parse(raw)
		if err != nil {
			return nil, mapper.ErrorToGRPC(entity.ErrIncorrectParameters)
		}
		guids = append(guids, guid)
	}

	products, err := h.srv.GetByGUIDs(ctx, guids)
	if err != nil {
		return nil, mapper.ErrorToGRPC(err)
	}

	resp := &catalogv1.GetProductsResponse{
		Products: make([]*catalogv1.Product, 0, len(products)),
	}
	found := make(map[uuid.UUID]struct{}, len(products))
	for _, p := range products {
		resp.Products = append(resp.Products, mcatv1.ProductToProto(p))
		found[p.GUID] = struct{}{}
	}

	resp.MissingGuids = make([]string, 0)
	for _, guid := range guids {
		if _, ok := found[guid]; !ok {
			resp.MissingGuids = append(resp.MissingGuids, guid.String())
		}
	}

	return resp, nil
}
