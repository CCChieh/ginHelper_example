package service

import (
	"github.com/ccchieh/ginHelper"
	"github.com/gin-gonic/gin"
)

type Hello struct {
	ginHelper.Param        //内嵌ginHelper的基本Param
	Name            string `form:"name" binding:"required"`
}

func (param *Hello) Service() {
	param.Ret = gin.H{"message": "Hello " + param.Name + "!"}
}
