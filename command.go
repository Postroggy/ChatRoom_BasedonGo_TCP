package main

type commandID int

const (
	CmdNick commandID = iota //start from 0
	CmdJoin
	CmdRooms
	CmdMsg
	CmdQuit
	CmdMember
)

type command struct {
	id     commandID
	client *client
	args   []string
}
