package main

import (
	"encoding/json"
	"example.com/deck"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	// "net/url"
	// "strings"
)

// TODO: read best-practices
// TODO: read best way to structure a go project (at least lib x main=cli=http
// server)
// TODO: on deck creation pass a shuffled arg and then shuffle accordingly
type HandlerContext struct {
	decks *map[uuid.UUID]deck.Deck
}

func NewHandlerContext(decks *map[uuid.UUID]deck.Deck) *HandlerContext {
	if decks == nil {
		panic("nil MongoDB session!")
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
	fmt.Println(err)
}

func (ctx *HandlerContext) Create(w http.ResponseWriter, r *http.Request) {
	newDeck := deck.NewDefaultDeck()
	(*ctx.decks)[newDeck.Guid] = newDeck
	createdDeck := intoCreatedDeck(newDeck)
	response, err := createdDeck.toJson()
	if err != nil {
		msg := fmt.Sprintf("%v\n", err)
		io.WriteString(w, msg)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, response)
	}
}

type CreatedDeck struct {
	Guid               uuid.UUID `json:"deck_id"`
	IsShuffled         bool      `json:"shuffled"`
	RemainingCardCount int       `json:"remaining"`
}

func intoCreatedDeck(deck deck.Deck) CreatedDeck {
	Guid := deck.Guid
	IsShuffled := deck.IsShuffled()
	RemainingCardCount := deck.RemainingCardCount()

	return CreatedDeck{
		Guid,
		IsShuffled,
		RemainingCardCount,
	}
}

func (d *CreatedDeck) toJson() (string, error) {
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
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

	deck := (*ctx.decks)[guid]
	openDeck := intoOpenDeck(deck)
	json, jsonErr := openDeck.toJson()
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("%v", jsonErr))
		return
	}

	io.WriteString(w, json)
}

func intoOpenDeck(deck deck.Deck) OpenDeck {
	Guid := deck.Guid
	IsShuffled := deck.IsShuffled()
	RemainingCardCount := deck.RemainingCardCount()
	Cards := []OpenCard{}
	for _, card := range deck.Cards {
		Cards = append(Cards, intoOpenCard(card))
	}

	return OpenDeck{
		Guid,
		IsShuffled,
		RemainingCardCount,
		Cards,
	}
}

type OpenDeck struct {
	Guid               uuid.UUID  `json:"deck_id"`
	IsShuffled         bool       `json:"shuffled"`
	RemainingCardCount int        `json:"remaining"`
	Cards              []OpenCard `json:"cards"`
}

type OpenCard struct {
	Rank string `json:"value"`
	Suit string `json:"suit"`
	Code string `json:"code"`
}

func intoOpenCard(card deck.Card) OpenCard {
	Rank := card.Rank.String()
	Suit := card.Suit.String()
	Code := card.Code()
	return OpenCard{
		Rank,
		Suit,
		Code,
	}
}

// some generics would come in handy here to avoid repeating the JSON
// marshalling logic yet again
func (d *OpenDeck) toJson() (string, error) {
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
