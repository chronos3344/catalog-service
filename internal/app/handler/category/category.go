package hcategory

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/handler"
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

	if err := req.Validate(); err != nil {
		http.Error(w, `{"error":"incorrect parameters"}`, http.StatusBadRequest)
		return
	}

	category, err := h.serviceCategory.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			http.Error(w, `{"error":"Category with this name already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Конвертируем entity в response DTO
	resp := entity.ResponseCategoryCreate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	category, err := h.serviceCategory.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Конвертируем entity в response DTO
	resp := entity.ResponseCategoryGet{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.serviceCategory.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Конвертируем []entity.Category в entity.ResponseCategoryList
	resp := make(entity.ResponseCategoryList, 0, len(categories))
	for _, cat := range categories {
		resp = append(resp, entity.ResponseCategoryGet{
			GUID:      cat.GUID,
			Name:      cat.Name,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

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

	if err := req.Validate(); err != nil {
		http.Error(w, `{"error":"incorrect parameters"}`, http.StatusBadRequest)
		return
	}

	category, err := h.serviceCategory.Update(r.Context(), guid, req.Name)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, entity.ErrAlreadyExists) {
			http.Error(w, `{"error":"Category with this name already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Конвертируем entity в response DTO
	resp := entity.ResponseCategoryUpdate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid UUID format"}`, http.StatusBadRequest)
		return
	}

	err = h.serviceCategory.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, entity.ErrCategoryHasProducts) {
			http.Error(w, `{"error":"category has linked products"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
