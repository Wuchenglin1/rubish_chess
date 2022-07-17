package api

import (
	"chess/server/model"
	"chess/server/service"
	"chess/server/tool"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var initBoard = [132]int{
	20, 19, 18, 17, 16, 17, 18, 19, 20,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 21, 0, 0, 0, 0, 0, 21, 0,
	22, 0, 22, 0, 22, 0, 22, 0, 22,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	14, 0, 14, 0, 14, 0, 14, 0, 14,
	0, 13, 0, 0, 0, 0, 0, 13, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	12, 11, 10, 9, 8, 9, 10, 11, 12,
}

var onlineUserMap sync.Map //用来存储所有在线的用户
var RM = sync.Map{}

func WS(c *gin.Context) {
	//获取用户token
	token := c.Query("token")
	roomNum := c.Query("room")
	if token == "" || roomNum == "" {
		tool.RespDataWithErr(c, 400, "token或者room不能为空")
		return
	}
	//解析token
	u, err := service.ParseToken(token)
	if err != nil {
		tool.RespDataWithErr(c, 400, err)
		return
	}
	u.RoomNum, err = strconv.Atoi(roomNum)
	if err != nil {
		tool.RespDataWithErr(c, 400, "房间号有误")
		return
	}
	if u.RoomNum > 20 || u.RoomNum < 0 {
		tool.RespDataWithErr(c, 400, "房间号不存在！(0-19)")
		return
	}

	//升级为websocket协议
	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		tool.RespDataWithErr(c, 400, err)
		return
	}
	//发送心跳ping
	node := &model.Node{
		User:  u,
		Conn:  conn,
		Send:  make(chan []byte),
		Heart: 1,
	}
	//把用户存入到map里面
	onlineUserMap.Store(u.ID, node)
	rim := InitChess()
	for k := range rim.Room {
		if rim.Room[k] == nil {
			node.Tag = k
			rim.Room[k] = node
		}
	}
	RM.Store(u.RoomNum, rim)
	go ListenProc(node)
	go RecProc(node)
	go heartBeat(conn, u.ID)
	fmt.Println(u.UserName, "加入了房间", u.RoomNum)
}

func InitChess() *model.RoomInfo {
	rim := new(model.RoomInfo)
	rim.Turn = 1
	rim.Turn = 1
	rim.IsGameOver = false
	rim.BoardInfo = initBoard
	return rim
}

// ListenProc 监听玩家的操作
func ListenProc(node *model.Node) {
	defer node.Conn.Close()
	for {
		select {
		case data := <-node.Send:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// RecProc 接收心跳和其他消息
func RecProc(node *model.Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			//这里一般是客户端关闭发来的消息
			fmt.Println(node.User.UserName, "离开了游戏")
			//清楚在线列表里的该玩家
			onlineUserMap.Delete(node.User.ID)
			node.Conn.Close()
			return
		}
		//处理数据
		ParseData(data)
	}
}

//心跳检测
func heartBeat(conn *websocket.Conn, uid uint) {
	value, _ := onlineUserMap.Load(uid)
	node, ok := value.(*model.Node)
	if !ok {
		panic("read userInfo error")
		return
	}
	for {
		node.Heart = 1
		time.Sleep(time.Second * 5)
		err := sendPing(conn)
		if err != nil {
			//玩家退出
			fmt.Println("send heart stop")
			node.Heart = 0
			return
		}
	}
}

func sendPing(conn *websocket.Conn) error {
	marshal, err := json.Marshal(&model.Msg{Cmd: 0})
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, marshal)
}

func ParseData(data []byte) {
	msg := model.Msg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("parse data error : %v", err.Error())
		return
	}
	switch msg.Cmd {
	case 0:
		//心跳
		value, _ := onlineUserMap.Load(msg.Uid)
		node, ok := value.(*model.Node)
		if ok {
			//收到心跳
			node.Heart = 2
		}
		return
	case 2:
		//接收坐标位置
		//先读取玩家信息
		value, ok := onlineUserMap.Load(msg.Uid)
		if !ok {
			fmt.Printf("load user error1 : %v", err)
			return
		}
		node, _ := value.(*model.Node)
		//读取棋盘信息
		value, ok = RM.Load(msg.RoomNum)
		var rim *model.RoomInfo
		if ok {
			rim = value.(*model.RoomInfo)
		} else {
			fmt.Println("没有该指针")
			return
		}
		//服务端处理下棋的数据
		Move(node, rim, msg.LastStep, msg.NextStep)
		fmt.Println("收到来自 ", node.User.UserName, " 的坐标位置：", msg.LastStep, msg.NextStep)
	case 10:
		value, ok := onlineUserMap.Load(msg.Uid)
		if !ok {
			log.Fatalf("load player error :%v", err)
		}
		node, ok := value.(*model.Node)
		if !ok {
			log.Fatalf("load plyer error : %v", err)
		}
		fmt.Println(node.User.UserName, "说:", msg.Info)
		//普通聊天
	case 20:
		//悔棋
	}
}

