package repository

import (
	"context"

	"catalog-service/internal/app/entity"
)

type (
	Category interface {
		Create(ctx context.Context, category entity.Category) (entity.Category, error) // Примерная сигнатура метода
		// ...
		// ...
		// ...
		// ...
	}

	Product interface {
		Create(ctx context.Context, product entity.Product) (entity.Product, error)
		// ...
		// ...
		// ...
	}
)
