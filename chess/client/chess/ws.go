package chess

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"strconv"
)

var (
	addr        = flag.String("addr", "110.42.184.72:8080", "http service address")
	token       string
	roomNum     int
	Message     Msg
	UpdateBoard [132]int
)

func Connect() {
	flag.Parse()
	log.SetFlags(0)

	sign := make(chan os.Signal, 1)
	room := strconv.Itoa(roomNum)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: "token=" + token + "&room=" + room}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial error %v", err)
	}
	go WriteMsg(c)
	go ReadMsg(c)
	for {
		select {
		case <-sign:
			log.Println("interrupt")
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Fatalf("write message error : %v", err)
			}
		}
	}
}

func ReadMsg(c *websocket.Conn) {
	defer c.Close()
	for {
		_, back, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("read message error : %v", err)
		}
		err = json.Unmarshal(back, &Message)
		if err != nil {
			fmt.Println(string(back))
			log.Fatalf("unmarshal message error : %v", err)
		}
		switch Message.Cmd {
		case 2:

		case 9:
			//系统消息
			UpdateBoard = Message.BoardInfo
			fmt.Println(Message.Info)
		case 10:
			//聊天消息
		}
	}
}

func WriteMsg(c *websocket.Conn) {
	defer c.Close()
	for {

		msg := <-SendMsgChan
		fmt.Println("已处理一条消息")
		marshal, err := json.Marshal(&msg)
		if err != nil {
			log.Fatalf("marshal message error : %v", err)
		}
		err = c.WriteMessage(websocket.TextMessage, marshal)
		if err != nil {
			log.Fatalf("send message error : %v", err)
			return
		}
	}
}
