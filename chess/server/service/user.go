package service

import (
	"chess/server/dao"
	"chess/server/model"
	"golang.org/x/crypto/bcrypt"
)

func SearchUserByName(u *model.User) (bool, error) {
	return dao.SearchUserByName(u)
}

func CreateUser(u model.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(password)
	return dao.CreateUser(u)
}

func Login(u *model.User) (uint, bool, error) {
	u1 := model.User{UserName: u.UserName}
	flag, err := dao.SearchUserByName(&u1)
	if err != nil {
		return 0, false, err
	}
	//账号不存在
	if !flag {
		return 0, false, nil
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u1.Password), []byte(u.Password)); err != nil {
		//密码错误
		return 0, false, nil
	} else {
		u = &u1
		return u.ID, true, nil
	}
}