//Move 各种象棋移动的逻辑处理
func Move(node *model.Node, rmi *model.RoomInfo, lastStep, nextStep []int) {

	startType := CheckPosition(rmi, lastStep[0], lastStep[1])
	endType := CheckPosition(rmi, nextStep[0], nextStep[1])
	fmt.Println("startType : ", startType)
	fmt.Println("endType : ", endType)
	if nextStep[0] == lastStep[0] && nextStep[1] == lastStep[1] {
		SendMessage(node, 9, "不要原地踏步哦", rmi.BoardInfo)
		return
	}
	if startType == 0 {
		SendMessage(node, 9, "无棋子，请重新指定", rmi.BoardInfo)
		return
	}
	//查看是否是当前应该走的一方的棋子
	flag := IsCheckSide(node, startType)
	if !flag {
		SendMessage(node, 9, "请选择自己的棋子喔~", rmi.BoardInfo)
		return
	}
	//查看是否是己方的棋子
	flag = IsOwnSideChess(startType, endType)
	if !flag {
		SendMessage(node, 9, "请不要乱走喔，目标点是己方棋", rmi.BoardInfo)
		return
	}
	tempX := nextStep[0] - lastStep[0]
	tempY := nextStep[1] - lastStep[1]
	switch startType {
	case model.Hongma:
		flag = Horse(rmi, nextStep[0], nextStep[1], tempX, tempY)
	case model.Heima:
		flag = Horse(rmi, nextStep[0], nextStep[1], tempX, tempY)
	case model.Hongche:
	case model.Heiche:
	case model.Hongpao:
	case model.Heipao:
	case model.Hongshi:
		flag = Shi(tempX, tempY, nextStep[0], nextStep[1])
	case model.Heishi:
		flag = Shi(tempX, tempY, nextStep[0], nextStep[1])
	case model.Hongbing:
		flag = Bing(startType, endType, lastStep[0], lastStep[1], nextStep[0], nextStep[1], tempX, tempY)
	case model.Heizu:
		flag = Bing(startType, endType, lastStep[0], lastStep[1], nextStep[0], nextStep[1], tempX, tempY)
	case model.Hongxiang:
		flag = Elephant(rmi, tempX, tempY, lastStep[0], lastStep[1], nextStep[0], nextStep[1])
	case model.Heixiang:
		flag = Elephant(rmi, tempX, tempY, lastStep[0], lastStep[1], nextStep[0], nextStep[1])
	case model.Hongshuai:
		flag = King(nextStep[0], nextStep[1], tempX, tempY)
	case model.Heijiang:
		flag = King(nextStep[0], nextStep[1], tempX, tempY)
	}
	Step(node, rmi, flag, lastStep, nextStep)

}

// CheckPosition 返回该坐标对应的棋类型
func CheckPosition(rmi *model.RoomInfo, x, y int) int {
	return rmi.BoardInfo[y*9+x]
}

//吃子
func Replace(rmi *model.RoomInfo, startX, startY, endX, endY int) {
	rmi.BoardInfo[endY*9+endX] = CheckPosition(rmi, startX, startY)
}

func ChangeTurn(turn int) int {
	return 1 - turn
}

func SendMessage(node *model.Node, cmd int64, data string, board [132]int) {
	msg := model.Msg{
		Uid:        node.User.ID,
		Cmd:        cmd,
		RoomNum:    node.User.RoomNum,
		Info:       data,
		IsGameOver: false,
		LastStep:   nil,
		NextStep:   nil,
		BoardInfo:  board,
	}
	marshal, err := json.Marshal(&msg)
	if err != nil {
		log.Fatalf("marshal error : %v", err)
	}
	node.Conn.WriteMessage(websocket.TextMessage, marshal)
}

//查看是否是当前应该走的一方的棋子
func IsCheckSide(node *model.Node, startType int) bool {
	if startType < 15 {
		//红
		if node.Tag == 0 {
			return false
		} else {
			return true
		}
	} else {
		//黑
		if node.Tag == 0 {
			return true
		} else {
			return false
		}
	}
}

//落子所在方格是否为己方子
func IsOwnSideChess(startType, endType int) bool {
	return startType*endType <= 0
}

func Horse(rim *model.RoomInfo, x, y, tempX, tempY int) bool {
	return tempX*tempX+tempY*tempY == 5 && CheckPosition(rim, x+tempX/2, y+tempY/2) == 0
}

func Elephant(rmi *model.RoomInfo, tempX, tempY, startX, startY, aminX, aminY int) bool {
	return CheckPosition(rmi, startX+tempX/2, startY+tempY/2) == 0 && tempX*tempX+tempY*tempY == 8 && startX/5 == aminX/5
}

func King(aimx, aimy, tempX, tempY int) bool {
	return aimx%7 >= 0 && aimx%7 <= 2 && aimy >= 3 && aimy <= 5 && tempX*tempX+tempY*tempY == 1
}

func Shi(tempX, tempY, aimX, aimY int) bool {
	return tempX*tempX+tempY*tempY == 2 && aimX%7 >= 0 && aimX%7 <= 2 && aimY >= 3 && aimY <= 5
}

func Bing(startType, endType, startX, startY, aimX, aimY, tempX, tempY int) bool {
	if math.Abs(float64(tempX)) == 1 && tempY <= 0 {
		if math.Abs(float64(tempX)) == 1 && tempY == 0 {
			if math.Abs(float64(tempY)) == 1 && tempX == 0 {
				if (startX/5 == 0 && startType > 0) || (startX/8 == 1 && startType < 0) {
					return true
				}
				return false
			}
		}
	}
	return false
}
func Pao(startX, startY, tempX, tempY, startType, endType int) bool {
	return false
}

func Car(startType, endType, startX, startY, tempX, tempY int) bool {
	return false
}

func Step(node *model.Node, rmi *model.RoomInfo, flag bool, lastStep []int, nextStep []int) {
	if flag {
		Replace(rmi, lastStep[0], lastStep[1], nextStep[0], nextStep[1])
	} else {
		SendMessage(node, 9, "违规处理操作！", rmi.BoardInfo)
	}
}
