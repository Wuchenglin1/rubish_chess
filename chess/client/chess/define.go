package chess

type ResponseLoginMessage struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refreshToken"`
		Uid          int    `json:"uid"`
	} `json:"data"`
	Info   string `json:"info"`
	Status int    `json:"status"`
}
type ResponseRegisterMessage struct {
	Data   string `json:"data"`
	Info   string `json:"info"`
	Status int    `json:"status"`
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

//var ReadMsgChan = make(chan Msg, 100)
var SendMsgChan = make(chan Msg, 100)

const (
	// ChessBoard 棋盘
	ChessBoard = 1
	//红帅
	hongshuai = 8
	//红士
	hongshi = 9
	//红象
	hongxiang = 10
	//红马
	hongma = 11
	//红车
	hongche = 12
	//红炮
	hongpao = 13
	//红兵
	hongbing = 14
	//黑将
	heijiang = -8
	//黑士
	heishi = -9
	//黑相
	heixiang = -10
	//黑马
	heima = -11
	//黑车
	heiche = -12
	//黑炮
	heipao = -13
	//黑卒
	heizu = -14
)

//棋子位置
const (
	offsetX   = 48  //棋子x的偏移量
	increaseX = 104 //棋子x的增量
	offsetY   = 46  //棋子y的偏移量
	increaseY = 93  //棋子y的增量
	//棋子大小就是 x*y = 104 * 93
)

// 棋盘初始设置
var initBoard = [132]int{
	/*
		1 + 9 + 1
		+	    +
		10	   10
		+		1
		1 + 9 + 1
	*/

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
