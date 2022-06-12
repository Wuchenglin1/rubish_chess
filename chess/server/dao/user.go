package dao

import (
	"chess/server/model"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// SearchUserByName 通过用户名查找用户
func SearchUserByName(u *model.User) (bool, error) {
	//先查redis中是否存在user
	val, err := Rdb.Get(context.Background(), u.UserName).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	data := strings.Split(val, ";")
	fmt.Println(data)
	_int, err := strconv.Atoi(data[0])
	u.ID = uint(_int)
	u.Password = data[1]
	//redis中存在user，返回true,nil
	//redis中不存在user
	//再查mysql
	db := Db.Where("user_name = ?", u.UserName).First(&u)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			//mysql中不存在user返回false,nil
			return false, nil
		} else {
			//err不为nil
			return false, db.Error
		}
	}
	_data := strconv.Itoa(int(u.ID)) + ";" + u.Password
	//mysql有，更新redis中的数据
	sCmd := Rdb.Set(context.Background(), u.UserName, _data, 0)
	if sCmd.Err() != nil {
		//err不为nil返回false,err
		return false, sCmd.Err()
	}
	//mysql中查到了user，返回true,nil
	return true, nil
}

// CreateUser 创建用户
func CreateUser(u model.User) error {
	//mysql中创建用户
	tx := Db.Begin()
	txDB := tx.Create(&u)
	if txDB.Error != nil {
		tx.Rollback()
		return txDB.Error
	}
	data := strconv.Itoa(int(u.ID)) + ";" + u.Password
	cmd := Rdb.Set(context.Background(), u.UserName, data, 0)
	if cmd.Err() != nil {
		tx.Rollback()
		return cmd.Err()
	}
	tx.Commit()
	return nil
}
