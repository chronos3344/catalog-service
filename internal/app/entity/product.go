package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products, alias:p" json:"bun_._base_model"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id,omitempty"`
	Guid         uuid.UUID `bun:"guid" json:"guid,omitempty"`
	Name         string    `bun:"name" json:"name,omitempty"`
	Description  string    `bun:"description" json:"description,omitempty"`
	Price        float64   `bun:"price" json:"price,omitempty"`
	CategoryGuid uuid.UUID `bun:"category_guid" json:"category_guid,omitempty"`
	CreatedAt    time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at" json:"updated_at"`
	Category     *Category `bun:"rel:belongs_to,join:category_guid=guid" json:"category,omitempty"`
}
// RequestProductCreate - модель для создания продукта
type RequestProductCreate struct {
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
}

// RequestProductUpdate - модель для обновления продукта
type RequestProductUpdate struct {
	Name         *string    `json:"name,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Price        *float64   `json:"price,omitempty"`
	CategoryGUID *uuid.UUID `json:"category_guid,omitempty"`
}

// ResponseProduct - модель для ответа с данными продукта
type ResponseProduct struct {
	GUID         uuid.UUID           `json:"guid"`
	Name         string              `json:"name"`
	Description  string              `json:"description,omitempty"`
	Price        float64             `json:"price"`
	CategoryGUID uuid.UUID           `json:"category_guid"`
	Category     *ResponseCategory   `json:"category,omitempty"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

// ResponseProductList - модель для списка продуктов
type ResponseProductList struct {
	Items      []ResponseProduct `json:"items"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}