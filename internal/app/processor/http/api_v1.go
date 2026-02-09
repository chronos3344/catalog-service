package rprocessor
//
//import (
//	"net/http"
//
//	"github.com/gorilla/mux"
//
//	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
//)
//
//func v1RegCategoryHandler(r1 *mux.Router, h rhandler.Category) {
//	reg(r1, http.MethodGet, "/category/create", ...)  // Делаем по аналогии с обработчиком health
//	reg(r1, ..., "/category/{category_guid}", ...)
//	reg(r1, ..., "/category/list", ...)
//}
//
//func v1RegProductHandler(r1 *mux.Router, h rhandler.Product) {
//	reg(r1, ..., "/product/create", ...)
//	...
//	...
//}