package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

//Client has three fields:
//	conn is the connection
//	name is the username of the client
//	commands is a channel that passes a command type to the server
//		commands is input only
type client struct {
	conn     net.Conn
	name     string
	commands chan<- command
	hand     []card
	bankroll int
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/start":
			c.commands <- command{
				id:     CMD_START,
				client: c,
				args:   args,
			}
		case "/fold":
			c.commands <- command{
				id:     CMD_FOLD,
				client: c,
				args:   args,
			}
		case "/check":
			c.commands <- command{
				id:     CMD_CHECK,
				client: c,
				args:   args,
			}
		case "/raise":
			c.commands <- command{
				id:     CMD_RAISE,
				client: c,
				args:   args,
			}
		case "/name":
			c.commands <- command{
				id:     CMD_NAME,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		case "/hand":
			c.printHand()
		default:
			log.Printf("not a valid command")
		}
	}
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *client) printHand() {
	for _, cards := range c.hand {
		c.msg(printCard(cards))
	}
}
