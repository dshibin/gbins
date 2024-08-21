package gredis

import (
	"context"
	"errors"
	"fmt"
	"gbins/glog"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type gredisLogHook struct {
}

const costTime = "gredis_costTime"

func (g *gredisLogHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	st := time.Now()
	c := ctx.(*gin.Context)
	c.Set(costTime, st)
	return c, nil
}

func (g *gredisLogHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	c := ctx.(*gin.Context)
	t, ok := c.Get(costTime)
	cost := time.Duration(0)
	if ok {
		cost = time.Since(t.(time.Time))
	}
	if cmd.Err() != nil {
		glog.Errorz(c, fmt.Sprintf("err : %s , %s", cmd.Err().Error(), fmt.Sprint(cmd.FullName(), cmd.Args())), zap.String(glog.ConstTotalTime, cost.String()))
	} else {
		glog.Debugz(c, fmt.Sprint(cmd.FullName(), cmd.Args()), zap.String(glog.ConstTotalTime, cost.String()))
	}
	if errors.Is(cmd.Err() , redis.Nil) {
		cmd.SetErr(nil)
	}
	return nil
}

func (g *gredisLogHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (g *gredisLogHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
