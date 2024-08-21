package gmiddleware

import (
	"encoding/json"
	"github.com/dshibin/gbins/glog"
	"github.com/dshibin/gbins/gret"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func GTotalTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		c.Next()
		ret, _ := c.Get(gret.GRet)
		retStr, code := []byte{}, int64(200)
		switch ret.(type) {
		case gret.RetM:
			retS := ret.(gret.RetM)
			code = retS.Code
			retStr, _ = json.Marshal(ret)
			glog.Infoz(c, string(retStr), zap.String(glog.ConstTotalTime, time.Since(st).String()), zap.Int64(glog.ConstRet, code))
		default:
		}
	}
}
