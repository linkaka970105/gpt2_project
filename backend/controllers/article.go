package controllers

import (
	"gpt2_project/backend/models"
)

type ArticleController struct {
	BaseController
}

type listArticleParams struct {
	Id             int    `form:"id"`
	Title          string `form:"title"`
	ProgressStatus int    `form:"progress_status"` // 论文过程类型：1=开题，2=中期，3=初稿，4=定稿
	Page           int    `form:"page"`
	PageNum        int    `form:"page_num"`

	IsDetail int `form:"is_detail"`
}

type editArticleParams struct {
	Id             int    `form:"id" json:"id"`
	Tuid           int    `form:"tuid" json:"tuid"`
	Suid           int    `form:"suid" json:"suid"`
	Atid           int    `form:"atid" json:"atid"`
	ProgressStatus int    `form:"progress_status" json:"progress_status"`
	FileUrl        string `form:"file_url" json:"file_url"`
	AppendixUrl    string `form:"appendix_url" json:"appendix_url"`
	Reason         string `form:"reason" json:"reason"`
	Status         int    `form:"status" json:"status"`
}

type delArticleParams struct {
	Id int `form:"id" valid:"Required"`
}

func (c *ArticleController) ListArticle() {
	params := listArticleParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	uid := c.getUid()
	u, err := models.GetUserInfoByUid(uid)
	if err != nil {
		c.echoErr(err)
		return
	}
	var suid int
	if u.Uid == 3 {
		suid = uid
	}
	if params.PageNum == 0 {
		params.PageNum = 4
	}
	if params.Page == 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.PageNum
	as, total, err := models.GetArticleList(offset, params.PageNum, params.Title, params.Id, params.ProgressStatus, suid)
	if err != nil {
		c.echoErr(err)
		return
	}
	if len(as) == 0 && params.IsDetail == 1 {
		ts, _ := models.QueryTopicByUid(u.Uid)
		title := ""
		if len(ts) > 0 {
			title = ts[0].Title
		}
		as = append(as, models.Article{
			SName:     u.Email,
			Title:     title,
		})
		total = 1
	}
	c.echoJSON(map[string]interface{}{"article_list": as, "total": total})
}

// EditArticle 修改article属性以及状态
func (c *ArticleController) EditArticle() {
	params := editArticleParams{}
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
	case 3:
		if params.Suid == 0 {
			params.Suid = u.Uid
		}
		if params.Id == 0 {
			var ts []models.Topic
			ts, err = models.QueryTopicByUid(u.Uid)
			if err != nil {
				c.echoErr(err)
				return
			}
			if len(ts) > 0 {
				params.Atid = ts[0].Id
			}
		}
	}
	err = models.InsertOrUpdateArticle(models.Article{
		Id:             params.Id,
		Tuid:           params.Tuid,
		Suid:           params.Suid,
		Atid:           params.Atid,
		ProgressStatus: params.ProgressStatus,
		FileUrl:        params.FileUrl,
		AppendixUrl:    params.AppendixUrl,
		Reason:         params.Reason,
		Status:         params.Status,
	})
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}
