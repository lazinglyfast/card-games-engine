package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

// I'd rather not expose the internals of Deck and Card but json parsing needs it
// as does go-cmp
// if a solid API were indeed the goal I'd then create a json serializable type and convert into it
type Deck struct {
	Cards []Card
	Guid  uuid.UUID
}

func deckFromJson(text string) (Deck, error) {
	deck := Deck{}
	bytes := []byte(text)
	err := json.Unmarshal(bytes, &deck)
	if err != nil {
		return deck, err
	}
	return deck, nil

}

func (d *Deck) toJson() (string, error) {
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

func newCard(rank Rank, suit Suit) Card {
	return Card{rank, suit}
}

type Suit int64

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "S"
	case Diamonds:
		return "D"
	case Clubs:
		return "C"
	case Hearts:
		return "H"
	}
	return "unknown suit"
}

type Rank int64

const (
	Ace Rank = iota
	V2
	V3
	V4
	V5
	V6
	V7
	V8
	V9
	V10
	Jack
	Queen
	King
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case V2:
		return "2"
	case V3:
		return "3"
	case V4:
		return "4"
	case V5:
		return "5"
	case V6:
		return "6"
	case V7:
		return "7"
	case V8:
		return "8"
	case V9:
		return "9"
	case V10:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	}
	return "unknown rank"
}

func newDefaultDeck() Deck {
	cards := []Card{}
	suits := []Suit{Spades, Diamonds, Clubs, Hearts}
	ranks := []Rank{Ace, V2, V3, V4, V5, V6, V7, V8, V9, V10, Jack, Queen, King}
	for _, suit := range suits {
		for _, rank := range ranks {
			card := Card{
				rank,
				suit,
			}
			cards = append(cards, card)
		}
	}
	return newDeck(cards)
}

func newDeck(cards []Card) Deck {
	guid := uuid.New()
	return Deck{
		cards,
		guid,
	}
}

func (d *Deck) shuffle() {
	n := len(d.Cards)
	indices := rand.Perm(n)
	cards := []Card{}
	for _, i := range indices {
		cards = append(cards, d.Cards[i])
	}
	d.Cards = cards
}

// do not resort to premature optimization (i.e a stack)
// will it ever be the bottle neck in our future games? Probably not
func (d *Deck) draw(count int) []Card {
	// early return on unhappy path
	if count < 1 {
		return []Card{}
	}

	n := len(d.Cards)
	to := n - count
	remainingCards := d.Cards[:to]
	drawnCards := d.Cards[to:]
	d.Cards = remainingCards
	return drawnCards
}
