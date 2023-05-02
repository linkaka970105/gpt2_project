package controllers

import (
	"fmt"
	"net/http"
)

// 通用响应码
const (
	Success   = 10000
	Err       = 10001
	ParamsErr = 10002
	TokenErr  = 10003
	NoAccess  = 10004
)

// 业务定义响应码，范围[12200，12299]。
const (
	PwdNotMatch = 12001
)

// 返回参数key
const (
	RespCodeKey = "code"
	RespMsgKey  = "message"
)

// RespMsg contains response message according to code.
var RespMsg = map[int]string{
	http.StatusOK:                  "成功",
	http.StatusNotFound:            "链接不存在",
	http.StatusInternalServerError: "服务器内部错误",
	Success:                        "成功",
	Err:                            "请求失败，请稍后再试",
	ParamsErr:                      "参数错误",
	TokenErr:                       "token不正确或者已经过期了",
	NoAccess:                       "没有访问权限",
	PwdNotMatch:                    "密码错误",
}

// Resp 业务返回结果
type Resp struct {
	Code int    `json:"code"`    // 业务返回状态码
	Msg  string `json:"message"` // 业务返回信息
}

// IsSucc return true if response success.
func (resp Resp) IsSucc() bool {
	return resp.Code == Success
}

// Err return `Error` struct contains code and msg.
func (resp Resp) Err() *ResError {
	return NewError(resp.Code, resp.Msg)
}

// New return a new result struct
func New() Resp {
	return Resp{Code: Success}
}

// ResError 业务响应错误
type ResError struct {
	Code int                    `json:"code"`    // 业务返回状态码
	Msg  string                 `json:"message"` // 业务返回信息
	Ext  map[string]interface{} `json:"-"`
}

// IsResError 判断错误类型
func IsResError(err error) bool {
	if _, ok := err.(ResError); ok {
		return true
	} else if _, ok := err.(*ResError); ok {
		return true
	}
	return false
}

func (e ResError) Error() string {
	return fmt.Sprintf("resp code %d with msg %s", e.Code, e.Msg)
}

// Map return map contains error info.
func (e ResError) Map() map[string]interface{} {
	m := make(map[string]interface{})
	m[RespCodeKey] = e.Code
	m[RespMsgKey] = e.Msg
	for k, v := range e.Ext {
		m[k] = v
	}
	return m
}

// NewError 创建一个新的业务响应错误
func NewError(code int, msg string) *ResError {
	return &ResError{Code: code, Msg: msg}
}

// NewErrorExt 创建一个新的业务响应错误
func NewErrorExt(code int, msg string, ext map[string]interface{}) *ResError {
	new := make(map[string]interface{})
	for k, v := range ext {
		new[k] = v
	}
	return &ResError{Code: code, Msg: msg, Ext: new}
}
