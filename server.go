package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CmdNick:
			s.nick(cmd.client, cmd.args)
		case CmdJoin:
			s.join(cmd.client, cmd.args)
		case CmdRooms:
			s.listRooms(cmd.client)
		case CmdMsg:
			s.msg(cmd.client, cmd.args)
		case CmdQuit:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s\n\t\t\"＊┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄＊", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
		room:     nil,
	}

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if c.nick == "anonymous" {
		c.msg("your name is `anonymous` plz change you nick name. usage: /nick NICK_NAME")
	}
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s＊┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄＊", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ",\n")))
}

func (s *server) msg(c *client, args []string) {
	if c.nick == "anonymous" {
		c.msg("your name is `anonymous` plz change you nick name. usage: /nick NICK_NAME  (•◡•) /")
	}
	if c.room == nil {
		s.listRooms(c)
		c.msg("please choose one of these rooms, or create a room on your own. usage：/join ROOM_NAME （っ＾▿＾）")
		return
	}
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("SEE YOU AGAIN ٩(˘◡˘)۶\n" +
		"＊┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄＊")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
