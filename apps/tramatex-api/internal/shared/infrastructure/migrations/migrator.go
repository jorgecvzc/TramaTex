package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gorm.io/gorm"
)

// RunMigrations ejecuta todas las migraciones en orden
func RunMigrations(db *gorm.DB) error {
	fmt.Println("🔄 Running database migrations...")

	// Migration 1 & 2: Now handled by v2_001_init_iam.sql (consolidated migrations)
	// Uncomment below if you need GORM auto-migration for development/testing:
	// if err := createUsersTable(db); err != nil {
	// 	return fmt.Errorf("failed to create users table: %w", err)
	// }
	// if err := seedAdminUser(db); err != nil {
	// 	return fmt.Errorf("failed to seed admin user: %w", err)
	// }

	// Migration 3: Execute SQL migrations from /app/migrations directory
	if err := executeSQLMigrations(db); err != nil {
		return fmt.Errorf("failed to execute SQL migrations: %w", err)
	}

	fmt.Println("✅ All migrations completed successfully")
	return nil
}

// executeSQLMigrations executes all SQL migration files in order
func executeSQLMigrations(db *gorm.DB) error {
	migrationsPath, err := resolveMigrationsPath()
	if err != nil {
		fmt.Println("  ⚠ No migrations directory found, skipping SQL migrations")
		return nil
	}

	// Read all .sql files
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("  ℹ No SQL migration files found")
		return nil
	}

	// Sort files to ensure execution order
	sort.Strings(files)

	// Get raw SQL database connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	fmt.Printf("  📂 Found %d SQL migration files\n", len(files))

	// Execute each migration file
	for _, file := range files {
		fileName := filepath.Base(file)

		// Skip if already executed (check by creating a migrations table)
		if err := ensureMigrationsTable(sqlDB); err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}

		// Check if migration was already executed
		executed, err := wasMigrationExecuted(sqlDB, fileName)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if executed {
			fmt.Printf("  ✓ %s (already executed)\n", fileName)
			continue
		}

		// Read and execute migration
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		// Execute the SQL and record migration in a single transaction
		// so that if the SQL fails, the migration is not recorded.
		tx, err := sqlDB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", fileName); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", fileName, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", fileName, err)
		}

		fmt.Printf("  ✓ Executed %s\n", fileName)
	}

	return nil
}

func resolveMigrationsPath() (string, error) {
	candidates := []string{
		"/app/migrations",
		"./migrations",
		"./apps/tramatex-api/migrations",
		"../migrations",
		"../../migrations",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			continue
		}

		files, err := filepath.Glob(filepath.Join(candidate, "*.sql"))
		if err != nil {
			continue
		}
		if len(files) > 0 {
			fmt.Printf("  📁 Using migrations directory: %s\n", candidate)
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no migration directory with SQL files found")
}

// ensureMigrationsTable creates the schema_migrations table if it doesn't exist
func ensureMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) UNIQUE NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

// wasMigrationExecuted checks if a migration was already executed
func wasMigrationExecuted(db *sql.DB, fileName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", fileName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// recordMigration records a migration as executed
func recordMigration(db *sql.DB, fileName string) error {
	_, err := db.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", fileName)
	return err
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

// seedAdminUser creates the default admin user if it doesn't exist
// Idempotent: doesn't fail if user already exists
func seedAdminUser(db *gorm.DB) error {
	adminEmail := "admin@tramatex.local"

	// Check if admin user already exists
	var count int64
	db.Model(&User{}).Where("email = ?", adminEmail).Count(&count)

	if count > 0 {
		fmt.Println("  ✓ Admin user already exists")
		return nil
	}

	// Create admin user with bcrypt hash for password: admin123
	// Hash generated with: bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	adminUser := User{
		ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Email:    adminEmail,
		Password: "$2a$10$Gd8.JP/L3j.vNvap81EpjuX7G4u5KKLmf10TSxmu779Mq/HdC/B9e",
		Role:     "admin",
		IsActive: true,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	fmt.Println("  ✓ Created admin user (admin@tramatex.local)")
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
