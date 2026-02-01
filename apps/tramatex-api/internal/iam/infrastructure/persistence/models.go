package persistence

import (
	"time"
)

// UserModel representa un usuario en la base de datos PostgreSQL
// Mapea la tabla 'users' a una estructura Go para GORM
type UserModel struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string     `gorm:"uniqueIndex;not null;type:varchar(255)"`
	PasswordHash string     `gorm:"not null;type:varchar(255)"`
	Role         string     `gorm:"not null;default:'commercial';type:varchar(50)"`
	IsActive     bool       `gorm:"not null;default:true"`
	CreatedAt    time.Time  `gorm:"autoCreateTime:milli;type:timestamp with time zone"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime:milli;type:timestamp with time zone"`
	DeletedAt    *time.Time `gorm:"type:timestamp with time zone"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (UserModel) TableName() string {
	return "users"
}

// SeedAdminUser retorna un UserModel para el usuario admin inicial
// Password: Admin@12345 (hasheado con bcrypt cost=10)
func SeedAdminUser() UserModel {
	return UserModel{
		ID:           "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Email:        "admin@tramatex.local",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36P4/tvW2",
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
