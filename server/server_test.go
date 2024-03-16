package main

import (
	"encoding/json"
	"example.com/deck"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateDeck(t *testing.T) {
	decks := make(map[uuid.UUID]deck.Deck)
	ctx := NewHandlerContext(&decks)

	t.Run("creates default deck", func(t *testing.T) {
		deck, _ := create(t, ctx, "/create")
		n := deck.RemainingCardCount
		if n != 52 {
			msg := "Expected default deck to have 52 cards, found %v instead"
			t.Errorf(msg, n)
		}

		if deck.IsShuffled {
			t.Errorf("Expected default deck to be unshuffled")
		}
	})

	t.Run("creates shuffled deck", func(t *testing.T) {
		deck, _ := create(t, ctx, "/create?shuffled=true")
		if !deck.IsShuffled {
			t.Errorf("Expected deck to be shuffled")
		}
	})

	t.Run("creates custom deck", func(t *testing.T) {
		deck, _ := create(t, ctx, "/create?cards=AS,KD,AC,2C,KH")
		n := deck.RemainingCardCount
		if n != 5 {
			msg := "Expected custom deck to have 5 cards, found %v instead"
			t.Errorf(msg, n)
		}
	})

	t.Run("fails to create deck with invalid card input", func(t *testing.T) {
		_, err := create(t, ctx, "/create?cards=A?,KD,AC,2C,KH")
		if err == nil {
			t.Errorf("Expected err due to invalid card input")
		}
	})
}

func create(t *testing.T, ctx *HandlerContext, url string) (CreatedDeck, error) {
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	ctx.Create(w, r)

	res := w.Result()
	defer res.Body.Close()
	jsonBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return CreatedDeck{}, err
	}

	var deck CreatedDeck
	if err := json.Unmarshal(jsonBytes, &deck); err != nil {
		return CreatedDeck{}, err
	}
	return deck, nil
}

func TestOpenDeck(t *testing.T) {
	decks := make(map[uuid.UUID]deck.Deck)
	ctx := NewHandlerContext(&decks)

	t.Run("opens deck", func(t *testing.T) {
		deck, _ := create(t, ctx, "/create")
		url := fmt.Sprintf("/open/%s", deck.Guid)
		openDeck, _ := open(t, ctx, url)
		n := openDeck.RemainingCardCount
		if n != 52 {
			msg := "Expected default open deck to have 52 cards, found %v instead"
			t.Errorf(msg, n)
		}

		if deck.IsShuffled {
			t.Errorf("Expected default open deck to be unshuffled")
		}
	})

	t.Run("fails to open non-existing deck", func(t *testing.T) {
		// astronomically unlikely to get the same guid (even more so than 1/52!)
		guid := "67e55044-10b1-426f-9247-bb680e5fe0c8"
		url := "/open/" + guid
		_, err := open(t, ctx, url)
		if err == nil {
			t.Errorf("Expected error opening a deck that does not exist")
		}
	})

	t.Run("fails in trying to open deck with invalid guid", func(t *testing.T) {
		// astronomically unlikely to get the same guid (even more so than 1/52!)
		url := fmt.Sprintf("/open/this-is-not-a-guid")
		_, err := open(t, ctx, url)
		if err == nil {
			t.Errorf("Expected error opening a deck with invalid guid")
		}
	})
}

func open(t *testing.T, ctx *HandlerContext, url string) (OpenDeck, error) {
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	ctx.Open(w, r)

	res := w.Result()
	defer res.Body.Close()
	jsonBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return OpenDeck{}, err
	}

	var deck OpenDeck
	if err := json.Unmarshal(jsonBytes, &deck); err != nil {
		return OpenDeck{}, err
	}
	return deck, nil
}

func TestDrawCards(t *testing.T) {
	decks := make(map[uuid.UUID]deck.Deck)
	ctx := NewHandlerContext(&decks)

	t.Run("draws one card by default if no count is passed in", func(t *testing.T) {
		deck, _ := create(t, ctx, "/create")

		url := fmt.Sprintf("/open/%s", deck.Guid)
		openDeck, _ := open(t, ctx, url)

		url = fmt.Sprintf("/draw/%s", openDeck.Guid)
		cards, _ := draw(t, ctx, url)
		n := len(cards)
		if n != 1 {
			msg := "Expected card draw with no count to have returned 1 card, found %v instead"
			t.Errorf(msg, n)
		}
	})

	t.Run("draws correct number of cards", func(t *testing.T) {
		deck, _ := create(t, ctx, "/create")

		openUrl := fmt.Sprintf("/open/%s", deck.Guid)
		openDeck, _ := open(t, ctx, openUrl)

		drawUrl := fmt.Sprintf("/draw/%s?count=3", openDeck.Guid)
		cards, _ := draw(t, ctx, drawUrl)
		n := len(cards)
		if n != 3 {
			msg := "Expected 3 cards to have been drawn. Found %v instead"
			t.Errorf(msg, n)
		}

		openDeck, _ = open(t, ctx, openUrl)
		if openDeck.RemainingCardCount != (52 - 3) {
			msg := "After having drawn 3 cards expected deck to have 49 cards. Found %v instead"
			t.Errorf(msg, n)
		}
	})
}

func draw(t *testing.T, ctx *HandlerContext, url string) ([]OpenCard, error) {
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	ctx.Draw(w, r)

	res := w.Result()
	defer res.Body.Close()
	jsonBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return []OpenCard{}, err
	}

	var cards []OpenCard
	if err := json.Unmarshal(jsonBytes, &cards); err != nil {
		return []OpenCard{}, err
	}
	return cards, nil
}

func TestContentsOfUrlPath(t *testing.T) {
	input := "http://127.0.0.1:8000/create/guid/?shuffled=true"
	output, _ := url.Parse(input)
	if output.Path != "/create/guid/" {
		t.Errorf("Expected /create/guid/ but found %v", output.Path)
	}
}

func TestExtractGuidFromUrlPath(t *testing.T) {
	path := "http://127.0.0.1:8000/draw/67e55044-10b1-426f-9247-bb680e5fe0c8/?count=10"
	actualGuid, _ := extractGuidFromUrlPath(path)
	expectedGuid, _ := uuid.Parse("67e55044-10b1-426f-9247-bb680e5fe0c8")
	if actualGuid != expectedGuid {
		msg := "Expected needle to equal %v but found %v instead"
		t.Errorf(msg, expectedGuid, actualGuid)
	}
}
