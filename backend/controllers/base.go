package controllers

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/core/validation"
	"net/http"
	"reflect"
	"strings"
)

type BaseController struct {
	beego.Controller
	AppCtx context.Context
}

func (c *BaseController) echoJSON(data map[string]interface{}) {
	c.Ctx.Output.SetStatus(http.StatusOK)
	responseCode, ok := data["code"].(int)
	if !ok {
		responseCode = Success
		data["code"] = responseCode
	}
	switch responseCode {
	case ParamsErr:
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
	case TokenErr:
		c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
	}
	if _, ok := data["message"]; !ok {
		if responseMsg := RespMsg[responseCode]; len(responseMsg) > 0 {
			data["message"] = responseMsg
		}
	}

	c.Data["json"] = data
	c.ServeJSON()
}

func (c *BaseController) echoCode(code int) {
	c.Ctx.Output.SetStatus(http.StatusBadRequest)
	c.echoJSON(map[string]interface{}{RespCodeKey: code})
}

func (c *BaseController) echoParamsErr() {
	c.Ctx.Output.SetStatus(http.StatusUnprocessableEntity)
	c.echoCode(ParamsErr)
}

func (c *BaseController) parseAndValid(params interface{}) (ok bool, err error) {
	objT := reflect.TypeOf(params)
	if objT.Kind() != reflect.Ptr {
		panic("params must be pointer")
	}
	if strings.Index(c.Ctx.Request.Header.Get("Content-Type"), "application/json") != -1 {
		if err = json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
			c.echoErr(err)
			return
		}
		// json不在这里验证
		ok = true
		return
	}
	if err = c.ParseForm(params); err != nil {
		c.echoParamsErr()
		return
	}
	valid := validation.Validation{}
	if ok, err = valid.Valid(params); !ok || err != nil {
		c.echoParamsErr()
		return
	}

	return
}

func (c *BaseController) echoErr(err error) {
	if resErr, ok := err.(ResError); ok {
		c.echoJSON(resErr.Map())
	} else if resErr, ok := err.(*ResError); ok {
		c.echoJSON(resErr.Map())
	} else {
		c.echoInternalErr(err)
	}
}

func (c *BaseController) echoInternalErr(err error) {
	c.Ctx.Output.SetStatus(http.StatusInternalServerError)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	c.echoErrMsg(http.StatusInternalServerError, errMsg)
}

func (c *BaseController) echoErrMsg(code int, msg string) {
	c.Ctx.Output.SetStatus(http.StatusBadRequest)
	c.echoJSON(map[string]interface{}{RespCodeKey: code, RespMsgKey: msg})
}

type NestPreparer interface {
	NestPrepare()
}

// NestFinisher 组合类调用Finish方法需要实现的interface.
type NestFinisher interface {
	NestFinish()
}

// Prepare runs after Init before request function execution.
func (c *BaseController) Prepare() {
	c.AppCtx = context.Background()
	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

// Finish runs after request function execution.
func (c *BaseController) Finish() {
	if app, ok := c.AppController.(NestFinisher); ok {
		app.NestFinish()
	}
}

func (c *BaseController) getUid() (uid int) {
	data := c.Ctx.Input.GetData("uid")
	uid = data.(int)
	return
}
