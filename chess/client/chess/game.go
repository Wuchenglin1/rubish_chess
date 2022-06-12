package chess

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
)

type Game struct {
	Selected   []int                 //选中的棋子
	NextStep   []int                 //目标位置
	MvLast     int                   //上一步移动的棋子
	IsGameOver bool                  //是否游戏结束
	images     map[int]*ebiten.Image //图片资源
	Turn       int                   //轮到谁走
	Position   [132]int              //棋盘上的棋子
	Size       int
}

func NewGame() bool {
	var userName string
	var password string
	var choice uint
	game := &Game{
		images:   make(map[int]*ebiten.Image),
		Position: initBoard,
	}
	if ok := game.loadResource(); !ok {
		return false
	}

	//设置窗口名称
	ebiten.SetWindowTitle("中国象棋")
	//设置窗口大小
	ebiten.SetWindowSize(768, 768)
	ebiten.SetMaxTPS(140)
	fmt.Println("请选择")
	fmt.Printf("1.注册\n2.登录\n")
	for _, err := fmt.Scan(&choice); err != nil && choice > 2; {
		fmt.Println("您的输入有误，请重新选择！")
	}

	fmt.Println("请输入账号密码")

	switch choice {
	case 1:
		for _, err1 := fmt.Scan(&userName, &password); err1 != nil || userName == "" || password == ""; {
			fmt.Println("您的输入有误，请重新输入")
		}
		err1 := Register(userName, password)
		if err1 != nil {
			return false
		} else {
			fmt.Println("注册成功！请登录吧~")
			for _, err1 = fmt.Scan(&userName, &password); err1 != nil; {
				fmt.Println("您的输入有误，请重新输入")
			}
			fmt.Println("请输入您想加入的房间号")
			for _, err1 = fmt.Scan(&roomNum); err1 != nil; {
				fmt.Println("输入有误，请重新输入")
			}
			err1 = Login(userName, password)
			if err1 != nil {
				return false
			} else {
				break
			}
		}
	case 2:
		for _, err1 := fmt.Scan(&userName, &password); err1 != nil; {
			fmt.Println("您的输入有误，请重新输入")
		}
		fmt.Println("请输入您想加入的房间号")
		for _, err1 := fmt.Scan(&roomNum); err1 != nil; {
			fmt.Println("输入有误，请重新输入")
		}

		err1 := Login(userName, password)
		if err1 != nil {
			return false
		} else {
			break
		}
	}
	go Connect()

	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("rungame error : %v", err)
		return false
	}
	return true
}

//loadResource 加载资源
func (g *Game) loadResource() bool {
	for k, v := range ChessMap {
		//因为这里棋盘文件有点特殊，是jpg格式文件，需要特殊处理
		if k == ChessBoard {
			img, err := jpeg.Decode(bytes.NewReader(v))
			if err != nil {
				log.Printf("load chessBoard error : %v", err)
				return false
			}
			_Image := ebiten.NewImageFromImage(img)
			g.images[k] = _Image
			continue
		}
		//加载图片,其他图片均为png
		img, _, err := image.Decode(bytes.NewReader(v))
		if err != nil {
			log.Printf("load img error : %v", err)
			return false
		}
		_Image := ebiten.NewImageFromImage(img)
		g.images[k] = _Image
	}
	return true
}

// LoadBoard 加载棋盘
func (g *Game) LoadBoard(screen *ebiten.Image) {
	//加载棋盘
	if v, ok := g.images[ChessBoard]; ok {
		screen.DrawImage(v, nil)
	}
	//加载棋子
	for x := 0; x < 10; x++ {
		for y := 0; y < 9; y++ {
			xPos := offsetX + y*increaseX //棋子相对x坐标
			yPos := offsetY + x*increaseY //棋子相对y坐标
			//加载棋子
			g.LoadSource(screen, g.images[g.Position[x*9+y]], xPos, yPos)
		}
	}
}

// LoadSource 与上之相连的就是加载图片资源
func (g *Game) LoadSource(screen, img *ebiten.Image, x, y int) {
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//视窗大小
	return 1024, 1024
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.IsGameOver {
			g.IsGameOver = false
		} else {
			if len(g.Selected) == 0 {
				g.Selected = make([]int, 2)
				x, y := ebiten.CursorPosition()
				x = (x - offsetX) / increaseX
				y = (y - offsetY) / increaseY
				g.Selected[0] = x
				g.Selected[1] = y
				fmt.Println("已选中第一步位置", x, y, g.Position[x*9+y])
			} else {
				g.NextStep = make([]int, 2)
				x, y := ebiten.CursorPosition()
				x = (x - offsetX) / increaseX
				y = (y - offsetY) / increaseY
				g.NextStep[0] = x
				g.NextStep[1] = y
				fmt.Println("已选中第二步位置", x, y, g.Position[x*9+y])
			}
		}
		//选择了第一步和第二步后发送位置到服务器，让服务器为其转发消息
		if len(g.NextStep) != 0 && len(g.Selected) != 0 {
			fmt.Println(g.Selected, g.NextStep)
			SendMsgChan <- Msg{
				Uid:      Uid,
				Cmd:      2,
				RoomNum:  roomNum,
				Info:     "",
				LastStep: g.Selected,
				NextStep: g.NextStep,
			}
			g.NextStep = nil
			g.Selected = nil
		}
		g.Position = UpdateBoard
		g.LoadBoard(screen)
	}
	//
	g.LoadBoard(screen)
}

// GetCoordinates 通过坐标获得
func GetCoordinates(x, y int) (float64, float64) {
	return float64(offsetX + y*increaseX), float64(offsetY + x*increaseY) //棋子相对y坐标
}
