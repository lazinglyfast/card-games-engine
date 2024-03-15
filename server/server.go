package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	r.HandleFunc("/draw/{deck_guid}", ctx.Draw)

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}

// GET /create?cards=A2,8C&shuffled where cards and shuffled are optional
func (ctx *HandlerContext) Create(w http.ResponseWriter, r *http.Request) {
	deck, err := deriveDeck(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("%v\n", err)
		io.WriteString(w, msg)
		return
	}

	shuffled := r.URL.Query().Get("shuffled")
	shouldShuffle := strings.ToLower(shuffled) == "true"
	if shouldShuffle {
		deck.Shuffle()
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
		deck := deck.NewDefaultDeck()
		return deck, nil
	}

	cards, err := parseCards(codes)
	if err != nil {
		return deck.NewEmptyDeck(), err
	}
	return deck.NewDeck(cards), nil
}

func parseCards(query string) ([]deck.Card, error) {
	cards := []deck.Card{}
	codes := strings.Split(query, ",")
	for _, code := range codes {
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
// GET /open/{guid}
func (ctx *HandlerContext) Open(w http.ResponseWriter, r *http.Request) {
	deck, err := retrieveDeck(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%v", err))
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

// GET /draw/{guid}?count=2 where count is optional
func (ctx *HandlerContext) Draw(w http.ResponseWriter, r *http.Request) {
	deck, err := retrieveDeck(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%v", err))
		return
	}

	param := r.URL.Query().Get("count")
	count, err := strconv.Atoi(param)
	if err != nil {
		count = 1
	}

	cards := deck.Draw(count)
	(*ctx.decks)[deck.Guid] = deck
	body, err := marshallOpenCardsToJson(IntoOpenCards(cards))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("%v", err))
		return
	}

	io.WriteString(w, body)

}

func retrieveDeck(ctx *HandlerContext, r *http.Request) (deck.Deck, error) {
	vars := mux.Vars(r)
	deck_guid, ok := vars["deck_guid"]
	if !ok {
		msg := "Failed to extract {deck_guid} from /open/{deck_guid}"
		return deck.NewEmptyDeck(), errors.New(msg)
	}

	guid, uuidErr := uuid.Parse(deck_guid)
	if uuidErr != nil {
		return deck.NewEmptyDeck(), uuidErr
	}

	foundDeck, ok := (*ctx.decks)[guid]
	if !ok {
		msg := fmt.Sprintf("There's no deck with identifier %v", guid)
		return deck.NewEmptyDeck(), errors.New(msg)
	}

	return foundDeck, nil
}
