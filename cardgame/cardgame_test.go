package cardgame

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// createMockServer создаёт тестовый HTTP сервер с заранее определённой колодой
// cards - список карт в порядке их вытягивания (первая карта будет вытянута первой)
func createMockServer(cards []Card) *httptest.Server {
	deckID := "test_deck_123"
	currentCardIndex := 0

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Обработка создания колоды
		if strings.Contains(r.URL.Path, "/new/shuffle/") {
			resp := DeckResponse{
				Success:   true,
				DeckID:    deckID,
				Shuffled:  true,
				Remaining: len(cards),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Обработка вытягивания карты
		if strings.Contains(r.URL.Path, "/draw/") {
			if currentCardIndex >= len(cards) {
				http.Error(w, "no cards left", http.StatusBadRequest)
				return
			}

			// Считываем параметр count
			countStr := r.URL.Query().Get("count")
			count := 1
			if countStr == "2" {
				count = 2
			}

			drawnCards := []Card{}
			for i := 0; i < count && currentCardIndex < len(cards); i++ {
				drawnCards = append(drawnCards, cards[currentCardIndex])
				currentCardIndex++
			}

			resp := DrawResponse{
				Success:   true,
				DeckID:    deckID,
				Cards:     drawnCards,
				Remaining: len(cards) - currentCardIndex,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
}

// TestPlayGameUserGuessesCorrectly проверяет случай, когда пользователь угадал
func TestPlayGameUserGuessesCorrectly(t *testing.T) {
	// Колода: 3 карты до дамы, затем дама
	cards := []Card{
		{Code: "7D", Value: "7", Suit: "DIAMONDS"},
		{Code: "3C", Value: "3", Suit: "CLUBS"},
		{Code: "KS", Value: "KING", Suit: "SPADES"},
		{Code: "QH", Value: "QUEEN", Suit: "HEARTS"},
	}

	server := createMockServer(cards)
	defer server.Close()

	var output bytes.Buffer
	client := &Client{
		baseURL: server.URL,
		client:  server.Client(),
		output:  &output,
	}

	// Пользователь угадал: нужно снять 4 карты
	result, err := client.PlayGame(4)

	if err != nil {
		t.Errorf("PlayGame() returned error: %v", err)
	}

	if !result {
		t.Error("Expected user to win, but they lost")
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Вы угадали!") {
		t.Errorf("Expected victory message, got: %s", outputStr)
	}
}

// TestPlayGameUserGuessesTooLow проверяет случай, когда пользователь угадал меньше
func TestPlayGameUserGuessesTooLow(t *testing.T) {
	// Колода: 2 карты до дамы, затем дама
	cards := []Card{
		{Code: "4C", Value: "4", Suit: "CLUBS"},
		{Code: "9D", Value: "9", Suit: "DIAMONDS"},
		{Code: "QC", Value: "QUEEN", Suit: "CLUBS"},
	}

	server := createMockServer(cards)
	defer server.Close()

	var output bytes.Buffer
	client := &Client{
		baseURL: server.URL,
		client:  server.Client(),
		output:  &output,
	}

	// Пользователь ошибся: угадал 10, а на самом деле 3
	result, err := client.PlayGame(10)

	if err != nil {
		t.Errorf("PlayGame() returned error: %v", err)
	}

	if result {
		t.Error("Expected user to lose, but they won")
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Вы проиграли!") {
		t.Errorf("Expected defeat message, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "Правильный ответ: 3") {
		t.Errorf("Expected correct answer 3, got: %s", outputStr)
	}
}

// TestPlayGameQueenIsFirstCard проверяет случай, когда дама - первая карта
func TestPlayGameQueenIsFirstCard(t *testing.T) {
	cards := []Card{
		{Code: "QS", Value: "QUEEN", Suit: "SPADES"},
	}

	server := createMockServer(cards)
	defer server.Close()

	var output bytes.Buffer
	client := &Client{
		baseURL: server.URL,
		client:  server.Client(),
		output:  &output,
	}

	// Пользователь угадал 1
	result, err := client.PlayGame(1)

	if err != nil {
		t.Errorf("PlayGame() returned error: %v", err)
	}

	if !result {
		t.Error("Expected user to win")
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "QUEEN of SPADES") {
		t.Errorf("Expected to see QUEEN of SPADES in output, got: %s", outputStr)
	}
}

// TestPlayGameMultipleCards проверяет вывод нескольких карт
func TestPlayGameMultipleCards(t *testing.T) {
	cards := []Card{
		{Code: "AC", Value: "ACE", Suit: "CLUBS"},
		{Code: "2H", Value: "2", Suit: "HEARTS"},
		{Code: "JD", Value: "JACK", Suit: "DIAMONDS"},
		{Code: "QD", Value: "QUEEN", Suit: "DIAMONDS"},
	}

	server := createMockServer(cards)
	defer server.Close()

	var output bytes.Buffer
	client := &Client{
		baseURL: server.URL,
		client:  server.Client(),
		output:  &output,
	}

	_, err := client.PlayGame(4)
	if err != nil {
		t.Errorf("PlayGame() returned error: %v", err)
	}

	outputStr := output.String()
	expectedCards := []string{
		"ACE of CLUBS",
		"2 of HEARTS",
		"JACK of DIAMONDS",
		"QUEEN of DIAMONDS",
	}

	for _, expectedCard := range expectedCards {
		if !strings.Contains(outputStr, expectedCard) {
			t.Errorf("Expected to see '%s' in output, got: %s", expectedCard, outputStr)
		}
	}
}
