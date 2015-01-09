package main

import (
	"ChatServer/Logger"
	"ChatServer/Servers"
	"encoding/json"
	"net"
	"os"
	"strings"
)

var serverList Servers.ServerList
var managerServer *net.TCPListener

func main() {
	setLogger()
	Logger.Info("Starting Manager Server...")
	getServerConfig()
	setupManagerServer()
}

func setLogger() {
	Logger.SetConsole(true)
	Logger.SetRollingDaily("../logs", "ManagerServer-logs.txt")
}

func getServerConfig() {
	serverConfig, err := os.Open("../config/Config.conf")
	defer serverConfig.Close()
	checkError(err)

	buf := make([]byte, 1024)
	length, err := serverConfig.Read(buf)
	checkError(err)

	err = json.Unmarshal(buf[:length], &serverList)
	checkError(err)

	Logger.Info("Get server config complete")
	return
}

func setupManagerServer() {
	address := serverList.ManagerSv[0].Ip + serverList.ManagerSv[0].Port
	addr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)

	managerServer, err = net.ListenTCP("tcp", addr)
	checkError(err)

	Logger.Info("Manager Server Start Success")
	Logger.Info("Manager Server name: " + serverList.ManagerSv[0].Name)
	Logger.Info("Listening ip: " + serverList.ManagerSv[0].Ip)
	Logger.Info("Listening port" + serverList.ManagerSv[0].Port)

	for {
		conn, err := managerServer.Accept()
		checkError(err)
		Logger.Info("Manager Server Accepted a new connection. A server from " + conn.RemoteAddr().String())
		go dealManagerServer(conn)
	}
}

func dealManagerServer(conn net.Conn) {
	defer conn.Close()

	var serverType int
	var serverName, ip, port string
	var id int

	for {
		buffer := make([]byte, 512)
		length, err := conn.Read(buffer)

		if err != nil {
			defer Logger.Info("Disconnected from " + serverName + " " + conn.RemoteAddr().String())
			switch serverType {
			case Servers.LOGIN_SERVER:
				serverList.LoginSv[id].IsAvailable = false
				serverList.LoginSv[id].Conn = nil
			case Servers.MSG_SERVER:
				serverList.MsgSv[id].IsAvailable = false
				serverList.MsgSv[id].Conn = nil
			case Servers.ROUTE_SERVER:
				serverList.RouteSv[id].IsAvailable = false
				serverList.RouteSv[id].Conn = nil
			case Servers.LOGIC_SERVER:
				serverList.LogicSv[id].IsAvailable = false
				serverList.LogicSv[id].Conn = nil
			}
			return
		}

		cmd := strings.Split(string(buffer[:length]), "|")

		if cmd[0] == "ONLINE" {
			if cmd[1] == "LOGIN_SERVER" {
				Logger.Info("Login Server Setup Request Arrived.")
				serverType = Servers.LOGIN_SERVER
				id, serverName, ip, port = findFreeServer(Servers.LOGIN_SERVER)
				if id != -1 {
					serverList.LoginSv[id].Conn = conn
					conn.Write([]byte("SETUP|" + serverName + "|" + ip + "|" + port))
					Logger.Info("Login Server Setup Info. name:" + serverName + " ip:" + ip + " port" + port)
				} else {
					conn.Write([]byte("UNAVAILABLE"))
					Logger.Warn("No Available Free Login Server")
					return
				}
			}
			if cmd[1] == "MSG_SERVER" {
				Logger.Info("Msg Server Setup Request Arrived.")
				serverType = Servers.MSG_SERVER
				id, serverName, ip, port = findFreeServer(Servers.MSG_SERVER)
				if id != -1 {
					serverList.MsgSv[id].Conn = conn
				}
			}
			if cmd[1] == "ROUTE_SERVER" {
				Logger.Info("Route Server Setup Request Arrived.")
				serverType = Servers.ROUTE_SERVER
				id, serverName, ip, port = findFreeServer(Servers.ROUTE_SERVER)
				if id != -1 {
					serverList.RouteSv[id].Conn = conn
				}
			}
			if cmd[1] == "LOGIC_SERVER" {
				Logger.Info("Logic Server Setup Request. Logic Server addr: " + conn.RemoteAddr().String())
				serverType = Servers.LOGIC_SERVER
				id, serverName, ip, port = findFreeServer(Servers.LOGIC_SERVER)
				if id != -1 {
					serverList.LogicSv[id].Conn = conn
				}
			}
		}
	}
}

func findFreeServer(serverType int) (int, string, string, string) {
	switch serverType {
	case Servers.LOGIN_SERVER:
		for i := 0; i < len(serverList.LoginSv[:]); i++ {
			if serverList.LoginSv[i].Name != "" && serverList.LoginSv[i].IsAvailable == false {
				serverList.LoginSv[i].IsAvailable = true
				return i, serverList.LoginSv[i].Name, serverList.LoginSv[i].Ip, serverList.LoginSv[i].Port
			}
		}
	case Servers.MSG_SERVER:
		for i := 0; i < len(serverList.MsgSv[:]); i++ {
			if serverList.MsgSv[i].Name != "" && serverList.MsgSv[i].IsAvailable == false {
				serverList.MsgSv[i].IsAvailable = true
				return i, serverList.MsgSv[i].Name, serverList.MsgSv[i].Ip, serverList.MsgSv[i].Port
			}
		}
	case Servers.ROUTE_SERVER:
		for i := 0; i < len(serverList.RouteSv[:]); i++ {
			if serverList.RouteSv[i].Name != "" && serverList.RouteSv[i].IsAvailable == false {
				serverList.RouteSv[i].IsAvailable = true
				return i, serverList.RouteSv[i].Name, serverList.RouteSv[i].Ip, serverList.RouteSv[i].Port
			}
		}
	case Servers.LOGIC_SERVER:
		for i := 0; i < len(serverList.LogicSv[:]); i++ {
			if serverList.LogicSv[i].Name != "" && serverList.LogicSv[i].IsAvailable == false {
				serverList.LogicSv[i].IsAvailable = true
				return i, serverList.LogicSv[i].Name, serverList.LogicSv[i].Ip, serverList.LogicSv[i].Port
			}
		}
	}
	return -1, "ERROR", "ERROR", "ERROR"
}

func checkError(err error) {
	if err != nil {
		Logger.Error(err.Error())
	}
}
