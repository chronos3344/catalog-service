package rprocessor

import (
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

func v1RegCategoryHandler(r1 *mux.Router, h rhandler.Category) {
	reg(r1, http.MethodPost, "/product/create", h.Create)
	reg(r1, http.MethodPost, "/category/create", h.Create)
	reg(r1, http.MethodGet, "/category/{category_guid}", h.Get)
	reg(r1, http.MethodPost, "/category/list", h.List)
	reg(r1, http.MethodPut, "/category/{category_guid}", h.Update)
	reg(r1, http.MethodDelete, "/category/{category_guid}", h.Delete)
}

func v1RegProductHandler(r1 *mux.Router, h rhandler.Product) {
	reg(r1, http.MethodPost, "/product/create", h.Create)
	reg(r1, http.MethodGet, "/product/{product_guid}", h.Get)
	reg(r1, http.MethodPost, "/product/list", h.List)
	reg(r1, http.MethodPut, "/product/{product_guid}", h.Update)
	reg(r1, http.MethodDelete, "/product/{product_guid}", h.Delete)
}
