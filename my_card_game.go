package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DeckResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
	Shuffled  bool   `json:"shuffled"`
}

type CardImages struct {
	SVG string `json:"svg"`
	PNG string `json:"png"`
}

type Card struct {
	Code   string     `json:"code"`
	Image  string     `json:"image"`
	Images CardImages `json:"images"`
	Value  string     `json:"value"`
	Suit   string     `json:"suit"`
}

type DrawResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Remaining int    `json:"remaining"`
}

func main() {
	var number_of_queen int
	fmt.Print("Введите число: ")
	_, _ = fmt.Scan(&number_of_queen)
	if number_of_queen > 52 {
		fmt.Println("Плохое число - не может быть больше 52")
		return
	}

	resp1, err := http.Get("https://deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1")
	if err != nil {
		panic("error in ger request")
	}
	defer resp1.Body.Close()

	body, err := io.ReadAll(resp1.Body)
	if err != nil {
		panic("error in read body")
	}

	var deckResp DeckResponse
	err = json.Unmarshal(body, &deckResp)
	if err != nil {
		panic("error in convert body")
	}

	deck_id := deckResp.DeckID

	url := fmt.Sprintf("https://deckofcardsapi.com/api/deck/%s/draw/?count=%d", deck_id, number_of_queen)

	// fmt.Println(url)
	resp2, err := http.Get(url)
	if err != nil {
		panic("error in ger request")
	}
	defer resp2.Body.Close()

	body, err = io.ReadAll(resp2.Body)
	if err != nil {
		panic("error in read body: " + err.Error())
	}
	// fmt.Println("Response:", string(body))

	var drawResp DrawResponse
	err = json.Unmarshal(body, &drawResp)
	if err != nil {
		panic("error in decode: " + err.Error())
	}

	if drawResp.Cards[number_of_queen-1].Value == "QUEEN" {
		fmt.Println("Ты удачливый!")
	} else {
		fmt.Println("Попробуй ещё раз :(")
	}

}
