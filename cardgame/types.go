package cardgame

// DeckResponse представляет ответ при создании колоды
type DeckResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

// DrawResponse представляет ответ при вытягивании карт
type DrawResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Remaining int    `json:"remaining"`
}

// Card представляет одну карту
type Card struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}

