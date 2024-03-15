package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"example.com/deck"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// TODO: read best-practices
// TODO: read best way to structure a go project (at least lib x main=cli=http
// server)
// TODO: on deck creation pass a shuffled arg and then shuffle accordingly
// FEATURE: authentication so that one person's deck cannot be drawn by someone
// else
type HandlerContext struct {
	decks *map[uuid.UUID]deck.Deck
}

func NewHandlerContext(decks *map[uuid.UUID]deck.Deck) *HandlerContext {
	if decks == nil {
		panic("nil decks!")
	}
	return &HandlerContext{decks}
}

func main() {
	decks := make(map[uuid.UUID]deck.Deck)
	ctx := NewHandlerContext(&decks)
	r := mux.NewRouter()
	r.HandleFunc("/create", ctx.Create)
	r.HandleFunc("/open/{deck_guid}", ctx.Open)

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}

func (ctx *HandlerContext) Create(w http.ResponseWriter, r *http.Request) {
	deck, err := deriveDeck(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("%v\n", err)
		io.WriteString(w, msg)
		return
	}

	(*ctx.decks)[deck.Guid] = deck
	createdDeck := intoCreatedDeck(deck)
	response, err := createdDeck.toJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("%v\n", err)
		io.WriteString(w, msg)
		return
	}

	io.WriteString(w, response)
}

func deriveDeck(r *http.Request) (deck.Deck, error) {
	codes := r.URL.Query().Get("cards")
	if codes == "" {
		return deck.NewDefaultDeck(), nil
	}

	cards, err := parseCards(codes)
	if err != nil {
		return deck.NewEmptyDeck(), err
	}
	return deck.NewDeck(cards), nil
}

func parseCards(codes string) ([]deck.Card, error) {
	cards := []deck.Card{}
	splittedCodes := strings.Split(codes, ",")
	for _, code := range splittedCodes {
		card, err := deck.ParseCard(code)
		fmt.Println(code, card, err)
		if err != nil {
			return []deck.Card{}, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// we could have used r.PathValue("deck_guid") on go 1.22
func (ctx *HandlerContext) Open(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deck_guid, ok := vars["deck_guid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Failed to extract {deck_guid} from /open/{deck_guid}"
		io.WriteString(w, msg)
		return
	}

	guid, uuidErr := uuid.Parse(deck_guid)
	if uuidErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%v", uuidErr))
		return
	}

	deck, ok := (*ctx.decks)[guid]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("There's no deck with identifier %v", guid)
		io.WriteString(w, msg)
		return
	}

	openDeck := intoOpenDeck(deck)
	json, jsonErr := openDeck.toJson()
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("%v", jsonErr))
		return
	}

	io.WriteString(w, json)
}
