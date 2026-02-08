package category

import (
	"net/http"

	rhandler "github.com/chronos3344/catalog-service/internal/app/handler"
)

type handler struct {
	serviceCategory service.Category
}

func NewHandler(serviceCategory service.Category) rhandler.Category {
	return &handler{serviceCategory}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	// Здесь будет располагаться логика обработки запроса

	// Напоминание!
	// Мы конвертируем запрос (r *http.Request) в нашу структуру запроса entity.RequestCategoryCreate.
	// Пока делаем это сами, позже на этом этапе будет реализована валидация.

	// Теперь конвертируем entity.RequestCategoryCreate в наш entity.Category
	// и только теперь вызываем метод сервиса для обработки.
}
