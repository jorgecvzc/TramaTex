package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	Email    string `gorm:"unique;type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
	Role     string `gorm:"type:varchar(50)"`
	IsActive bool   `gorm:"default:true"`
}

func (User) TableName() string {
	return "users"
}

func main() {
	// Connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=tramatex password=tramatex dbname=tramatex port=5432 sslmode=disable"
	}

	fmt.Println("🔍 Conectando a la base de datos...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error conectando: %v", err)
	}

	// 1. Verificar si la tabla existe
	fmt.Println("\n📋 Verificando tabla 'users'...")
	if !db.Migrator().HasTable(&User{}) {
		fmt.Println("❌ Tabla 'users' NO EXISTE. Creándola...")
		if err := db.Migrator().CreateTable(&User{}); err != nil {
			log.Fatalf("❌ Error creando tabla: %v", err)
		}
		fmt.Println("✅ Tabla 'users' creada exitosamente")
	} else {
		fmt.Println("✅ Tabla 'users' existe")
	}

	// 2. Contar usuarios
	var count int64
	db.Model(&User{}).Count(&count)
	fmt.Printf("📊 Total de usuarios: %d\n", count)

	// 3. Buscar usuario admin
	var admin User
	result := db.Where("email = ?", "admin@tramatex.local").First(&admin)
	
	if result.Error != nil {
		fmt.Println("\n❌ Usuario admin NO encontrado. Creándolo...")
		
		// Generar hash para "admin123"
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ Error generando hash: %v", err)
		}

		admin = User{
			ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			Email:    "admin@tramatex.local",
			Password: string(hash),
			Role:     "admin",
			IsActive: true,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("❌ Error creando usuario: %v", err)
		}

		fmt.Println("✅ Usuario admin creado exitosamente")
		fmt.Printf("   Hash generado: %s\n", string(hash))
	} else {
		fmt.Println("\n✅ Usuario admin encontrado:")
		fmt.Printf("   ID: %s\n", admin.ID)
		fmt.Printf("   Email: %s\n", admin.Email)
		fmt.Printf("   Role: %s\n", admin.Role)
		fmt.Printf("   Active: %v\n", admin.IsActive)
		fmt.Printf("   Password Hash: %s\n", admin.Password)
	}

	// 4. Verificar el hash de la contraseña
	fmt.Println("\n🔐 Verificando contraseña 'admin123'...")
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte("admin123"))
	if err != nil {
		fmt.Printf("❌ La contraseña 'admin123' NO coincide con el hash\n")
		fmt.Printf("   Error: %v\n", err)
		
		// Regenerar hash
		fmt.Println("\n🔄 Regenerando hash...")
		newHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ Error generando nuevo hash: %v", err)
		}

		db.Model(&admin).Update("password", string(newHash))
		fmt.Println("✅ Hash actualizado exitosamente")
		fmt.Printf("   Nuevo hash: %s\n", string(newHash))
	} else {
		fmt.Println("✅ La contraseña 'admin123' es correcta")
	}

	fmt.Println("\n✅ Diagnóstico completado. Intenta hacer login con:")
	fmt.Println("   Email: admin@tramatex.local")
	fmt.Println("   Password: admin123")
}
