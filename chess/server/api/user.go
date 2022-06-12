package api

import (
	"chess/server/model"
	"chess/server/service"
	"chess/server/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Register(c *gin.Context) {
	u := model.User{
		UserName: c.PostForm("userName"),
		Password: c.PostForm("password"),
	}
	if u.UserName == "" || len(u.UserName) < 6 || len(u.Password) > 20 {
		tool.RespDataWithErr(c, 400, "用户名长度不符合规范，请输入长度为6-20个字符的用户名")
		return
	}
	if u.Password == "" || len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RespDataWithErr(c, 400, "密码长度不符合规范，请输入长度为6-20个长度的密码")
		return
	}
	flag, err := service.SearchUserByName(&u)
	if err != nil {
		fmt.Println("searchUserByName error : ", err)
		tool.RespDataWithErr(c, 400, err)
		return
	}
	if flag {
		tool.RespDataWithErr(c, 400, "用户已存在")
		return
	}
	err = service.CreateUser(u)
	if err != nil {
		tool.RespDataWithErr(c, 400, err)
		return
	}
	tool.RespDataSuccess(c, 200, "注册成功")
}

func Login(c *gin.Context) {
	u := model.User{
		UserName: c.PostForm("userName"),
		Password: c.PostForm("password"),
	}
	if u.UserName == "" || len(u.UserName) < 6 || len(u.Password) > 20 {
		tool.RespDataWithErr(c, 400, "用户名长度不符合规范，请输入长度为6-20个字符的用户名")
		return
	}
	if u.Password == "" || len(u.Password) < 6 || len(u.Password) > 20 {
		tool.RespDataWithErr(c, 400, "密码长度不符合规范，请输入长度为6-20个长度的密码")
		return
	}
	id, flag, err := service.Login(&u)
	u.ID = id
	if err != nil {
		tool.RespDataWithErr(c, 400, err)
		return
	}
	if !flag {
		tool.RespDataWithErr(c, 400, "账号密码不相符或账号不存在！")
		return
	}
	//创建refreshToken和accessToken
	accessToken, err1 := service.CreateToken(u, "access_token", time.Hour*2)
	refreshToken, err2 := service.CreateToken(u, "refresh_token", time.Hour*24*2)
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
		tool.RespDataWithErr(c, 400, "创建token错误")
		return
	}
	tool.RespDataSuccess(c, 200, gin.H{
		"uid":          id,
		"access_token": accessToken,
		"refreshToken": refreshToken,
	})
}
