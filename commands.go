package main

type commandID int

const (
	CMD_JOIN commandID = iota
	CMD_NAME
	CMD_FOLD
	CMD_CHECK
	CMD_RAISE
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
