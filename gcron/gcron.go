//@author : bins
//@date : 2022/10/27 10:59

package gcron

import (
	"github.com/dshibin/gbins/glog"
	"github.com/dshibin/gbins/grouter/gmiddleware"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var cronClient = cron.New(cron.WithSeconds())

// 定时任务管理器
func init() {
	cronClient.Start()
}

// 挂载定时任务
func AddCron(fc func(c *gin.Context), times string) {
	_, err := cronClient.AddFunc(times, DoCron(fc))
	if err != nil {
		panic("cron init error" + err.Error())
	}
}

// 执行挂载的定时任务
func DoCron(fc func(c *gin.Context)) func() {
	return func() {
		ctx := &gin.Context{}
		ctx.Set(gmiddleware.SeqKey, gmiddleware.CreateSeq())
		defer func() {
			err := recover()
			if err == nil {
				glog.Error(ctx, "DoCron err : ", err)
			}
		}()
		fc(ctx)
	}
}
