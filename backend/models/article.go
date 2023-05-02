package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Article struct {
	Id             int    `orm:"pk" json:"id"`
	Tuid           int    `json:"tuid"`
	Suid           int    `json:"suid"`
	Atid           int    `json:"atid"`
	ProgressStatus int    `json:"progress_status"`
	FileUrl        string `json:"file_url"`
	AppendixUrl    string `json:"appendix_url"`
	Reason         string `json:"reason"`
	Status         int    `json:"status"`

	ShowStatus int       `json:"show_status"`
	School     int       `json:"school"`
	Title      string    `json:"title"`
	TName      string    `json:"t_name"`
	SName      string    `json:"s_name"`
	StudentNo  string    `json:"student_no"`
	Subject    int       `json:"subject"`
	Ct         time.Time `json:"-"`
	CtStr      string    `json:"ct"`
	Ut         time.Time `json:"-"`
}

func GetArticleList(offset, pageNum int, title string, id, ProgressStatus, suid int) (as []Article, total int, err error) {
	o := orm.NewOrm()
	sqlTpl := `select a.*, su.name as s_name,tu.name as t_name,at.title,su.school,su.subject,su.student_no
from gpt_project.article a
left join gpt_project.users su on a.suid = su.uid
left join gpt_project.users tu on a.tuid = tu.uid
left join gpt_project.article_topic at on a.at_id = at.id
%s
order by a.ct desc
limit ?,?`
	countSql := `select count(1)
from gpt_project.article a
left join gpt_project.users su on a.suid = su.uid
left join gpt_project.users tu on a.tuid = tu.uid
left join gpt_project.article_topic at on a.at_id = at.id
%s`
	var where []string
	var args []interface{}
	if title != "" {
		where = append(where, fmt.Sprintf("at.title like '%%%s%%'", title))
	}
	if id > 0 {
		where = append(where, "a.id = ?")
		args = append(args, id)
	}
	if ProgressStatus > 0 {
		where = append(where, "a.progress_status = ?")
		args = append(args, ProgressStatus)
	}
	if suid > 0 {
		where = append(where, "su.uid = ?")
		args = append(args, suid)
	}
	countArgs := args
	args = append(args, offset, pageNum)
	whereFields := ""
	if len(where) > 0 {
		whereFields = "where " + strings.Join(where, " and ")
	}
	sqlTpl = fmt.Sprintf(sqlTpl, whereFields)
	fmt.Println("sqlTpl: ", sqlTpl)
	fmt.Println("args: ", args)
	_, err = o.Raw(sqlTpl, args...).QueryRows(&as)
	if err != nil {
		return
	}
	for i := range as {
		as[i].CtStr = as[i].Ct.Format("2006/01/02 15:04")
		if as[i].Status == 0 {
			as[i].ShowStatus = 1 // 待通过
		} else if as[i].Status == 1 {
			as[i].ShowStatus = 2 // 通过
		} else if as[i].Status == 2 {
			as[i].ShowStatus = 3 // 打回
		}
	}
	countSql = fmt.Sprintf(countSql, whereFields)
	err = o.Raw(countSql, countArgs...).QueryRow(&total)
	if err != nil {
		return
	}
	return
}

func InsertOrUpdateArticle(t Article) (err error) {
	sqlTpl := `insert into gpt_project.article (id, tuid, suid, at_id, progress_status, file_url, appendix_url, reason, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
on duplicate key update file_url=if(values(file_url) != '', values(file_url), file_url),
                        appendix_url=if(values(appendix_url) != '', values(appendix_url), appendix_url),
                        reason=if(values(reason) != '', values(reason), reason),
                        progress_status=if(values(progress_status) != 0, values(progress_status), progress_status),
                        status=if(values(status) != 0, values(status), status)`
	args := make([]interface{}, 0)
	args = append(args, t.Id, t.Tuid, t.Suid, t.Atid, t.ProgressStatus, t.FileUrl, t.AppendixUrl, t.Reason, t.Status)
	o := orm.NewOrm()
	_, err = o.Raw(sqlTpl, args...).Exec()
	return
}

func DelArticle(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("delete from gpt_project.article_topic where id = ?", id).Exec()
	return
}
