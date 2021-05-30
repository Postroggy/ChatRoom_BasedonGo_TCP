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
		case CmdMember:
			s.listMember(cmd.client)
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

	roomName := args[1] //room 的名字

	r, ok := s.rooms[roomName] //如果存在就可以直接加入
	if !ok {                   //如果不存在就自己创建一个
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client), //ip地址对应每个用户
		}
		s.rooms[roomName] = r // 创建新房间完成
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s :),enjoy your trip here", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ",\n")))
}
func (s *server) listMember(c *client) {
	if c.nick == "anonymous" {
		c.msg("your name is `anonymous` plz change you nick name. usage: /nick NICK_NAME  (•◡•) /")
	}
	if c.room == nil {
		s.listRooms(c)
		c.msg("please choose one of these rooms, or create a room on your own. usage：/join ROOM_NAME （っ＾▿＾）")
		return
	}
	var memberList []string
	mem := s.rooms[c.room.name].members
	for _, client := range mem {
		memberList = append(memberList, client.nick)
	}
	c.msg(fmt.Sprintf("member in this chat room: %s\n", strings.Join(memberList, ",\n")))

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
	if c.room != nil { //如果当前已经加入了房间，则可以退出当前房间
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
