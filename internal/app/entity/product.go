package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products, alias:p" json:"bun_._base_model"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id,omitempty"`
	Guid         int64     `bun:"guid" json:"guid,omitempty"`
	Name         string    `bun:"name" json:"name,omitempty"`
	Description  string    `bun:"description" json:"description,omitempty"`
	Price        float64   `bun:"price" json:"price,omitempty"`
	CategoryGuid string    `bun:"category_guid" json:"category_guid,omitempty"`
	CreatedAt    time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at" json:"updated_at"`
	Category     *Category `bun:"rel:belongs_to,join:category_guid=GUID" json:"category,omitempty"`
}
