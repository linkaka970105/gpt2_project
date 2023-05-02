package controllers

import (
	"gpt2_project/backend/models"
)

type MessageController struct {
	BaseController
}

type createMsgParams struct {
	Id         int    `form:"id" json:"id"`                                    // notice id
	ToUid      int    `form:"to_uid" json:"to_uid"`                            // 收件人uid，所有人则为0
	Category   int    `form:"category" valid:"Required" json:"category"`       // 消息类别，1=个人消息，2=公告'
	MsgToName  string `form:"msg_to_name" json:"msg_to_name"`                  // 姓名或学号
	MsgTitle   string `form:"msg_title" valid:"Required" json:"msg_title"`     // 消息标题
	MsgContent string `form:"msg_content" valid:"Required" json:"msg_content"` // 消息内容
}

type listMsgParams struct {
	Id       int    `form:"id"`
	Title    string `form:"title"`
	Category int    `form:"category"`
	Page     int    `form:"page"`
	PageNum  int    `form:"page_num"`
}

type delMsgParams struct {
	Id int `form:"id" valid:"Required"`
}

func (c *MessageController) CreateMsg() {
	params := createMsgParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	uid := c.getUid()
	if params.MsgToName != "" {
		u, _ := models.GetUserInfoByNameOrSN(params.MsgToName)
		params.ToUid = u.Uid
	}
	if params.ToUid == 0 {
		params.Category = 2
	} else {
		params.Category = 1 // 个人消息
	}
	err := models.InsertMsg(models.Message{
		Id:         params.Id,
		Category:   params.Category,
		MsgFrom:    uid,
		MsgTo:      params.ToUid,
		MsgTitle:   params.MsgTitle,
		MsgContent: params.MsgContent,
	})
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}

func (c *MessageController) ListMsg() {
	params := listMsgParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	if params.PageNum == 0 {
		params.PageNum = 4
	}
	if params.Page == 0 {
		params.Page = 1
	}
	uid := c.getUid()
	offset := (params.Page - 1) * params.PageNum
	msgList, total, err := models.GetMsgList(uid, offset, params.PageNum, params.Title, params.Category, params.Id)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{"msg_list": msgList, "total": total})
}

func (c *MessageController) DelMsg() {
	params := delMsgParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	err := models.DelMsg(params.Id)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}
