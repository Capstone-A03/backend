package sql

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        *uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	CreatedAt *time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
}
