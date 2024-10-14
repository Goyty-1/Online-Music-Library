package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

// Структура для хранения данных о песне
type SongDetail struct {
	Group       string `json:"group"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Загрузка данных о песне из JSON-файла
func GetSongDetailFromJSON(group, song string) (*SongDetail, error) {
	// Чтение данных из локального JSON-файла
	data, err := ioutil.ReadFile("data/song_mock.json")
	if err != nil {
		log.Printf("ERROR: Не удалось прочитать файл JSON: %v\n", err)
		return nil, err
	}

	var songs []SongDetail
	if err := json.Unmarshal(data, &songs); err != nil {
		log.Printf("ERROR: Ошибка при парсинге JSON: %v\n", err)
		return nil, err
	}

	// Поиск песни по параметрам group и song
	for _, s := range songs {
		if s.Group == group && s.Title == song {
			return &s, nil
		}
	}

	// Если песня не найдена
	return nil, errors.New("song not found")
}
