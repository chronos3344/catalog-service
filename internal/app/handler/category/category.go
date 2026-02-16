package hcategory

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	serviceCategory service.Category
}

func NewHandler(serviceCategory service.Category) rhandler.Category {
	return &handler{serviceCategory}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestCategoryCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Name) > 255 {
		http.Error(w, `{"error":"Name too long"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceCategory.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, `{"error":"Category with this name already exists"}`, http.StatusConflict)
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
	guidStr := vars["category_guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceCategory.Get(r.Context(), guid)
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

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	resp, err := h.serviceCategory.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["category_guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	var req entity.RequestCategoryUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Name) > 255 {
		http.Error(w, `{"error":"Name too long"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.serviceCategory.Update(r.Context(), guid, req.Name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrConflict) {
			http.Error(w, `{"error":"Category with this name already exists"}`, http.StatusConflict)
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
	guidStr := vars["category_guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	err = h.serviceCategory.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
