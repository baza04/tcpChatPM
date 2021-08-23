package main

import "net"

type room struct {
	name    string
	members map[net.Conn]*client
}
