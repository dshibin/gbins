//@author : bins
//@date : 2023/4/3 11:40

package gmiddleware

import (
	"context"
	"github.com/dshibin/gbins/gconf"
	"github.com/dshibin/gbins/gret"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		if gconf.GConfig().Server.Timeout == 0 {
			gconf.GConfig().Server.Timeout = 3000
		}

		ctx, cancel := context.WithTimeout(c, time.Duration(gconf.GConfig().Server.Timeout)*time.Millisecond)
		c.Request = c.Request.WithContext(ctx)
		r := gret.TimeOut
		finish := make(chan struct{}, 1)

		defer func() {
			cancel()
		}()

		// 子协程
		go func() {
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			c.Abort()
			gret.Ret(c, r)
		case <-finish:
		}

		ret, exist := c.Get(gret.GRet)
		if !exist {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": r.Code, "msg": r.Msg, "seq": c.GetString(gret.Seqkey)})
			return
		}
		r, ok := ret.(gret.RetM)
		if !ok {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": r.Code, "msg": r.Msg, "seq": c.GetString(gret.Seqkey)})
			return
		}
		if r.Data == nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": r.Code, "msg": r.Msg, "seq": c.GetString(gret.Seqkey)})
		} else {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": r.Code, "msg": r.Msg, "data": r.Data, "seq": c.GetString(gret.Seqkey)})
		}
	}
}
