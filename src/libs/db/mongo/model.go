package mongo

import (
	"time"
)

type Model struct {
	ID        *ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt *time.Time `bson:"created_at,omitempty" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty" json:"updatedAt,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}
