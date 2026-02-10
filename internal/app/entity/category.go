package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories, alias:c" json:"bun_._base_model"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id,omitempty"`
	Guid      int64     `bun:"GUID,type:uuid_generate_v4()" json:"guid,omitempty"`
	Name      string    `bun:"name" json:"name,omitempty"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at" json:"updated_at"`
}
