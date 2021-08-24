package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	reader := bufio.NewReader(c.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = strings.Trim(msg, " \r\n")
		fmt.Println(msg)

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		}
	}

}

func (c *client) err(err error) {
	c.conn.Write([]byte(fmt.Sprintf("ERROR: %s\n", err.Error())))
}

func (c *client) writeMsg(msg string) {
	c.conn.Write([]byte(fmt.Sprintf("[%s]: %s\n", c.nick, msg)))
}
