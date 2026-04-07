package hcategory

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/handler"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/chronos3344/catalog-service/internal/pkg/http/binding"
	"github.com/chronos3344/catalog-service/internal/pkg/http/httph"
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

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	category, err := h.serviceCategory.Create(r.Context(), req.Name)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
			httph.SendError(w, http.StatusBadRequest, err)
		default:
			httph.SendError(w, http.StatusInternalServerError, err)
		}
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
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	category, err := h.serviceCategory.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			httph.SendError(w, http.StatusNotFound, entity.ErrNotFound)
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
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
		httph.SendError(w, http.StatusInternalServerError, err)
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
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	var req entity.RequestCategoryUpdate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	category, err := h.serviceCategory.Update(r.Context(), guid, req.Name)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			httph.SendError(w, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, entity.ErrAlreadyExists) {
			httph.ErrorApply(w, http.StatusConflict, "Категория с таким названием уже существует")
			httph.SendError(w, http.StatusConflict, err)
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
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
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
		return
	}

	err = h.serviceCategory.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			httph.SendError(w, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, entity.ErrCategoryHasProducts) {
			httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
			httph.SendError(w, http.StatusBadRequest, entity.ErrIncorrectParameters)
			return
		}
		httph.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
