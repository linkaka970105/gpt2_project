package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Topic struct {
	Id         int       `orm:"pk" json:"id"`
	Tuid       int       `json:"-"` // 选题老师uid
	Suid       int       `json:"-"` // 学生uid
	TName      string    `json:"t_name"`
	SName      string    `json:"s_name"`
	School     int       `json:"school"`
	Title      string    `json:"title"`       // 选题标题
	Content    string    `json:"content"`     // 选题内容
	Type       int       `json:"type"`        // 选题类型,1=论文，2=设计
	Category   int       `json:"category"`    // 选题二级类别,1=理论研究型，2=应用研究型，3=其他
	Source     int       `json:"source"`      // 选题来源,1=科研，2=社会生产实践，3=其他,
	SelectType int       `json:"select_type"` // 选题方式,1=学生自拟，2=老师推荐
	Status     int       `json:"status"`      // 状态，1=通过，0未通过,2打回
	CtStr      string    `json:"ct"`
	Publisher  string    `json:"publisher"`
	ShowStatus int       `json:"show_status"` // 1=待选取,2=待通过,3=通过,4=打回
	Reason     string    `json:"reason"`
	Ct         time.Time `json:"-"`
	Ut         time.Time `json:"-"`
}

func GetTopicList(offset, pageNum, source int, title, tName string, id int) (ts []Topic, total int, err error) {
	o := orm.NewOrm()
	sqlTpl := `select at.*,ut.name as t_name, us.name as s_name, ut.school
from gpt_project.article_topic at
left join gpt_project.users ut on at.tuid = ut.uid
left join gpt_project.users us on at.suid = us.uid
%s 
order by ct desc
limit ?,?`
	countSql := `select count(1)
from gpt_project.article_topic at
left join gpt_project.users ut on at.tuid = ut.uid
left join gpt_project.users us on at.suid = us.uid
%s`
	var where []string
	var args []interface{}
	if title != "" {
		where = append(where, fmt.Sprintf("at.title like '%%%s%%'", title))
	}
	if source > 0 {
		where = append(where, "at.source = ?")
		args = append(args, source)
	}
	if id > 0 {
		where = append(where, "at.id = ?")
		args = append(args, id)
	}
	if tName != "" {
		where = append(where, "((select_type = 1 and us.name = ?) or (select_type = 2 and ut.name = ?))")
		args = append(args, tName, tName)
	}
	countArgs := args
	args = append(args, offset, pageNum)
	whereFields := ""
	if len(where) > 0 {
		whereFields = "where " + strings.Join(where, " and ")
	}
	sqlTpl = fmt.Sprintf(sqlTpl, whereFields)
	_, err = o.Raw(sqlTpl, args...).QueryRows(&ts)
	if err != nil {
		return
	}
	for i := range ts {
		ts[i].CtStr = ts[i].Ct.Format("2006/01/02 15:04")
		if ts[i].SelectType == 1 {
			ts[i].Publisher = ts[i].SName
		} else {
			ts[i].Publisher = ts[i].TName
		}
		if ts[i].Suid == 0 {
			ts[i].ShowStatus = 1 // 带选取
		} else if ts[i].Status == 0 {
			ts[i].ShowStatus = 2 // 待通过
		} else if ts[i].Status == 1 {
			ts[i].ShowStatus = 3 // 通过
		} else if ts[i].Status == 2 {
			ts[i].ShowStatus = 4 // 打回
		}
	}
	countSql = fmt.Sprintf(countSql, whereFields)
	err = o.Raw(countSql, countArgs...).QueryRow(&total)
	if err != nil {
		return
	}
	return
}

func InsertOrUpdateTopic(t Topic) (err error) {
	sqlTpl := `insert into gpt_project.article_topic (id, tuid, suid, title, content,
                                           type, category, source, select_type, status, reason)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
on duplicate key update title=if(values(title) != '', values(title), title),
                        content=if(values(content) != '', values(content), content),
                        reason=if(values(reason) != '', values(reason), reason),
                        type=if(values(type) != 0, values(type), type),
                        category=if(values(category) != 0, values(category), category),
                        source=if(values(source) != 0, values(source), source),
                        select_type=if(values(select_type) != 0, values(select_type), select_type),
                        status=if(values(status) != 0, values(status), status),
                        suid=if(values(suid) != 0 and status != 1, values(suid), suid)`
	args := make([]interface{}, 0)
	args = append(args, t.Id, t.Tuid, t.Suid, t.Title, t.Content,
		t.Type, t.Category, t.Source, t.SelectType, t.Status, t.Reason)
	o := orm.NewOrm()
	_, err = o.Raw(sqlTpl, args...).Exec()
	return
}

func DelTopic(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("delete from gpt_project.article_topic where id = ?", id).Exec()
	return
}

func QueryTopicByUid(suid int) (ts []Topic, err error) {
	o := orm.NewOrm()
	sqlTpl := `select id,title
from gpt_project.article_topic
where suid = ? and status = 1
order by ct desc
limit 1`
	_, err = o.Raw(sqlTpl, suid).QueryRows(&ts)
	return
}
