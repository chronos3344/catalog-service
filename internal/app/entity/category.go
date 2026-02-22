package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories,alias:c"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	GUID      uuid.UUID `bun:"guid,type:uuid,pk,default:gen_random_uuid()" json:"guid"`
	Name      string    `bun:"name,notnull,unique" json:"name"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

type RequestCategoryCreate struct {
	Name string `json:"name"`
}

type ResponseCategoryCreate struct {
	GUID uuid.UUID `json:"guid"`
	Name string    `json:"name"`
}

type ResponseCategoryGet struct {
	GUID uuid.UUID `json:"guid"`
	Name string    `json:"name"`
}

type RequestCategoryUpdate struct {
	Name string `json:"name"`
}

type ResponseCategoryUpdate struct {
	GUID uuid.UUID `json:"guid"`
	Name string    `json:"name"`
}

type ResponseCategoryList []ResponseCategoryGet
