package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	rooms    map[string]*Room
	Commands chan Command
}

func NewServer() *Server {
	return &Server{
		rooms:    make(map[string]*Room),
		Commands: make(chan Command),
	}
}

func (s *Server) Run() {
	for cmd := range s.Commands {
		switch cmd.Id {
		case CmdNick:
			s.Nick(cmd.Client, cmd.Args)
		case CmdJoin:
			s.Join(cmd.Client, cmd.Args)
		case CmdRooms:
			s.Listrooms(cmd.Client)
		case CmdMsg:
			s.Msg(cmd.Client, cmd.Args)
		case CmdQuit:
			s.Quit(cmd.Client)
		case CmdMember:
			s.listMember(cmd.Client)
		}
	}
}

func (s *Server) NewClient(conn net.Conn) {
	log.Printf("new client has joined: %s\n\t\t\"＊┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄＊", conn.RemoteAddr().String())

	c := &Client{
		Conn:     conn,
		Nick:     "anonymous",
		Commands: s.Commands,
		Room:     nil,
	}

	c.Readinput()
}

func (s *Server) Nick(c *Client, args []string) {
	if len(args) < 2 {
		c.Msg("Nick is required. usage: /Nick NAME")
		return
	}

	c.Nick = args[1]
	c.Msg(fmt.Sprintf("all right, I will call you %s", c.Nick))
}

func (s *Server) Join(c *Client, args []string) {
	if c.Nick == "anonymous" {
		c.Msg("your name is `anonymous` plz change you Nick name. usage: /Nick NICK_NAME")
	}
	if len(args) < 2 {
		c.Msg("Room name is required. usage: /Join ROOM_NAME")
		return
	}

	roomName := args[1] //Room 的名字

	r, ok := s.rooms[roomName] //如果存在就可以直接加入
	if !ok {                   //如果不存在就自己创建一个
		r = &Room{
			name:    roomName,
			members: make(map[net.Addr]*Client), //ip地址对应每个用户
		}
		s.rooms[roomName] = r // 创建新房间完成
	}
	r.members[c.Conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.Room = r

	r.broadcast(c, fmt.Sprintf("%s joined the Room", c.Nick))

	c.Msg(fmt.Sprintf("welcome to %s :),enjoy your trip here", roomName))
}

func (s *Server) Listrooms(c *Client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	if rooms != nil {
		c.Msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ",\n")))
	} else {
		c.Msg(fmt.Sprintf("THERE IS NO ROOM AVAILABLE. CREATE YOUR ROOM, usage: /Join yourroomid"))
	}
}
func (s *Server) listMember(c *Client) {
	if c.Nick == "anonymous" {
		c.Msg("your name is `anonymous` plz change you Nick name. usage: /Nick NICK_NAME  (•◡•) /")
	}
	if c.Room == nil {
		s.Listrooms(c)
		c.Msg("please choose one of these rooms, or create a Room on your own. usage：/Join ROOM_NAME （っ＾▿＾）")
		return
	}
	var memberList []string
	mem := s.rooms[c.Room.name].members
	for _, client := range mem {
		memberList = append(memberList, client.Nick)
	}
	c.Msg(fmt.Sprintf("member in this chat Room:\n%s", strings.Join(memberList, ",\n")))

}
func (s *Server) Msg(c *Client, args []string) {
	if c.Nick == "anonymous" {
		c.Msg("your name is `anonymous` plz change you Nick name. usage: /Nick NICK_NAME  (•◡•) /")
	}
	if c.Room == nil {
		c.Msg("you haven't Join any of the rooms.")
		s.Listrooms(c)
		c.Msg("please choose one of these rooms, or create a Room on your own. usage：/Join ROOM_NAME （っ＾▿＾）")
		return
	}
	if len(args) < 2 {
		c.Msg("message is required, usage: /Msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.Room.broadcast(c, c.Nick+": "+msg)
}

func (s *Server) Quit(c *Client) {
	log.Printf("client has left the chat: %s", c.Conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.Msg("SEE YOU AGAIN ٩(˘◡˘)۶\n" +
		"＊┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄＊")
	c.Conn.Close()
}

func (s *Server) quitCurrentRoom(c *Client) {
	if c.Room != nil { //如果当前已经加入了房间，则可以退出当前房间
		oldRoom := s.rooms[c.Room.name]
		delete(s.rooms[c.Room.name].members, c.Conn.RemoteAddr()) //删除用户
		log.Printf("client has left: %s\n\t\t\"＊┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄┅┄＊", c.Nick)
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the Room", c.Nick))
	}
}
