package rprocessor

import (
	"github.com/gorilla/mux"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
)

func v1RegCategoryHandler(r1 *mux.Router, h rhandler.Category) {
	reg(r1, "/category/create", h.Create, "POST")
	reg(r1, "/category/{category_guid}", h.Get, "GET")
	reg(r1, "/category/list", h.List, "POST")
	reg(r1, "/category/{category_guid}", h.Update, "PUT")
	reg(r1, "/category/{category_guid}", h.Delete, "DELETE")
}

func v1RegProductHandler(r1 *mux.Router, h rhandler.Product) {
	reg(r1, "/product/create", h.Create, "POST")
	reg(r1, "/product/{product_guid}", h.Get, "GET")
	reg(r1, "/product/list", h.List, "POST")
	reg(r1, "/product/{product_guid}", h.Update, "PUT")
	reg(r1, "/product/{product_guid}", h.Delete, "DELETE")
}
