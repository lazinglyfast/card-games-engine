// A generic backend engine for card games
package deck

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"sort"
)

// we could make a deck so extensible that it could work with any number of
// cards, suits and ranks or include other concepts entirely
// (i.e. a healing card)
// but if that's not an immediate or foreseeable requirement there's no need
// to over-engineer. Complexity must be tamed and one of the most effective
// ways to do that is to not add more code

type Deck struct {
	Cards []Card
	Guid  uuid.UUID
}

func (d *Deck) RemainingCardCount() int {
	return len(d.Cards)
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c *Card) Code() string {
	// using the ASCII subset of UTF-8 so this is ok
	rank := c.Rank.String()[0]
	suit := c.Suit.String()[0]
	return fmt.Sprintf("%s%s", string(rank), string(suit))
}

func newCard(rank Rank, suit Suit) Card {
	return Card{rank, suit}
}

func ParseCard(code string) (Card, error) {
	n := len(code)
	rank, rankErr := parseRank(code[:(n - 1)])
	suit, suitErr := parseSuit(code[n-1:])
	if suitErr != nil || rankErr != nil {
		msg := fmt.Sprintf("Failed to parse %s into Card", code)
		return defaultCard(), errors.New(msg)
	}
	return newCard(rank, suit), nil

}

func defaultCard() Card {
	return newCard(0, 0)
}

type Suit int64

// suits are specified in increasing order according to domain rules
// (Spades > Diamonds > Clubs > Hearts)
const (
	Hearts Suit = iota
	Clubs
	Diamonds
	Spades
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "SPADES"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	case Hearts:
		return "HEARTS"
	}
	return "UNKNOWN SUIT"
}

func parseSuit(suit string) (Suit, error) {
	switch suit {
	case "S":
		return Spades, nil
	case "D":
		return Spades, nil
	case "C":
		return Spades, nil
	case "H":
		return Spades, nil
	default:
		{
			msg := fmt.Sprintf("Invalid suit code: %v", suit)
			return 0, errors.New(msg)
		}
	}
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
		return "ACE"
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
		return "JACK"
	case Queen:
		return "QUEEN"
	case King:
		return "KING"
	}
	return "UNKNOWN RANK"
}

func parseRank(rank string) (Rank, error) {
	switch rank {
	case "A":
		return Ace, nil
	case "2":
		return V2, nil
	case "3":
		return V3, nil
	case "4":
		return V4, nil
	case "5":
		return V5, nil
	case "6":
		return V6, nil
	case "7":
		return V7, nil
	case "8":
		return V8, nil
	case "9":
		return V9, nil
	case "10":
		return V10, nil
	case "J":
		return Jack, nil
	case "Q":
		return Queen, nil
	case "K":
		return King, nil
	default:
		{
			msg := fmt.Sprintf("Invalid rank code: %v", rank)
			return 0, errors.New(msg)
		}
	}
}

func NewDefaultDeck() Deck {
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
	return NewDeck(cards)
}

func NewDeck(cards []Card) Deck {
	guid := uuid.New()
	return Deck{
		cards,
		guid,
	}
}

func NewEmptyDeck() Deck {
	cards := []Card{}
	return NewDeck(cards)
}

func (d *Deck) Shuffle() {
	n := d.RemainingCardCount()
	indices := rand.Perm(n)
	cards := []Card{}
	for _, i := range indices {
		cards = append(cards, d.Cards[i])
	}
	d.Cards = cards
}

// this method (and the whole design really) would need to be rethinked if
// joker cards get introduced
func (d *Deck) IsShuffled() bool {
	n := d.RemainingCardCount()
	for i := 0; i < (n - 1); i += 1 {
		card := d.Cards[i]
		nextCard := d.Cards[i+1]

		suitOutOfOrder := card.Suit < nextCard.Suit
		if suitOutOfOrder {
			return true
		}

		if card.Suit == nextCard.Suit {
			rankOutOfOrder := card.Rank > nextCard.Rank
			if rankOutOfOrder {
				return true
			}
		}
	}
	return false
}

func (d *Deck) unshuffle() {
	sort.SliceStable(d.Cards, func(i, j int) bool {
		if d.Cards[i].Suit == d.Cards[j].Suit {
			// in a unshuffled deck lowest ranks come first
			// as defined by domain rules
			return d.Cards[i].Rank < d.Cards[j].Rank
		}
		// in a unshuffled deck highest suits come first
		// as defined by domain rules
		return d.Cards[i].Suit > d.Cards[j].Suit
	})
}

// do not resort to premature optimization (i.e a stack)
// will this ever be the bottle neck in our future games? Probably not
func (d *Deck) Draw(count int) []Card {
	// early return on unhappy path: less nesting improves code readability
	if count < 1 {
		return []Card{}
	}

	n := d.RemainingCardCount()
	if count > n {
		count = n
	}

	to := n - count
	remainingCards := d.Cards[:to]
	drawnCards := d.Cards[to:]
	d.Cards = remainingCards
	return drawnCards
}
