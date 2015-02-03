package Servers

import (
	"net"
)

const (
	LOGIN_SERVER = iota
	MSG_SERVER
	ROUTE_SERVER
	LOGIC_SERVER
)

type ServerList struct {
	LoginSv   []Server
	MsgSv     []Server
	RouteSv   []Server
	LogicSv   []Server
	ManagerSv []Server
}

type Server struct {
	Name    string
	Ip      string
	Port    string
	IsUsing bool
	Conn    net.Conn
}
