package main

import (
	"ChatServer/Logger"
	"ChatServer/Servers"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var serverList Servers.ServerList
var managerClient *net.Conn
var loginServerListener *net.TCPListener

const HEARTBEAT_INTERVAL = 1

func main() {
	setLogger()
	Logger.Info("Starting Login Server...")
	getServerConfig()
	setupLoginServer()
}

func setLogger() {
	Logger.SetConsole(true)
	Logger.SetRollingDaily("../logs", "LoginServer-logs.txt")
}

func getServerConfig() {
	serverConfig, err := os.Open("../Config/Config.conf")
	defer serverConfig.Close()
	checkError(err)

	buf := make([]byte, 1024)
	length, err := serverConfig.Read(buf)
	checkError(err)

	err = json.Unmarshal(buf[:length], &serverList)
	checkError(err)

	Logger.Info("Get server config complete")
}

func setupLoginServer() {
	address := serverList.ManagerSv[0].Ip + serverList.ManagerSv[0].Port

	managerClient, err := net.Dial("tcp", address)
	checkError(err)
	defer managerClient.Close()

	managerClient.Write([]byte("ONLINE|LOGIN_SERVER"))
	Logger.Info("Request Login Server Setup Info. Manager Server addr: " + address)

	for {
		buffer := make([]byte, 512)
		length, err := managerClient.Read(buffer)
		if err != nil {
			defer Logger.Info("Disconnected from " + serverList.ManagerSv[0].Name + " " + managerClient.RemoteAddr().String())
			Logger.Info("Login Server Closed")
			exitServer(false)
		}

		cmd := strings.Split(string(buffer[:length]), "|")
		if cmd[0] == "STOP" {
			Logger.Info("Login Server Closed")
			exitServer(true)
		}
		if cmd[0] == "SETUP" {
			Logger.Info("Login Server Setup Info. name:" + cmd[1] + " ip:" + cmd[2] + " port:" + cmd[3])
			go createLoginServer(cmd[1], cmd[2], cmd[3])
		}
		if cmd[0] == "UNAVAILABLE" {
			Logger.Warn("No Available Free Login Server")
			exitServer(false)
		}
		//Todo:other operation
	}
}

func createLoginServer(name, ip, port string) {
	address := ip + port
	addr, err := net.ResolveTCPAddr("tcp", address)
	checkError(err)

	loginServerListener, err := net.ListenTCP("tcp", addr)
	checkError(err)

	Logger.Info("Login Server Start Success")
	Logger.Info("Login Server name: " + name)
	Logger.Info("Listening ip: " + ip)
	Logger.Info("Listening port" + port)

	for {
		conn, err := loginServerListener.Accept()
		checkError(err)
		Logger.Info("Login Server Accepted a new connection. A client from " + conn.RemoteAddr().String())
		go dealLoginServer(conn)

		//r := rand.Intn(len(serverList.Connector[:]))
		//conn.Write([]byte(serverList.Connector[r].Ip + ":" + serverList.Connector[r].Port))
	}
}

func sendHeartBeat() {
	Logger.Info(time.Now().Format("2006-01-02 15:04:05") + "----发送心跳")
	time.AfterFunc(HEARTBEAT_INTERVAL*time.Second, sendHeartBeat)
}

func dealLoginServer(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 512)
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		Logger.Error(err.Error())
	}
}

func exitServer(bDirectExit bool) {
	if bDirectExit {
		os.Exit(0)
	} else {
		fmt.Println("按回车键退出...")
		var str string
		fmt.Scanf("%v", &str)
		os.Exit(0)
	}
}
