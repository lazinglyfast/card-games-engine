// go test ./... # runs all tests

package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestJsonRoundtrip(t *testing.T) {
	cards := []Card{
		newCard(Ace, Spades),
		newCard(King, Hearts),
		newCard(Ace, Clubs),
	}
	deck := newDeck(cards)

	json, err := deck.toJson()
	if err != nil {
		t.Errorf("JSON-serialization of deck failed")
	}

	reconstructedDeck, err := deckFromJson(json)
	if err != nil {
		t.Errorf("JSON-deserialization of deck failed")
	}

	if !cmp.Equal(deck, reconstructedDeck) {
		t.Errorf("JSON roundtrip resulted in a different deck. Before: %v, After: %v", deck, reconstructedDeck)
	}
}

func TestDraw(t *testing.T) {
	deck := newDefaultDeck()
	actualCards := deck.draw(1)
	if len(actualCards) != 1 {
		t.Errorf("Expected to have drawn 1 card, instead drew %d", len(actualCards))
	}

	actualCard := actualCards[0]
	expectedCard := newCard(King, Hearts)
	if actualCard != expectedCard {
		t.Errorf("Expected drawn card of a unshuffled deck to be %v. Found %v instead", expectedCard, actualCard)
	}

	// do not introduce non-simple logic otherwise we'll end up having to write tests for tests
	// keep it simple
	if len(deck.Cards) != 51 {
		t.Errorf("Expected deck to have 51 cards after having drawn 1 card")
	}
}

func TestShuffle(t *testing.T) {
	deck := newDefaultDeck()
	deck.shuffle()
	unshuffledDeck := newDefaultDeck()
	// there's a 1/52! chance of this actually being true
	// instead of over-engineering let's accept this, at least for initial prototypes
	if cmp.Equal(deck, unshuffledDeck) {
		t.Errorf("Expected deck to be shuffled")
	}
}

func TestNewDeck(t *testing.T) {
	cards := []Card{
		newCard(Ace, Spades),
		newCard(King, Hearts),
		newCard(Ace, Clubs),
	}
	deck := newDeck(cards)
	n := len(deck.Cards)
	if n != 3 {
		t.Errorf("Expected default deck to have 3 cards but it has %d instead", n)
	}

	firstCard := deck.Cards[0]
	expectedCard := newCard(Ace, Spades)
	if firstCard != expectedCard {
		t.Errorf("Expected card to be %v but found %v instead", expectedCard, firstCard)
	}

	secondCard := deck.Cards[1]
	expectedCard = newCard(King, Hearts)
	if secondCard != expectedCard {
		t.Errorf("Expected card to be %v but found %v instead", expectedCard, secondCard)
	}

	thirdCard := deck.Cards[2]
	expectedCard = newCard(Ace, Clubs)
	if thirdCard != expectedCard {
		t.Errorf("Expected card to be %v but found %v instead", expectedCard, thirdCard)
	}
}

func TestDefaultDeck(t *testing.T) {
	deck := newDefaultDeck()
	n := len(deck.Cards)
	if n != 52 {
		t.Errorf("Expected default deck to have 52 cards but it has %d instead", n)
	}

	firstCard := deck.Cards[0]
	expectedCard := newCard(Ace, Spades)
	if firstCard != expectedCard {
		t.Errorf("Expected card to be %v but found %v instead", expectedCard, firstCard)
	}

	lastCard := deck.Cards[51]
	expectedCard = newCard(King, Hearts)
	if lastCard != expectedCard {
		t.Errorf("Expected card to be %v but found %v instead", expectedCard, lastCard)
	}

	card := deck.Cards[26]
	expectedCard = newCard(Ace, Clubs)
	if card != expectedCard {
		t.Errorf("Expected card to be %v but found %v instead", expectedCard, card)
	}
}
