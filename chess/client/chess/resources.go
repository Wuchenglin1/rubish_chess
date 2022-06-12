package chess

import (
	"chess/client/img"
	_ "embed"
)

var ChessMap = map[int][]byte{
	ChessBoard: img.Chessboard,
	heiche:     img.Heiche,
	heijiang:   img.Heijiang,
	heima:      img.Heima,
	heipao:     img.Heipao,
	heishi:     img.Heishi,
	heixiang:   img.Heixiang,
	heizu:      img.Heizu,
	hongbing:   img.Hongbing,
	hongche:    img.Hongche,
	hongma:     img.Hongma,
	hongpao:    img.Hongpao,
	hongshi:    img.Hongshi,
	hongshuai:  img.Hongshuai,
	hongxiang:  img.Hongxiang,
}
