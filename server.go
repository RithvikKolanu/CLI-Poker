package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/chehsunliu/poker"
)

//Server has 4 fields
//	commands is the command input from each go routine
//	members is a map of net addresses to client pointers
//  deck is the deck of cards used by all clients
//  pool is the collection of money each round
type server struct {
	commands    chan command
	members     map[net.Addr]*client
	memberorder []*client
	deck        []card
	pool        int
	roundmaxbet int
	flop        []card
	rounds      int
}

func newServer() *server {
	return &server{
		commands:    make(chan command),
		members:     make(map[net.Addr]*client),
		memberorder: make([]*client, 0),
		deck:        newDeck(),
		pool:        0,
		roundmaxbet: 0,
		flop:        make([]card, 0),
		rounds:      0,
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_START:
			s.start()
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_FOLD:
			s.fold(cmd.client, cmd.args)
		case CMD_CHECK:
			s.check(cmd.client, cmd.args)
		case CMD_RAISE:
			s.raise(cmd.client, cmd.args)
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		case CMD_FLOP:
			c := cmd.client
			c.printCardsClient(s.flop)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		name:     "anon",
		commands: s.commands,
		hand:     newHand(),
		bankroll: 500,
		roundbet: 0,
		matched:  true,
	}

	c.readInput()
}

func (s *server) join(c *client, args []string) {
	s.members[c.conn.RemoteAddr()] = c
	s.msg_all(fmt.Sprintf("joined server: %s", c.name))
}

func (s *server) start() {
	log.Printf("Starting game")
	s.msg_all("\nStarting new round")
	s.deck = newDeck()
	s.deck = shuffle(s.deck)
	for _, cli := range s.members {
		s.memberorder = append(s.memberorder, cli)
		cli.roundbet = 0
		s.roundmaxbet = 0
		cli.hand = newHand()
		dealcard := card{suit: "", number: ""}
		for i := 0; i < 2; i++ {
			s.deck, dealcard = deal(s.deck)
			cli.hand = append(cli.hand, dealcard)
		}
	}
	s.flop = make([]card, 0)
	flopcard := card{suit: "", number: ""}
	for i := 0; i < 3; i++ {
		s.deck, flopcard = deal(s.deck)
		s.flop = append(s.flop, flopcard)
	}
	s.msg_all("The flop is: ")
	s.printDeck(s.flop)
	s.msg_all(s.memberorder[0].name + "'s turn")
}

//fold removes the player from the ordered list fo players
func (s *server) fold(c *client, args []string) {
	log.Printf("Fold")
	s.msg_all(c.name + " folded")
	i := IndexOf(s.memberorder, c)
	s.memberorder = append(s.memberorder[:i], s.memberorder[i+1:]...)
	fmt.Printf("%d", i)
	if i <= 0 {
		s.turn(s.memberorder[i], i)
	} else {
		s.turn(s.memberorder[i-1], i)
	}
}

func (s *server) check(c *client, args []string) {
	if !c.matched {
		c.msg("cannot check, must raise value to match: " + strconv.Itoa(s.roundmaxbet))
		return
	}
	log.Printf("check")
	s.msg_all(c.name + " checked")
	s.turn(c, IndexOf(s.memberorder, c)+1)
}

func (s *server) raise(c *client, args []string) {
	if len(args) == 1 {
		c.msg("Error, value was not entered, try again")
		return
	}
	betval, err := strconv.Atoi(args[1])
	if err != nil {
		c.msg("Error, value entered was not a number, try again")
		return
	}

	c.roundbet += betval
	if c.roundbet < s.roundmaxbet {
		c.msg("Error, value entered was not high enough to match the raise: " + args[1])
		c.roundbet -= betval
		return
	} else if c.roundbet > s.roundmaxbet {
		s.roundmaxbet = c.roundbet
		c.matched = true
		for _, cli := range s.memberorder {
			if cli != c {
				cli.matched = false
			}
		}
	} else if c.roundbet == s.roundmaxbet {
		c.matched = true
	}
	log.Printf("raise")
	s.msg_all(c.name + " raised a total of: " + strconv.Itoa(c.roundbet))
	s.turn(c, IndexOf(s.memberorder, c)+1)
}

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("name set to %s", c.name))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())
	s.msg_all(c.name + "has disconnected")
	c.conn.Close()
}

//logic for connecting with other clients
func (s *server) msg_all(msg string) {
	for _, cli := range s.members {
		cli.msg(msg)
	}
}

func (s *server) turn(prev *client, index int) {
	if index > len(s.memberorder)-1 {
		for i, cli := range s.memberorder {
			if !cli.matched {
				s.msg_all(s.memberorder[i].name + " has not matched the bet")
			}
		}
		for _, cli := range s.memberorder {
			if !cli.matched {
				s.msg_all(s.memberorder[0].name + " 's turn," + " everyone must raise to: " + strconv.Itoa(s.roundmaxbet))
				return
			}
		}
		//end of round logic
		if s.rounds >= 2 {
			cli := s.memberorder[0]
			clirank := s.handrank(s.memberorder[0])
			for i, n := range s.memberorder {
				if s.handrank(n) > clirank {
					clirank = s.handrank(n)
					cli = s.memberorder[i]
				}
			}
			for _, i := range s.memberorder {
				s.addToPot(i, i.roundbet)
			}
			s.msg_all("With rank: " + poker.RankString(clirank) + ", " + cli.name + " has won: " + strconv.Itoa(s.pool))
			s.collect(cli)

			s.start()

		} else {
			s.roundmaxbet = 0
			flopcard := card{suit: "", number: ""}
			s.deck, flopcard = deal(s.deck)
			s.flop = append(s.flop, flopcard)
			s.msg_all("The new flop is: ")
			s.printDeck(s.flop)
			fmt.Printf("%d\n", s.rounds)
			s.rounds++
			s.msg_all(s.memberorder[0].name + "'s turn")
		}

	} else if index < 0 {
		s.msg_all(s.memberorder[0].name + "'s turn")
	} else {
		s.msg_all(s.memberorder[index].name + "'s turn")
	}
}

func IndexOf(slice []*client, val *client) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}
