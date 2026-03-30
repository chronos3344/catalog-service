package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:category,alias:c"`

	ID        int64     `bun:"id,autoincrement" json:"id"`
	GUID      uuid.UUID `bun:"guid,type:uuid" json:"guid"`
	Name      string    `bun:"name,notnull,unique" json:"name"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type ResponseCategoryCreate struct {
	GUID uuid.UUID `json:"guid"`
	Name string    `json:"name"`
}

type ResponseCategoryGet struct {
	GUID uuid.UUID `json:"guid"`
	Name string    `json:"name"`
}

type ResponseCategoryUpdate struct {
	GUID uuid.UUID `json:"guid"`
	Name string    `json:"name"`
}

type ResponseCategoryList []ResponseCategoryGet

////////////////////////////////////////////////////////////////////////////////
///// HTTP REQUEST & RESPONSE //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type RequestCategoryCreate struct {
	Name string `json:"name"`
}

func (r RequestCategoryCreate) Validate() error {
	if r.Name == "" {
		return ErrIncorrectParameters
	}
	return nil
}

type RequestCategoryUpdate struct {
	Name string `json:"name"`
}

func (r RequestCategoryUpdate) Validate() error {
	if r.Name == "" {
		return ErrIncorrectParameters
	}
	return nil
}

type ResponseCategory struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
