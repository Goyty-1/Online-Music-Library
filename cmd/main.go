package main

import (
	"Online-Music-Library/config"
	_ "Online-Music-Library/docs"
	"Online-Music-Library/handlers"
	"Online-Music-Library/mock"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

// @title Online Music Library API
// @version 1.0
// @description Это API для управления онлайн библиотекой песен
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	log.Println("Запуск приложения...")

	// старт мок-сервера
	go mock.StartMockServer()

	// подключение к базе данных
	config.ConnectDatabase()

	// Создание нового роутера
	router := mux.NewRouter()

	// Регистрация маршрутов
	router.HandleFunc("/songs", handlers.GetSongs).Methods("GET")
	router.HandleFunc("/songs/{id}", handlers.GetSongByID).Methods("GET")
	router.HandleFunc("/songs", handlers.CreateSong).Methods("POST")
	router.HandleFunc("/songs/{id}", handlers.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", handlers.DeleteSong).Methods("DELETE")

	// Маршрут для Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Запуск HTTP-сервера
	port := ":" + getEnv("APP_PORT", "8080")
	log.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

// Функция для получения переменной окружения с значением по умолчанию
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
