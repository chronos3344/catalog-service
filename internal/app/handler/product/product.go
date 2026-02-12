package hproduct

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/chronos3344/catalog-service/internal/app/service"
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
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	if req.Price <= 0 {
		http.Error(w, `{"error":"Invalid price value"}`, http.StatusBadRequest)
		return
	}

	if len(req.Name) > 255 {
		http.Error(w, `{"error":"Name too long"}`, http.StatusBadRequest)
		return
	}

	if req.CategoryGUID == uuid.Nil {
		http.Error(w, `{"error":"Category GUID is required"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceProduct.Create(r.Context(), req.Name, req.Price, req.CategoryGUID, req.Description)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, `{"error":"Product with this name already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["product_guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceProduct.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Product not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestProductList
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		http.Error(w, `{"error":"Min price cannot be greater than max price"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceProduct.List(r.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["product_guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	var req entity.RequestProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	if req.Price != nil && *req.Price <= 0 {
		http.Error(w, `{"error":"Invalid price value"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceProduct.Update(r.Context(), guid, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Product not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, `{"error":"Product with this name already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["product_guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	err = h.serviceProduct.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Product not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
