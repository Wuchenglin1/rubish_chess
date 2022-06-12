package tool

import (
	"bufio"
	"chess/server/model"
	"encoding/json"
	"os"
)

var cfg *model.Cfg

func InitConfig() {
	//读取config文件
	file, err := os.Open("../config/chess_config/config.json")
	if err != nil {
		//直接报错
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *model.Cfg {
	return cfg
}
