package main

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	commands chan command
	members  map[net.Addr]*client
}

func newServer() *server {
	return &server{
		commands: make(chan command),
		members:  make(map[net.Addr]*client),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
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
	}

	c.readInput()
}

func (s *server) join(c *client, args []string) {
	s.members[c.conn.RemoteAddr()] = c
	c.msg("joined server")
	s.msg_all(c, fmt.Sprintf("joined server: %s", c.name))
}

func (s *server) fold(c *client, args []string) {
	log.Printf("Fold")
	s.msg_all(c, "folded")
}
func (s *server) check(c *client, args []string) {
	log.Printf("check")
	s.msg_all(c, "check")
}
func (s *server) raise(c *client, args []string) {
	log.Printf("raise")
	s.msg_all(c, "raise")
}
func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("name set to %s", c.name))
}
func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())
	c.conn.Close()
}

//logic for connecting with other clients
func (s *server) msg_all(sender *client, msg string) {
	for addr, cli := range s.members {
		if addr != sender.conn.RemoteAddr() {
			cli.msg(msg)
		}
	}
}
