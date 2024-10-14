package mock

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Статический набор данных
var songsDB = map[string]map[string]SongDetail{

	// статичные данные для демонстрации работы мок-сервера
	"Muse": {
		"Supermassive Black Hole": {
			ReleaseDate: "2006-07-16",
			Text:        "Ooh baby, don't you know I suffer?",
			Link:        "https://youtube.com",
		},
		"Starlight": {
			ReleaseDate: "2006-07-16",
			Text:        "Far away, this ship is taking me far away...",
			Link:        "https://youtube.com",
		},
	},
	"Radiohead": {
		"Creep": {
			ReleaseDate: "1992-09-21",
			Text:        "I wish I was special, you're so very special.",
			Link:        "https://youtube.com/watch?v=6nW3D-7ccU4",
		},
	},
}

func StartMockServer() {
	// Инициализация маршрутизатора Gin
	testRouter := gin.Default()

	// Эмуляция внешнего API
	testRouter.POST("/info", func(c *gin.Context) {
		var input struct {
			Group string `json:"group"`
			Title string `json:"title"`
		}

		// Получаем данные из тела POST-запроса
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println("DEBUG: Неверные параметры запроса.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
			return
		}

		// проверка наличия параметров
		if input.Group == "" || input.Title == "" {
			log.Println("DEBUG: Отсутствуют параметры group или title.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
			return
		}

		// Поиск информации о песне
		if songDetail, found := songsDB[input.Group][input.Title]; found {
			log.Printf("INFO: Успешный запрос для группы: %s, песни: %s\n", input.Group, input.Title)
			c.JSON(http.StatusOK, songDetail)
		} else {
			log.Printf("DEBUG: Песня не найдена для группы: %s, песни: %s\n", input.Group, input.Title)
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		}
	})

	// Запуск сервера на порту 8081
	if err := testRouter.Run(":8081"); err != nil {
		log.Fatalf("ERROR: Ошибка при запуске тестового сервера: %v", err)
	}
}
