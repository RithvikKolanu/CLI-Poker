package main

import (
	"math/rand"
	"time"
)

type card struct {
	suit   string
	number string
}

func newDeck() []card {
	deck := make([]card, 0, 52)
	suittypes := []string{"clubs", "diamonds", "hearts", "spades"}
	cardtypes := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

	for _, s := range suittypes {
		for _, c := range cardtypes {
			deck = append(deck, card{s, c})
		}
	}
	return deck
}

func printDeck(deck []card) {
	for _, c := range deck {
		printCard(c)
	}
}

func shuffle(deck []card) []card {
	shuffled := make([]card, len(deck))
	rand.Seed(time.Now().UnixNano())
	random := rand.Perm(len(deck))
	for i, j := range random {
		shuffled[j] = deck[i]
	}
	return shuffled
}

func deal(deck []card) ([]card, card) {
	c := deck[0]
	for i := range deck {
		if i == len(deck)-1 {
			break
		}
		deck[i] = deck[i+1]
	}
	newdeck := deck[:len(deck)-1]
	return newdeck, c
}

func printCard(c card) string {
	cardprint := c.number + " of " + c.suit
	return cardprint
}

func newHand() []card {
	deck := make([]card, 0)
	return deck
}
