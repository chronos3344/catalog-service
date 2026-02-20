package rprocessor

import (
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
)

func v1RegCategoryHandler(r1 *mux.Router, h rhandler.Category) {
	reg(r1, "POST", "/category/create", h.Create)
	reg(r1, "GET", "/category/{category_guid}", h.Get)
	reg(r1, "POST", "/category/list", h.List)
	reg(r1, "PUT", "/category/{category_guid}", h.Update)
	reg(r1, "DELETE", "/category/{category_guid}", h.Delete)
}

func v1RegProductHandler(r1 *mux.Router, h rhandler.Product) {
	reg(r1, "POST", "/product/create", http.HandlerFunc(h.Create))
	reg(r1, "GET", "/product/{product_guid}", h.Get)
	reg(r1, "POST", "/product/list", h.List)
	reg(r1, "PUT", "/product/{product_guid}", h.Update)
	reg(r1, "DELETE", "/product/{product_guid}", h.Delete)
}
