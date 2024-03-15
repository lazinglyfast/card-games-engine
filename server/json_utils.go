package main

import (
	"encoding/json"
	"example.com/deck"
	"github.com/google/uuid"
)

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

func intoOpenDeck(deck deck.Deck) OpenDeck {
	Guid := deck.Guid
	IsShuffled := deck.IsShuffled()
	RemainingCardCount := deck.RemainingCardCount()
	Cards := IntoOpenCards(deck.Cards)

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

func IntoOpenCards(cards []deck.Card) []OpenCard {
	openCards := []OpenCard{}
	for _, card := range cards {
		openCards = append(openCards, intoOpenCard(card))
	}
	return openCards
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

func marshallOpenCardsToJson(cards []OpenCard) (string, error) {
	jsonBytes, err := json.Marshal(cards)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
