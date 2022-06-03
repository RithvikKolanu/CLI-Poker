package main

import (
	"fmt"

	"github.com/chehsunliu/poker"
)

func (s *server) handrank(c *client) int32 {
	totalhand := make([]card, 0)
	totalhand = append(totalhand, c.hand...)
	totalhand = append(totalhand, s.flop...)
	eval := make([]poker.Card, 0)
	for _, i := range totalhand {
		str := i.number + i.suit
		num := poker.NewCard(str)
		eval = append(eval, num)
	}
	fmt.Println(eval)

	rank := poker.Evaluate([]poker.Card(eval))
	return rank
}
