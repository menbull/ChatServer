服务器职责：

ManagerServer:
负责管理其他服务器，与其他服务器形成一个集群，一个可用的闭环。
功能包括: 各所属服务器的启停统一控制与协调，查看各服务器日志，故障以及初步的管理功能(暂未定)。


{"ArgAmount":int,"Args":"string"}
SHUTDOWN
RESTART
根据游戏逻辑的不同再有不同的标准。。暂时还没有想好

2000端口是Manager
3000是数据库
5000端口是Gate
5001往后是Connector

都是子服务器连接父服务器
Manager<-LoginServer
       <-Connector<-Logic
          Database<-Logic
		
		
//Manager端消息类型
const (
	"SETUP"       //创建子服务器
	"UNAVAILABLE" //没有可用子服务器
	"STOP"        //停止子服务器
	//other
)

//其他子Server共用消息类型
const (
	"ONLINE"    //请求创建子服务器
	"HEARTBEAT" //子服务器心跳消息
	//other
)

