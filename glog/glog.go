package glog

import (
	"fmt"
	"github.com/dshibin/gbins/gconf"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger = &zap.Logger{}

const (
	Seqkey = "gMiddlewareSeq"

	coreConsole = "console"
	coreJson    = "json"

	ConstTotalTime = "costtime"
	ConstApp       = "app"
	ConstPath      = "path"
	ConstEnv       = "env"
	ConstSeq       = "seq"
	ConstRet       = "ret"
)

func init() {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/gbin.log",
		MaxSize:    20,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
	logger = zap.New(zapcore.NewTee(
		zapcore.NewCore(newEncode(coreJson), zapcore.AddSync(lumberJackLogger), zap.DebugLevel),
		zapcore.NewCore(newEncode(coreConsole), zapcore.AddSync(os.Stdout), zap.DebugLevel),
	), zap.AddCallerSkip(1), zap.AddCaller(), defaultFields())
	defer logger.Sync()
}

func newEncode(formatter string) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "l",
		NameKey:        "n",
		CallerKey:      "c",
		MessageKey:     "m",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	if formatter == coreJson {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	return encoder
}

func defaultFields() zap.Option {
	return zap.Fields(
		zap.String(ConstEnv, gconf.GConfig().Global.EnvName),
		zap.String(ConstApp, gconf.GConfig().Server.App),
	)
}

func Info(ctx *gin.Context, arg ...interface{}) {
	logger.Info(fmt.Sprint(arg...), zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))
	fmt.Print("")
}

func Infof(ctx *gin.Context, msg string, arg ...interface{}) {
	logger.Info(fmt.Sprintf(msg, arg...), zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))
	fmt.Print("")
}

func Infoz(ctx *gin.Context, msg string, arg ...zap.Field) {
	logger.Info(msg, append(arg, zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))...)
	fmt.Print("")
}

func Debug(ctx *gin.Context, arg ...interface{}) {
	logger.Debug(fmt.Sprint(arg...), zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))
	fmt.Print("")
}

func Debugf(ctx *gin.Context, msg string, arg ...interface{}) {
	logger.Debug(fmt.Sprintf(msg, arg...), zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))
	fmt.Print("")
}

func Debugz(ctx *gin.Context, msg string, arg ...zap.Field) {
	logger.Debug(msg, append(arg, zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))...)
	fmt.Print("")
}

func Error(ctx *gin.Context, arg ...interface{}) {
	logger.Error(fmt.Sprint(arg...), zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))
	fmt.Print("")
}

func Errorf(ctx *gin.Context, msg string, arg ...interface{}) {
	logger.Error(fmt.Sprintf(msg, arg...), zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))
	fmt.Print("")
}

func Errorz(ctx *gin.Context, msg string, arg ...zap.Field) {
	logger.Error(msg, append(arg, zap.String(ConstPath, getPath(ctx)), zap.String(ConstSeq, ctx.GetString(Seqkey)))...)
	fmt.Print("")
}

func getPath(ctx *gin.Context) (path string) {
	if ctx.Request != nil && ctx.Request.URL != nil && ctx.Request.URL.Path != "" {
		path = ctx.Request.URL.Path
	}
	return
}
