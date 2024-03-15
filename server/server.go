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
// TODO: read best way to structure a go project (at least lib x main=cli=http server)
// TODO: on deck creation pass a shuffled arg and then shuffle accordingly
var channel = make(chan func(decks *map[uuid.UUID]deck.Deck)) // <- Make channel and assign to var

func main() {
	// deck := deck.NewDefaultDeck()
	// fmt.Printf("Hello, deck %v", deck)
	// fmt.Printf("Hello")
	decks := make(map[uuid.UUID]deck.Deck)

	go func() {
		for {
			handler, ok := <-channel

			if !ok {
				return
			}

			handler(&decks)

			// fmt.Println("hey you", request.URL)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/open/{deck_guid}", open)
	r.HandleFunc("/create", create)
	err := http.ListenAndServe(":8000", r)
	fmt.Println(err)
}

func create(w http.ResponseWriter, r *http.Request) {
	newDeck := deck.NewDefaultDeck()
	channel <- func(decks *map[uuid.UUID]deck.Deck) {
		(*decks)[newDeck.Guid] = newDeck
	}

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

type CreatedDeck struct {
	Guid               uuid.UUID `json:"deck_id"`
	IsShuffled         bool      `json:"shuffled"`
	RemainingCardCount int       `json:"remaining"`
}

func (d *CreatedDeck) toJson() (string, error) {
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// we could have used r.PathValue("deck_guid") on go 1.22 (latest)
func open(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deck_guid, ok := vars["deck_guid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Error parsing request to /open/{deck_guid}")
	}

	guid, _ := uuid.Parse(deck_guid)
	channel <- func(decks *map[uuid.UUID]deck.Deck) {
		deck := (*decks)[guid]
		temp, _ := deck.ToJson()
		io.WriteString(w, temp)
	}
}
