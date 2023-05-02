package controllers

import (
	"gpt2_project/backend/models"
)

type TopicController struct {
	BaseController
}

type editTopicParams struct {
	Id         int    `orm:"pk" json:"id"`
	Tuid       int    `form:"tuid" json:"tuid"`
	Suid       int    `form:"suid" json:"suid"`               //
	Title      string `form:"title" json:"title"`             // 选题标题
	Content    string `form:"content" json:"content"`         // 选题内容
	Type       int    `form:"type" json:"type"`               // 选题类型,1=论文，2=设计
	Category   int    `form:"category" json:"category"`       // 选题二级类别,1=理论研究型，2=应用研究型，3=其他
	Source     int    `form:"source" json:"source"`           // 选题来源,1=科研，2=社会生产实践，3=其他,
	SelectType int    `form:"select_type" json:"select_type"` // 选题方式,1=学生自拟，2=老师推荐
	Status     int    `form:"status" json:"status"`           // 状态，1=确认通过，2=打回，0=待选择或确认
	Reason     string `form:"reason" json:"reason"`           // reason
}

type listTopicParams struct {
	Id      int    `form:"id"`
	TName   string `form:"t_name"`
	Title   string `form:"title"`
	Source  int    `form:"source"`
	Page    int    `form:"page"`
	PageNum int    `form:"page_num"`
}

type delTopicParams struct {
	Id int `form:"id" valid:"Required"`
}

func (c *TopicController) ListTopic() {
	params := listTopicParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	if params.PageNum == 0 {
		params.PageNum = 4
	}
	if params.Page == 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.PageNum
	ts, total, err := models.GetTopicList(offset, params.PageNum, params.Source, params.Title, params.TName, params.Id)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{"topic_list": ts, "total": total})
}

// EditTopic 修改topic属性以及状态
func (c *TopicController) EditTopic() {
	params := editTopicParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	uid := c.getUid()
	u, err := models.GetUserInfoByUid(uid)
	if err != nil {
		c.echoErr(err)
		return
	}
	switch u.Uid {
	// 1=管理员，2=老师，3=学生
	case 1, 2:
		if params.Tuid == 0 {
			params.Tuid = u.Uid
		}
		if params.Id == 0 {
			params.SelectType = 2
		}
	case 3:
		if params.Suid == 0 {
			params.Suid = u.Uid
		}
		if params.Id == 0 {
			params.SelectType = 1
		}
	}
	err = models.InsertOrUpdateTopic(models.Topic{
		Id:         params.Id,
		Tuid:       params.Tuid,
		Suid:       params.Suid,
		Title:      params.Title,
		Content:    params.Content,
		Type:       params.Type,
		Category:   params.Category,
		Source:     params.Source,
		SelectType: params.SelectType,
		Status:     params.Status,
		Reason:     params.Reason,
	})
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}

func (c *TopicController) DelTopic() {
	params := delTopicParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	err := models.DelTopic(params.Id)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}
