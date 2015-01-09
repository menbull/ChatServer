package main

import (
	"chatserver/logger"
	"chatserver/structs"
	"encoding/json"
	"net"
)

var managerClient *net.Conn
var serverList Structs.ServerList
var SERVER_NAME string

func main() {
	setLogger()
	Logger.Info("Starting Channel Server...")
	getServerConfig()
	setupManagerClient()
}

func setLogger() {
	Logger.SetConsole(true)
	Logger.SetRollingDaily("../logs", "Channel-logs.txt")
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

	managerClient.Write([]byte("ONLINE|CHANNEL_SERVER"))

	for {
		buffer := make([]byte, 512)
		length, err := managerClient.Read(buffer)
		checkError(err)

		cmd := strings.Split(string(buffer[:length]), "|")

		if cmd[0] == "STOP" {
			Logger.Info("Channel server closed")
			os.Exit(0)
		}
		if cmd[0] == "SETUP" {
			Logger.Info("Now my name is " + cmd[1])
			SERVER_NAME = cmd[1]
			Logger.Info("Listening port " + cmd[2])
			go setupClientHandler(cmd[2])
		}
	}
}
