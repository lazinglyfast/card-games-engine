package deck

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDraw(t *testing.T) {
	deck := NewDefaultDeck()
	drawnCards := deck.Draw(1)
	if len(drawnCards) != 1 {
		msg := "Expected to have drawn 1 card, instead drew %d"
		t.Errorf(msg, len(drawnCards))
	}

	actualCard := drawnCards[0]
	expectedCard := newCard(King, Hearts)
	if actualCard != expectedCard {
		msg := "Expected drawn card of a unshuffled deck to be %v. Found %v instead"
		t.Errorf(msg, expectedCard, actualCard)
	}

	// do not introduce non-simple logic otherwise we'll end up having to
	// write tests for tests. Keep it simple
	if deck.RemainingCardCount() != 51 {
		t.Errorf("Expected deck to have 51 cards after having drawn 1 card")
	}
}

func TestDrawFromEmptyDeck(t *testing.T) {
	deck := NewEmptyDeck()
	drawnCards := deck.Draw(10)
	if len(drawnCards) != 0 {
		msg := "Expected empty deck to have to cards, instead found %v"
		t.Errorf(msg, drawnCards)
	}
}

func TestDrawMoreCardsThanAvailable(t *testing.T) {
	cards := []Card{
		newCard(Ace, Spades),
		newCard(King, Hearts),
		newCard(Ace, Clubs),
	}

	deck := NewDeck(cards)

	drawnCards := deck.Draw(5)
	if len(drawnCards) != 3 {
		msg := "Expected at most 3 cards to be drawn since the deck had only 3 cards, instead found %v"
		t.Errorf(msg, len(drawnCards))
	}
}

func TestIsShuffled(t *testing.T) {
	unshuffledDeck := NewDefaultDeck()
	if unshuffledDeck.IsShuffled() {
		t.Errorf("Expected default deck to be unshuffled")
	}

	cards := []Card{
		newCard(Ace, Spades),
		newCard(King, Hearts),
		newCard(Ace, Clubs),
	}

	shuffledDeck := NewDeck(cards)
	if !shuffledDeck.IsShuffled() {
		t.Errorf("Expected deck to be shuffled")
	}
}

func TestShuffle(t *testing.T) {
	shuffledDeck := NewDefaultDeck()
	shuffledDeck.Shuffle()
	unshuffledDeck := NewDefaultDeck()
	// there's a 1/52! chance of this actually being true
	// instead of over-engineering let's ignore that
	// (after all that is beyond astronomically large)
	if cmp.Equal(shuffledDeck.Cards, unshuffledDeck.Cards) {
		t.Errorf("Expected deck to be shuffled")
	}
}

func TestUnshuffle(t *testing.T) {
	cards := []Card{
		newCard(Ace, Spades),
		newCard(King, Hearts),
		newCard(Ace, Clubs),
	}

	shuffledDeck := NewDeck(cards)
	if !shuffledDeck.IsShuffled() {
		t.Errorf("Expected deck %v to be shuffled", shuffledDeck)
	}

	shuffledDeck.unshuffle()

	if shuffledDeck.IsShuffled() {
		t.Errorf("Expected deck %v to be unshuffled", shuffledDeck)
	}

	cards = []Card{
		newCard(Ace, Spades),
		newCard(Ace, Clubs),
		newCard(King, Hearts),
	}
	expectedDeck := NewDeck(cards)

	if !cmp.Equal(shuffledDeck.Cards, expectedDeck.Cards) {
		msg := "Cards %v differ from expected %v"
		t.Errorf(msg, shuffledDeck.Cards, expectedDeck.Cards)
	}
}

func TestNewDeck(t *testing.T) {
	cards := []Card{
		newCard(Ace, Spades),
		newCard(King, Hearts),
		newCard(Ace, Clubs),
	}
	deck := NewDeck(cards)
	n := len(deck.Cards)
	if n != 3 {
		msg := "Expected default deck to have 3 cards but it has %d instead"
		t.Errorf(msg, n)
	}

	firstCard := deck.Cards[0]
	expectedCard := newCard(Ace, Spades)
	if firstCard != expectedCard {
		msg := "Expected card to be %v but found %v instead"
		t.Errorf(msg, expectedCard, firstCard)
	}

	secondCard := deck.Cards[1]
	expectedCard = newCard(King, Hearts)
	if secondCard != expectedCard {
		msg := "Expected card to be %v but found %v instead"
		t.Errorf(msg, expectedCard, secondCard)
	}

	thirdCard := deck.Cards[2]
	expectedCard = newCard(Ace, Clubs)
	if thirdCard != expectedCard {
		msg := "Expected card to be %v but found %v instead"
		t.Errorf(msg, expectedCard, thirdCard)
	}
}

func TestDefaultDeck(t *testing.T) {
	deck := NewDefaultDeck()
	n := len(deck.Cards)
	if n != 52 {
		msg := "Expected default deck to have 52 cards but it has %d instead"
		t.Errorf(msg, n)
	}

	firstCard := deck.Cards[0]
	expectedCard := newCard(Ace, Spades)
	if firstCard != expectedCard {
		msg := "Expected card to be %v but found %v instead"
		t.Errorf(msg, expectedCard, firstCard)
	}

	lastCard := deck.Cards[51]
	expectedCard = newCard(King, Hearts)
	if lastCard != expectedCard {
		msg := "Expected card to be %v but found %v instead"
		t.Errorf(msg, expectedCard, lastCard)
	}

	card := deck.Cards[26]
	expectedCard = newCard(Ace, Clubs)
	if card != expectedCard {
		msg := "Expected card to be %v but found %v instead"
		t.Errorf(msg, expectedCard, card)
	}
}
