package mcategory

import (
	"context"

	"catalog-service/internal/app/repository"
)

type (
	service struct {
		repoCategory repository.Category
	}
)

// Наша функция-билдер
func NewService(repoCategory repository.Category) service.Category {
	return &service{
		repoCategory: repoCategory,
	}
}

// Опять же реализовываем наш интерфейс, здесь уже расписываем логику
func (s *service) Create(ctx context.Context, name string, price float64, category_guid uuid.UUID, description string) (entity.ResponseCategoryCreate, error) {
	panic("implement me!")

	// Здесь распологается основная логика.
	// Например мы хотим создать товар, который уже существует, мы должные вернуть ошибку в таком случае.
	// Мы например можем проверить базу данных на дубликаты, по параметрам которые нам пришли в запросе
	// или попытаться вставить данные, как есть, но получить ошибку при попытке обработки от БД. На ваш выбор.

	// Пока можно не следовать формату возврата ошибок из описания API, но основные кейсы должны обрабатываться корректно.
}
