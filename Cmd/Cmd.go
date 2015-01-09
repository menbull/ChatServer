package Cmd

import (
	"net"
)

type ServerCommand struct {
	Args []string
}

type Connection struct {
	uid    int
	socket *net.Conn
}
