package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Conn     net.Conn
	Nick     string
	Room     *Room
	Commands chan<- Command
}

func (c *Client) Readinput() {
	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.Commands <- Command{
				Id:     CmdNick,
				Client: c,
				Args:   args,
			}
		case "/join":
			c.Commands <- Command{
				Id:     CmdJoin,
				Client: c,
				Args:   args,
			}
		case "/rooms":
			c.Commands <- Command{
				Id:     CmdRooms,
				Client: c,
			}
		case "/msg":
			c.Commands <- Command{
				Id:     CmdMsg,
				Client: c,
				Args:   args,
			}
		case "/quit":
			c.Commands <- Command{
				Id:     CmdQuit,
				Client: c,
			}
		case "/member":
			c.Commands <- Command{
				Id:     CmdMember,
				Client: c,
			}
		default:
			c.err(fmt.Errorf("try again,unknown Command: %s", cmd))
		}
	}
}

func (c *Client) err(err error) {
	c.Conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *Client) Msg(msg string) {
	c.Conn.Write([]byte("> " + msg + "\n"))
}
