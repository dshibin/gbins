//@author : bins
//@date : 2024/8/22 16:32

package gerrs

import (
	"errors"
	"fmt"
)

type Error struct {
	Code  int64       `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	cause error
}

var (
	RetOK      = New(0, "success")
	RetUnKnown = New(-10001, "unknown error")
	RetTimeout = New(-10002, "time out")
	RetMoreReq = New(-10003, "Frequent network requests")
)

// 获取错误信息
func (e *Error) Error() string {
	if e == nil {
		return RetOK.(*Error).Msg
	}
	if e.cause != nil {
		return fmt.Sprintf("code: %d, msg: %s, data: %v , caused: %s", e.Code, e.Msg, e.Data, e.cause.Error())
	}
	return fmt.Sprintf("code: %d, msg: %s, data: %v", e.Code, e.Msg, e.Data)
}

// 获取错误码
func Code(e error) int64 {
	if e == nil {
		return RetOK.(*Error).Code
	}
	err, ok := e.(*Error)
	if !ok && !errors.As(e, &err) {
		return RetUnKnown.(*Error).Code
	}
	if err == nil {
		return RetOK.(*Error).Code
	}
	return err.Code
}

// 获取错误信息
func Msg(e error) string {
	if e == nil {
		return RetOK.(*Error).Msg
	}
	err, ok := e.(*Error)
	if !ok || !errors.As(e, &err) {
		return RetUnKnown.(*Error).Msg
	}
	if err == nil {
		return RetOK.(*Error).Msg
	}
	return err.Msg
}

func New(code int64, msg string) error {
	return &Error{Code: code, Msg: msg}
}

func Wrap(err error, code int64, msg string) error {
	return &Error{Code: code, Msg: msg, Data: err, cause: err}
}

func (e *Error) Unwrap() error {
	return e.cause
}
