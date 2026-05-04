package hcategory

import (
	"errors"
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
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	category, err := h.serviceCategory.Create(r.Context(), req.Name)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.ErrorApply(w, http.StatusBadRequest, "Категория с таким именем уже существует")
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, "Ошибка сервера")
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

	httph.SendEncoded(w, r, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		return
	}

	category, err := h.serviceCategory.Get(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			return
		}
		httph.ErrorApply(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	// Конвертируем entity в response DTO
	resp := entity.ResponseCategoryGet{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.serviceCategory.List(r.Context())
	if err != nil {
		httph.ErrorApply(w, http.StatusInternalServerError, "Ошибка сервера")
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

	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		return
	}

	var req entity.RequestCategoryUpdate
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	category, err := h.serviceCategory.Update(r.Context(), guid, req.Name)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			return
		}
		if errors.Is(err, entity.ErrAlreadyExists) {
			httph.ErrorApply(w, http.StatusConflict, "Категория с таким названием уже существует")
			return
		}
		httph.ErrorApply(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	// Конвертируем entity в response DTO
	resp := entity.ResponseCategoryUpdate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guidStr := vars["guid"]

	guid, err := uuid.Parse(guidStr)
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, "Неверный формат UUID")
		return
	}

	err = h.serviceCategory.Delete(r.Context(), guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			httph.ErrorApply(w, http.StatusNotFound, "Категория не найдена")
			return
		}
		if errors.Is(err, entity.ErrCategoryHasProducts) {
			httph.ErrorApply(w, http.StatusConflict, "Нельзя удалить категорию с товарами")
			return
		}
		httph.ErrorApply(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
