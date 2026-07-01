package model

import "time"

// Ini adalah representasi tabel di PostgreSQL nanti
type SessionLog struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Events string `gorm:"type:jsonb" json:"events"`

	CreatedAt time.Time `json:"created_at"`
}
