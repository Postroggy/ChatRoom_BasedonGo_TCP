package server
type commandID int

const (
	CmdNick commandID = iota //start from 0
	CmdJoin
	CmdRooms
	CmdMsg
	CmdQuit
	CmdMember
)

type Command struct {
	Id     commandID
	Client *Client
	Args   []string
}
