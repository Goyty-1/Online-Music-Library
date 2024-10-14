package handlers

import (
	"Online-Music-Library/config"
	"Online-Music-Library/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type GetSondsRequest struct {
	Group       string `json:"group"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
}

// GetSongs godoc
// @Summary Получить список песен
// @Description Возвращает список всех песен, сохраненных в базе данных
// @Tags songs
// @Produce  json
// @Success 200 {array} models.Song
// @Failure 500 {string} string "Ошибка сервера"
// @Router /songs [get]
func GetSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на получение всех песен")

	var songs []models.Song
	query := config.DB

	//фильтрация по параметрам
	group := r.URL.Query().Get("group")
	title := r.URL.Query().Get("title")
	releaseDate := r.URL.Query().Get("release_date")

	if group != "" {
		query = query.Where("group = ?", group)
	}
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if releaseDate != "" {
		query = query.Where("release_date = ?", releaseDate)
	}

	// Пагинация
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Получение данных
	if err := query.Limit(pageSize).Offset(offset).Find(&songs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Песни успешно получены")

	// Установка заголовков для пагинации
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongByID godoc
// @Summary Получить песню по ID
// @Description Возвращает песню по её идентификатору
// @Tags songs
// @Produce  json
// @Param id path int true "ID песни"
// @Success 200 {object} models.Song
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /songs/{id} [get]
func GetSongByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Получен запрос на получение песни с ID=%s\n", id)

	var song models.Song
	if err := config.DB.First(&song, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Песня не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Песня с ID=%s успешно найдена\n", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// CreateSong godoc
// @Summary Создать новую песню
// @Description Создает новую песню с обогащением данных через внешний API
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body struct { Group string "example:Muse"; Title string "example:Supermassive Black Hole" } true "Параметры песни"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /songs [post]
func CreateSong(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на создание новой песни")

	var input struct {
		Group string `json:"group"`
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Вызов внешнего API для получения дополнительных данных
	songDetail, err := FetchSongDetails(input.Group, input.Title)
	if err != nil {
		log.Println("Ошибка при запросе данных от внешнего API:", err)
		http.Error(w, "Ошибка получения данных из внешнего API", http.StatusInternalServerError)
		return
	}

	// Создание новой песни с обогащенными данными
	song := models.Song{
		Group:       input.Group,
		Title:       input.Title,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	if err := config.DB.Create(&song).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// UpdateSong godoc
// @Summary Обновить песню по ID
// @Description Обновляет информацию о песне по её идентификатору
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "ID песни"
// @Param song body struct { Group string "example:Muse"; Title string "example:Supermassive Black Hole"; ReleaseDate string "example:2006-07-16"; Text string "example:Lyrics..."; Link string "example:https://youtube.com" } true "Обновленные данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /songs/{id} [put]
func UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := config.DB.First(&song, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Песня не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var input struct {
		Group       *string `json:"group"`
		Title       *string `json:"title"`
		ReleaseDate *string `json:"release_date"`
		Text        *string `json:"text"`
		Link        *string `json:"link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновление полей, если они присутствуют в запросе
	if input.Group != nil {
		song.Group = *input.Group
	}
	if input.Title != nil {
		song.Title = *input.Title
	}
	if input.ReleaseDate != nil {
		song.ReleaseDate = *input.ReleaseDate
	}
	if input.Text != nil {
		song.Text = *input.Text
	}
	if input.Link != nil {
		song.Link = *input.Link
	}

	if err := config.DB.Save(&song).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// DeleteSong godoc
// @Summary Удалить песню по ID
// @Description Удаляет песню из базы данных по её идентификатору
// @Tags songs
// @Param id path int true "ID песни"
// @Success 204 {string} string "Песня удалена"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /songs/{id} [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := config.DB.Delete(&models.Song{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Вспомогательная функция для вызова внешнего API
func FetchSongDetails(group, title string) (*models.Song, error) {

	// Формируется URL для запроса на мок-сервер (без параметров в URL)
	apiURL := fmt.Sprintf("%s/info", os.Getenv("API_BASE_URL"))

	// Создается JSON-данные для POST-запроса
	requestData := map[string]string{
		"group": group,
		"title": title,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Ошибка при сериализации данных для POST-запроса:", err)
		return nil, err
	}

	// Выполняется POST запрос на сервер с JSON-данными
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Ошибка при вызове внешнего API через POST:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверка на успешный статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Внешнее API вернуло статус %d: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("внешнее API вернуло статус %d", resp.StatusCode)
	}

	// Структура для получения данных о песне
	var songDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	// Декодируется JSON-ответ от сервера
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Println("Ошибка декодирования ответа внешнего API:", err)
		return nil, err
	}

	// Возвращается информация о песне в формате модели Song
	return &models.Song{
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}, nil
}
