package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// RunMigrations ejecuta todas las migraciones en orden
func RunMigrations(db *gorm.DB) error {
	fmt.Println("🔄 Running database migrations...")

	// Migration 1: Create users table (idempotent - doesn't fail if table exists)
	if err := createUsersTable(db); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	fmt.Println("✅ All migrations completed successfully")
	return nil
}

// createUsersTable crea la tabla de usuarios si no existe
// Idempotent: no falla si la tabla ya existe
func createUsersTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&User{}) {
		if err := db.Migrator().CreateTable(&User{}); err != nil {
			return err
		}
		fmt.Println("  ✓ Created users table")
	} else {
		fmt.Println("  ✓ Users table already exists")
	}
	return nil
}

// User domain model para migraciones
type User struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string     `gorm:"unique;type:varchar(255);not null"`
	Password  string     `gorm:"type:varchar(255);not null"`
	Role      string     `gorm:"type:varchar(50);default:'commercial'"`
	IsActive  bool       `gorm:"default:true"`
	CreatedAt time.Time  `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime:milli"`
	DeletedAt *time.Time `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
