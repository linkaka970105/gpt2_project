package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Message struct {
	Id          int       `orm:"pk" json:"id"`
	Category    int       `json:"category"`
	MsgFrom     int       `json:"msg_from"`
	MsgTo       int       `json:"msg_to"`
	MsgTitle    string    `json:"msg_title"`
	MsgContent  string    `json:"msg_content"`
	Ct          time.Time `json:"-"`
	Ut          time.Time `json:"-"`
	CtStr       string    `json:"ct"`
	MsgToName   string    `json:"msg_to_name"`
	MsgFromName string    `json:"msg_from_name"`
}

func GetMsgList(uid, offset, pageNum int, title string, category int, id int) (msgList []Message, total int, err error) {
	o := orm.NewOrm()
	sqlTpl := `select *
from gpt_project.notice
%s
order by ct desc
limit ?,?`
	countSql := `select count(1)
from gpt_project.notice
%s`
	where := []string{"msg_to in (0, ?) or msg_from = ?"}
	args := []interface{}{uid, uid}
	if title != "" {
		where = append(where, fmt.Sprintf("msg_title like '%%%s%%'", title))
	}
	if category > 0 {
		where = append(where, "category = ?")
		args = append(args, category)
	}
	if id > 0 {
		where = append(where, "id = ?")
		args = append(args, id)
	}
	countArgs := args
	args = append(args, offset, pageNum)
	sqlTpl = fmt.Sprintf(sqlTpl, "where "+strings.Join(where, " and "))
	_, err = o.Raw(sqlTpl, args...).QueryRows(&msgList)
	fmt.Printf("sqlTpl: %+v\n", sqlTpl)
	fmt.Printf("args: %+v\n", args)
	if err != nil {
		return
	}
	for i := range msgList {
		msgList[i].CtStr = msgList[i].Ct.Format("2006/01/02 15:04")
		if msgList[i].MsgFrom != 0 {
			msgList[i].MsgFromName = ""
		}
		if msgList[i].MsgTo != 0 {
			msgList[i].MsgToName = ""
		}
	}
	countSql = fmt.Sprintf(countSql, "where "+strings.Join(where, " and "))
	err = o.Raw(countSql, countArgs...).QueryRow(&total)
	if err != nil {
		return
	}
	return
}

func InsertMsg(msg Message) (err error) {
	o := orm.NewOrm()
	sqlTpl := `insert into gpt_project.notice (id,category, msg_from, msg_to, msg_title, msg_content)
values (?, ?, ?, ?, ?, ?)
on duplicate key update msg_title=if(values(msg_title) != '', values(msg_title), msg_title),
                        msg_content=if(values(msg_content) != '', values(msg_content), msg_content)`
	_, err = o.Raw(sqlTpl, msg.Id, msg.Category, msg.MsgFrom, msg.MsgTo, msg.MsgTitle, msg.MsgContent).Exec()
	return
}

func DelMsg(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("delete from gpt_project.notice where id = ?", id).Exec()
	return
}
