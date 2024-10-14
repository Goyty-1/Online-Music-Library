package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model

	Group       string `json:"group" gorm:"not null"`
	Title       string `json:"title" gorm:"not null"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
