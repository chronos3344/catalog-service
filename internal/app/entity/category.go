package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories, alias:c" json:"bun_._base_model"`

	ID        uint32    `bun:"id,pk,autoincrement" json:"id,omitempty"`
	Guid      uuid.UUID `bun:"guid,pk,type:uuid_generate_v4()" json:"guid,omitempty"`
	Name      string    `bun:"name" json:"name,omitempty"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at" json:"updated_at"`
}

// RequestCategoryCreate - модель для создания категории
type RequestCategoryCreate struct {
	Name string `json:"name"`
}

// RequestCategoryUpdate - модель для обновления категории
type RequestCategoryUpdate struct {
	Name string `json:"name"`
}

// ResponseCategory - модель для ответа с данными категории
type ResponseCategory struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// ResponseCategoryList - модель для списка категорий
type ResponseCategoryList struct {
	Items      []ResponseCategory `json:"items"`
	TotalCount int                `json:"total_count"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

