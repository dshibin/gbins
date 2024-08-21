package gdb

import (
	"encoding/json"
	"fmt"
	"github.com/dshibin/gbins/glog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"time"
)

const (
	callBackBeforeName = "callBackBeforeName"
	callBackAfterName  = "callBackAfterName"
	startTime          = "startTime"
)

type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "gbin_ormplugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:before_query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:before_update").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:before_row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:before_raw").Register(callBackBeforeName, before)

	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:after_row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:after_raw").Register(callBackAfterName, after)
	return
}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

func after(db *gorm.DB) {
	statement := db.Statement
	if statement == nil {
		return
	}

	_ctx := statement.Context.(*gin.Context)
	_ts, exist := db.InstanceGet(startTime)
	if !exist {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	sqlInfo, _ := json.Marshal(map[string]interface{}{
		"sql":   sql,
		"stack": utils.FileWithLineNum(),
		"Rows":  db.Statement.RowsAffected,
	})

	if statement.Error != nil {
		glog.Errorz(_ctx, fmt.Sprintf("err : %s , %s", statement.Error.Error(), string(sqlInfo)), zap.String(glog.ConstTotalTime, time.Since(ts).String()))
	} else {
		glog.Debugz(_ctx, string(sqlInfo), zap.String(glog.ConstTotalTime, time.Since(ts).String()))
	}
	return
}
