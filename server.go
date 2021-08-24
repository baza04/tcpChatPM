package main

import (
	"errors"
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
		default:
			fmt.Println("ERROR: unexpected commandID:", cmd.id)
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.writeMsg(fmt.Sprintf("Ok, I will call you as %s\n", args[1]))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok { // what about closing rooms?
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	s.quitCurrentRoom(c)

	r.members[c.conn.RemoteAddr()] = c
	c.room = r

	c.room.broadcast(c, fmt.Sprintf("%s has joined to the room", c.nick))
	c.writeMsg(fmt.Sprintf("Welcome to %s", roomName))
}

func (s *server) listRooms(c *client, args []string) {
	list := make([]string, 0, len(s.rooms))
	for name := range s.rooms {
		list = append(list, name)
		c.writeMsg(fmt.Sprintf("exist rooms:\n%s\n", strings.Join(list, "\n")))
	}
}

func (s *server) msg(c *client, args []string) {
	if c.room != nil {
		c.err(errors.New("you must join to the room first"))
		return
	}
	msg := fmt.Sprintf("[%s]: %s", c.nick, strings.Join(args[1:], " "))
	c.room.broadcast(c, msg)
}

func (s *server) quit(c *client, args []string) {
	log.Printf("client has disconnected %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.writeMsg("sad to see your go :(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
