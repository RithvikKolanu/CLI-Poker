package main

import (
	"fmt"
	"log"
	"net"
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
}

func newServer() *server {
	return &server{
		commands:    make(chan command),
		members:     make(map[net.Addr]*client),
		memberorder: make([]*client, 0),
		deck:        newDeck(),
		pool:        0,
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_START:
			s.start(cmd.client, cmd.args)
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
	}

	c.readInput()
}

func (s *server) join(c *client, args []string) {
	s.members[c.conn.RemoteAddr()] = c
	s.memberorder = append(s.memberorder, c)
	s.msg_all(fmt.Sprintf("joined server: %s", c.name))
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
	log.Printf("check")
	s.msg_all(c.name + " checked")
	s.turn(c, IndexOf(s.memberorder, c)+1)
}
func (s *server) raise(c *client, args []string) {
	log.Printf("raise")
	s.msg_all(c.name + " raised")
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

func (s *server) start(c *client, args []string) {
	log.Printf("Starting game")
	s.msg_all("Starting game")
	s.deck = newDeck()
	s.deck = shuffle(s.deck)
	for _, cli := range s.members {
		cli.hand = newHand()
		dealcard := card{suit: "", number: ""}
		for i := 0; i < 3; i++ {
			s.deck, dealcard = deal(s.deck)
			cli.hand = append(cli.hand, dealcard)
		}
	}
}

func (s *server) turn(prev *client, index int) {
	fmt.Println(index)
	fmt.Println(len(s.memberorder))
	if index > len(s.memberorder)-1 {
		s.msg_all(s.memberorder[0].name + "'s turn")
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
