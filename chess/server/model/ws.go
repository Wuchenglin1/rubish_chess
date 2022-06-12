package model

import "github.com/gorilla/websocket"

//节点信息
type Node struct {
	User  User
	Conn  *websocket.Conn
	Send  chan []byte
	Heart int //0:断开心跳,1:开始发送ping,2:接收pong
	Tag   int //标记的序号,0为黑方，1为红方
}

//房间信息
type RoomInfo struct {
	Turn       int      `json:"turn"` //0或1 0是黑方，1是红方
	Room       [2]*Node `json:"room"` //两个玩家的uid
	IsGameOver bool     `json:"isGameOver"`
	BoardInfo  [132]int `json:"boardInfo"`
}

type Msg struct {
	Uid        uint     `json:"uid"`
	Cmd        int64    `json:"cmd"`
	RoomNum    int      `json:"roomNum"` //房间号
	Info       string   `json:"info"`
	IsGameOver bool     `json:"isGameOver"`
	LastStep   []int    `json:"step"`
	NextStep   []int    `json:"step2"`
	BoardInfo  [132]int `json:"boardInfo"`
}

const (
	//红帅
	Hongshuai = 8
	//红士
	Hongshi = 9
	//红相
	Hongxiang = 10
	//红马
	Hongma = 11
	//红车
	Hongche = 12
	//红炮
	Hongpao = 13
	//红兵
	Hongbing = 14
	//黑将
	Heijiang = -8
	//黑士
	Heishi = -9
	//黑相
	Heixiang = -10
	//黑马
	Heima = -11
	//黑车
	Heiche = -12
	//黑炮
	Heipao = -13
	//黑卒
	Heizu = -14
)
