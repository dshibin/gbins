//@author : bins
//@date : 2022/5/8 20:43

package greq

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"strings"
)

// 绑定json数据
func BindJson(ctx *gin.Context, req interface{}) (err error) {
	err = ctx.ShouldBindBodyWith(&req, binding.JSON)
	return
}

// 获取全部请求体 ，json 形式
func GetAllBodyJson(ctx *gin.Context) (ret string) {
	var body []byte
	if cb, ok := ctx.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		b, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}
		ctx.Set(gin.BodyBytesKey, b)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		body = b
	}
	ret = strings.ReplaceAll(string(body), "\r\n", "")
	ret = strings.ReplaceAll(ret, "\n", "")
	ret = strings.ReplaceAll(ret, " ", "")
	return
}
