package config

import (
	"Online-Music-Library/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Загружаются переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	// Формируется строка подключения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	// Миграция схемы базы данных
	if err := db.AutoMigrate(&models.Song{}); err != nil {
		log.Fatal("Ошибка миграции базы данных: ", err)
	}

	DB = db
	log.Println("База данных успешно подключена и мигрирована!")
	log.Println("Swagger для тестирования - ", "http://localhost:8080/swagger/index.html")
}
