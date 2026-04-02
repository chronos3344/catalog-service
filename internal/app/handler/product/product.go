package hproduct

import (
	"encoding/json"
	"errors"
	"github.com/chronos3344/catalog-service/pkg/http/httph"
	"log"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	serviceProduct service.Product
}

func NewHandler(serviceProduct service.Product) rhandler.Product {
	return &handler{serviceProduct}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestProductCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	if err := req.Validate(); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	product, err := h.serviceProduct.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			httph.SendError(w, http.StatusNotFound, entity.ErrNotFound)
			return
		}
		if errors.Is(err, entity.ErrAlreadyExists) {
			httph.SendError(w, http.StatusConflict, err)
			httph.ErrorApply(w, http.StatusConflict, "Товар с таким названием уже существует")
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	resp := entity.ResponseProductCreate{
		GUID:         product.GUID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	product, err := h.serviceProduct.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Продукт не найден")
			httph.SendError(w, http.StatusNotFound, err)
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	resp := entity.ResponseProductGet{
		GUID:         product.GUID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.serviceProduct.List(r.Context())
	if err != nil {
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}
	resp := make(entity.ResponseProductList, 0, len(products))
	for _, p := range products {
		resp = append(resp, entity.ResponseProductGet{
			GUID:         p.GUID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CategoryGUID: p.CategoryGUID,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	var req entity.RequestProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	if err := req.Validate(); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	product, err := h.serviceProduct.Update(r.Context(), guid, req)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Продукт не найден")
			httph.SendError(w, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, entity.ErrAlreadyExists) {
			httph.ErrorApply(w, http.StatusConflict, "Товар с таким названием уже существует")
			httph.SendError(w, http.StatusConflict, err)
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	resp := entity.ResponseProductUpdate{
		GUID:         product.GUID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	err = h.serviceProduct.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Продукт не найден")
			httph.SendError(w, http.StatusNotFound, err)
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
