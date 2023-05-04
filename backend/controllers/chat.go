package controllers

import (
	"encoding/json"
	"gpt2_project/backend/models"
)

type ChatController struct {
	BaseController
}

type chatParams struct {
	Id      int    `form:"id" json:"id"`
	Message string `form:"message" json:"message"`
}

func (c *ChatController) Chat() {
	uid := c.getUid()
	var params chatParams
	if c.Ctx.Request.Header.Get("Content-Type") == "application/json" {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &params); err != nil {
			c.echoErr(err)
			return
		}
	} else {
		if ok, err := c.parseAndValid(&params); !ok || err != nil {
			return
		}
	}
	//reply := "gpt的回复内容：" + params.Message
	reply, err := models.GetGptResp(uid, params.Message)
	if err != nil {
		c.echoErr(err)
		return
	}
	if params.Id > 0 {
		err := models.ChatRecord(uid, params.Id, params.Message, reply)
		if err != nil {
			c.echoErr(err)
			return
		}
	}
	c.echoJSON(map[string]interface{}{"reply": reply})
}
