package controllers

import (
	"encoding/json"
	"gpt2_project/backend/models"
)

type ExperimentController struct {
	BaseController
}

type experimentReplyParams struct {
	ID     int    `form:"id" json:"id" valid:"Required"`
	Answer string `form:"answer" json:"answer"`
}

type Answer struct {
	ID     int    `json:"id"`
	Answer []int  `json:"answer"`
	Text   string `json:"text"`
}

type questionnaireReq struct {
	ID      int             `json:"id"`
	Answers []models.Answer `json:"answers"`
}

func (c *ExperimentController) Experiment() {
	uid := c.getUid()
	// 获取实验
	ex, err := models.GetExperiment(uid)
	if err != nil {
		c.echoErr(err)
		return
	}
	if ex.Id > 0 {
		yes, err1 := models.HasParticipated(uid, ex.Id, 1)
		if err1 != nil {
			c.echoErr(err1)
			return
		}
		if yes == 0 {
			err = models.ParticipatedRecord(uid, ex.Id, "", 1)
			if err != nil {
				c.echoErr(err)
				return
			}
		}
	}
	c.echoJSON(map[string]interface{}{"experiment": ex})
}

func (c *ExperimentController) ExperimentReply() {
	uid := c.getUid()
	params := experimentReplyParams{}
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
	err := models.ParticipatedRecord(uid, params.ID, params.Answer, 0)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}

func (c *ExperimentController) QuestionnaireReply() {
	uid := c.getUid()
	params := questionnaireReq{}
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
	err := models.QuestionReply(uid, params.ID, params.Answers)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}
