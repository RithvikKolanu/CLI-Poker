package main

type commandID int

const (
	CMD_JOIN commandID = iota
	CMD_NAME
	CMD_FOLD
	CMD_CHECK
	CMD_RAISE
	CMD_QUIT
	CMD_START
)

//Three fields
//	id: is the type of command
//	client: is the client that sent the command
//  args: is the text the client entered
type command struct {
	id     commandID
	client *client
	args   []string
}
