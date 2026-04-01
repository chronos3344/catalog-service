package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product,alias:p"`

	ID           int64     `bun:"id,autoincrement" json:"id"`
	GUID         uuid.UUID `bun:"guid,type:uuid,notnull,pk" json:"guid"`
	Name         string    `bun:"name,notnull" json:"name"`
	Description  *string   `bun:"description" json:"description"`
	Price        float64   `bun:"price,type:decimal(12,2),notnull" json:"price"`
	CategoryGUID uuid.UUID `bun:"category_guid,type:uuid,notnull" json:"category_guid"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

// RequestProductCreate - модель для создания продукта
type RequestProductCreate struct {
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  *string   `json:"description"`
}

func (r RequestProductCreate) Validate() error {
	if r.Name == "" || r.Price <= 0 || r.CategoryGUID == uuid.Nil {
		return ErrIncorrectParameters
	}
	return nil
}

type RequestProductUpdate struct {
	Name         *string    `json:"name"`
	Price        *float64   `json:"price"`
	CategoryGUID *uuid.UUID `json:"category_guid"`
	Description  *string    `json:"description"`
}

func (r RequestProductUpdate) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return ErrIncorrectParameters
	}
	if r.Price != nil && *r.Price <= 0 {
		return ErrIncorrectParameters
	}
	if r.CategoryGUID != nil && *r.CategoryGUID == uuid.Nil {
		return ErrIncorrectParameters
	}
	return nil
}

type ResponseProductCreate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResponseProductGet struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResponseProductUpdate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        float64   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResponseProductList []ResponseProductGet
