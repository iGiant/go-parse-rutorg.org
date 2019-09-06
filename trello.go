package main

import (
	"fmt"
	"github.com/VojtechVitek/go-trello"
	"os"
)

const listID = "5d71f9e38bebb648f1de6e30"

var (
	api, token string
)

func init() {
	api = os.Getenv("TrelloApiKey")
	token = os.Getenv("TrelloToken")
}

func getFilms() ([]movie, error) {
	if api == "" || token == "" {
		return []movie{}, fmt.Errorf("ошибка получения переменных среды")
	}
	client, err := trello.NewAuthClient(api, &token)
	if err != nil {
		return []movie{}, fmt.Errorf("ошибка идентификации")
	}
	list, err := client.List(listID)
	if err != nil {
		return []movie{}, fmt.Errorf("ошибка получения информации о колонке с фильмами")
	}
	cards, err := list.Cards()
	if err != nil {
		return []movie{}, fmt.Errorf("ошибка получения информации о карточках")
	}
	result := make([]movie, 0)
	for _, card := range cards {
		film := movie{show: "1"}
		film.name = card.Name
		result = append(result, film)
	}
	return result, nil
}
