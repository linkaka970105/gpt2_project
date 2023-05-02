package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"gpt2_project/backend/util"
	"strings"
)

type Experiment struct {
	Id            int           `json:"id"`
	Topic         string        `json:"experiment_topic"`
	GuidePages    []GuidePage   `json:"guide_pages"`
	Questionnaire Questionnaire `json:"questionnaire"`
}

type GuidePage struct {
	Id                  int    `json:"id"`
	Content             string `json:"content"`
	NextButton          int    `json:"next_button"`
	YesOrNo             int    `json:"yes_or_no"`
	Chat                int    `json:"chat"`
	ChatTimes           int    `json:"chat_times"`
	NextPage            int    `json:"next_page"`
	YesPage             int    `json:"yes_page"`
	NoPage              int    `json:"no_page"`
	Countdown           int    `json:"countdown"`
	AnswerPage          int    `json:"answer_page"`
	AnswerTimeCountdown int    `json:"answer_time_countdown"`
}

type Questionnaire struct {
	Id        int        `json:"id"`
	Guide     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Id      int      `json:"id"`
	Content string   `json:"content"`
	Choice  []string `json:"choice"`
	Choices string   `json:"-"`
}

type Answer struct {
	Id     int    `json:"id"`
	Answer []int  `json:"answer"`
	Text   string `json:"text"`
}

type questionnaireReq struct {
	Id      int      `json:"id"`
	Answers []Answer `json:"answers"`
}

func GetExperiment(uid int) (e Experiment, err error) {
	// 1. 获取可用于展示的实验
	o := orm.NewOrm()
	sqlTpl := `select *
from gpt_project.experiment
where status = 1
order by ct desc`
	es := make([]Experiment, 0)
	_, err = o.Raw(sqlTpl).QueryRows(&es)
	if len(es) == 0 {
		return
	}
	// 2. 遍历判断用户是否已经做过该实验
	for _, ex := range es {
		fmt.Println("---- exid: ", ex.Id)
		yes, err1 := HasParticipated(uid, ex.Id, 0)
		if err1 != nil {
			err = err1
			return
		}
		if yes == 1 {
			continue
		}
		e.Id = ex.Id
		e.Topic = ex.Topic
		sqlTpl = `select *
			from gpt_project.guide_page
			where exid = ?`
		_, err = o.Raw(sqlTpl, ex.Id).QueryRows(&e.GuidePages)
		if err != nil {
			return
		}
		sqlTpl = `select *
			from gpt_project.questionnaire
			where exid = ? limit 1`

		err = o.Raw(sqlTpl, ex.Id).QueryRow(&e.Questionnaire)
		if err != nil {
			if err == orm.ErrNoRows {
				err = nil
			} else {
				return
			}
		}
		sqlTpl = `select *
		from gpt_project.question
		where qnid = ?
		order by id`

		_, err = o.Raw(sqlTpl, e.Questionnaire.Id).QueryRows(&e.Questionnaire.Questions)
		if err != nil {
			return
		}
		for i := range e.Questionnaire.Questions {
			e.Questionnaire.Questions[i].Choice = strings.Split(e.Questionnaire.Questions[i].Choices, ";")
		}
		return
	}
	return
}

func ParticipatedRecord(uid, exid int, content string, access int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("insert into gpt_project.experiment_reply (uid, exid, content, access) values (?, ?, ?, ?)",
		uid, exid, content, access).Exec()
	return
}

func HasParticipated(uid, exid int, access int) (yes int, err error) {
	o := orm.NewOrm()
	sqlTpl := `select 1
from gpt_project.experiment_reply
where uid = ?
  and exid = ?
  and access = ?`
	err = o.Raw(sqlTpl, uid, exid, access).QueryRow(&yes)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		} else {
			return
		}
	}
	return
}

func QuestionReply(uid, qnid int, answers []Answer) (err error) {
	o := orm.NewOrm()
	execSql := `insert into gpt_project.question_reply (uid, qnid, qid, reply_text, choices) values %s`
	placeholder := "(?, ?, ?, ?, ?)"
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)
	for i := 0; i < len(answers); i++ {
		placeholders = append(placeholders, placeholder)
		args = append(args, uid, qnid, answers[i].Id, answers[i].Text, util.IntArrayToBinary(util.Index2Int(answers[i].Answer)))
	}
	execSql = fmt.Sprintf(execSql, strings.Join(placeholders, ","))
	_, err = o.Raw(execSql, args...).Exec()
	return
}

func ChatRecord(uid, exid int, message, reply string) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("insert into gpt_project.experiment_chat_history (uid, exid, query, gpt_reply) values (?, ?, ?, ?)",
		uid, exid, message, reply).Exec()
	return
}