package tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func RespDataWithErr(c *gin.Context, status int, data interface{}) {
	fmt.Println(data)
	c.JSON(200, gin.H{
		"status": status,
		"info":   "failed",
		"data":   data,
	})
}

func RespDataSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(200, gin.H{
		"status": status,
		"info":   "success",
		"data":   data,
	})
}
