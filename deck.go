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
	suittypes := []string{"c", "d", "h", "s"}
	cardtypes := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	for _, s := range suittypes {
		for _, c := range cardtypes {
			deck = append(deck, card{s, c})
		}
	}
	return deck
}

func (s *server) printDeck(deck []card) {
	for _, c := range deck {
		s.msg_all(printCard(c))
	}
}

func shuffle(deck []card) []card {
	shuffled := make([]card, len(deck))
	rand.Seed(time.Now().UnixNano())
	random := rand.Perm(len(deck))
	for x := 0; x < 52; x++ {
		for i, j := range random {
			shuffled[j] = deck[i]
		}
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
	suits := ""
	numbers := ""
	switch c.suit {
	case "c":
		suits = "clubs"
	case "s":
		suits = "spades"
	case "d":
		suits = "diamonds"
	case "h":
		suits = "hearts"
	}
	switch c.number {
	case "K":
		numbers = "King"
	case "Q":
		numbers = "Queen"
	case "J":
		numbers = "Jack"
	case "A":
		numbers = "Ace"
	default:
		numbers = c.number
	}
	cardprint := "     " + numbers + " of " + suits
	return cardprint
}

func newHand() []card {
	deck := make([]card, 0)
	return deck
}
