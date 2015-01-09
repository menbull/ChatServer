package main

import (
	"chatserver/Logger"
	"chatserver/structs"
	"encoding/json"
	"net"
	"os"
	"strings"
)

var managerClient *net.Conn
var clientListener *net.TCPListener
var serverList Structs.ServerList
var SERVER_NAME string

func main() {
	setLogger()
	Logger.Info("Starting Connector Server...")
	getServerConfig()
	setupManagerClient()
}

func setLogger() {
	Logger.SetConsole(true)
	Logger.SetRollingDaily("../logs", "Connector-logs.txt")
}

func getServerConfig() {
	serverConfig, err := os.Open("../config/servers.conf")
	defer serverConfig.Close()
	checkError(err)

	buf := make([]byte, 1024)
	length, err := serverConfig.Read(buf)
	checkError(err)

	err = json.Unmarshal(buf[:length], &serverList)
	checkError(err)

	Logger.Info("Get server config complete")
}

func setupManagerClient() {
	managerClient, err := net.Dial("tcp", serverList.Manager[0].Ip+serverList.Manager[0].Port)
	checkError(err)
	defer managerClient.Close()

	managerClient.Write([]byte("ONLINE|CONNECTOR_SERVER"))

	for {
		buffer := make([]byte, 512)
		length, err := managerClient.Read(buffer)
		checkError(err)

		cmd := strings.Split(string(buffer[:length]), "|")

		if cmd[0] == "STOP" {
			Logger.Info("Connector server closed")
			os.Exit(0)
		}
		if cmd[0] == "SETUP" {
			Logger.Info("Now my name is " + cmd[1])
			SERVER_NAME = cmd[1]
			Logger.Info("Listening port " + cmd[2])
			go setupClientListener(cmd[2])
		}
	}
}

func setupClientListener(p string) {
	laddr, err := net.ResolveTCPAddr("tcp", p)
	checkError(err)
	clientListener, err = net.ListenTCP("tcp", laddr)
	checkError(err)
	for {
		conn, err := clientListener.Accept()
		checkError(err)
		Logger.Info("A client from " + conn.RemoteAddr().String() + " is online")
	}
}

func checkError(err error) {
	if err != nil {
		Logger.Error(err.Error())
	}
}
