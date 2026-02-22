package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	GUID         uuid.UUID `bun:"guid,type:uuid,pk,default:gen_random_uuid()" json:"guid"`
	Name         string    `bun:"name,notnull,unique" json:"name"`
	Description  *string   `bun:"description" json:"description"`
	Price        float64   `bun:"price,type:decimal(12,3),notnull" json:"price"`
	CategoryGUID uuid.UUID `bun:"category_guid,type:uuid,notnull" json:"category_guid"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	Category     *Category `bun:"rel:belongs_to,join:category_guid=guid" json:"category,omitempty"`
}

// RequestProductCreate - модель для создания продукта
type RequestProductCreate struct {
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  *string   `json:"description"`
}

type ResponseProductCreate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  *string   `json:"description"`
}

type ResponseProductGet struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  *string   `json:"description"`
}

type RequestProductUpdate struct {
	Name         *string    `json:"name"`
	Price        *float64   `json:"price"`
	CategoryGUID *uuid.UUID `json:"category_guid"`
	Description  *string    `json:"description"`
}

type ResponseProductUpdate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  *string   `json:"description"`
}

type RequestProductList struct {
	CategoryGUID *uuid.UUID `json:"category_guid"`
	MinPrice     *float64   `json:"min_price"`
	MaxPrice     *float64   `json:"max_price"`
}

type ResponseProductList []ResponseProductGet
