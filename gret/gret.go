package gret

import (
	"github.com/gin-gonic/gin"
)

const (
	GRet   = "gbin_ret"
	Seqkey = "gMiddlewareSeq"
)

type RetM struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RetNo(code int64, msg string) RetM {
	return RetM{Code: code, Msg: msg}
}

func Ret(c *gin.Context, r RetM, data ...interface{}) {
	if r.Code == 0 {
		r.Msg = "success"
	}
	if len(data) > 0 && r.Code == 0 {
		r.Data = data[0]
	}
	c.Abort()
	c.Set(GRet, r)
	return
}
