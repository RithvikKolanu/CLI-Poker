package main

import "fmt"

func main() {
	newdeck := newDeck()
	newdeck = shuffle(newdeck)
	newdeck, c := deal(newdeck)
	hand1 := newHand()
	hand1 = append(hand1, c)
	fmt.Print("------\n")
	printDeck(hand1)
	fmt.Printf("length of deck: %d\n", len(newdeck))
	fmt.Printf("Length of Hand: %d\n", len(hand1))
}
