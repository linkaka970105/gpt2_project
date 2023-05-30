package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	"gpt2_project/backend/util"
	"strconv"
	"strings"
	"time"
)

var defaultTokenExpireSec = 24 * 3600 // 默认token保存时间
var PwdNotMatchErr = errors.New("password not match")

type Users struct {
	Uid      int    `orm:"pk" json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type BasicInfo struct {
	LoginTitle      string `json:"login_title"`
	IndexBottomTips string `json:"index_bottom_tips"`
}

type ListUserQuery struct {
	School    int    `form:"school"`           // 学院
	Class     int    `form:"class"`            // 班级
	JobTitle  int    `form:"job_title"`        // 职位
	Name      string `form:"name"`             // 人员名称
	Subject   int    `form:"subject"`          // 专业
	Sex       int    `form:"sex"`              // 性别
	StudentNo string `form:"student_no"`       // 学号
	Type      int    `form:"type" json:"type"` // 账号类型
	Page      int    `form:"page"`
	PageNum   int    `form:"page_num"`
}

func GetUserInfoByEmail(email, pwd string) (u Users, err error) {
	o := orm.NewOrm()
	sqlTpl := `select *
			from gpt_project.users
			where email = ?`
	err = o.Raw(sqlTpl, email).QueryRow(&u)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		} else {
			return
		}
	}
	if u.Uid == 0 {
		execSql := `insert into gpt_project.users (email, password)
values (?, ?)`
		_, err = o.Raw(execSql, email, pwd).Exec()
		if err != nil {
			return
		}
		err = o.Raw(sqlTpl, email).QueryRow(&u)
		if err != nil {
			return
		}
	} else {
		// 判断密码是否匹配
		if u.Password != pwd {
			err = PwdNotMatchErr
			return
		}
	}
	return
}

func GetUserInfoByUid(uid int) (u Users, err error) {
	o := orm.NewOrm()
	sqlTpl := `select *
			from gpt_project.users
			where uid = ?`
	err = o.Raw(sqlTpl, uid).QueryRow(&u)
	if err != nil {
		return
	}
	if u.Uid == 0 {
		err = errors.New("user not found")
		return
	}
	return
}

func GetUserInfoByNameOrSN(name string) (u Users, err error) {
	o := orm.NewOrm()
	sqlTpl := `select *
			from gpt_project.users
			where student_no = ? or name = ? limit 1`
	err = o.Raw(sqlTpl, name, name).QueryRow(&u)
	if err != nil {
		return
	}
	if u.Uid == 0 {
		err = errors.New("user not found")
		return
	}
	return
}

func InsertOrUpdateUser(u Users) (err error) {
	sqlTpl := `insert into gpt_project.users (email, password)
values (?, ?)`
	args := make([]interface{}, 0)
	args = append(args, u.Uid, u.Password)
	o := orm.NewOrm()
	_, err = o.Raw(sqlTpl, args...).Exec()
	return
}

func ListUsers(q ListUserQuery) (users []Users, total int, err error) {
	o := orm.NewOrm()
	sqlTpl := `select * from gpt_project.users %s order by uid desc
limit ?,?`
	countSql := `select count(1) from gpt_project.users %s`
	offset := (q.Page - 1) * q.PageNum
	args := make([]interface{}, 0)
	where := make([]string, 0)
	if q.School != 0 {
		where = append(where, "school = ?")
		args = append(args, q.School)
	}
	if q.Class != 0 {
		where = append(where, "class = ?")
		args = append(args, q.Class)
	}
	if q.JobTitle != 0 {
		where = append(where, "job_title = ?")
		args = append(args, q.JobTitle)
	}
	if q.Subject != 0 {
		where = append(where, "subject = ?")
		args = append(args, q.Subject)
	}
	if q.Sex != 0 {
		where = append(where, "sex = ?")
		args = append(args, q.Sex)
	}
	if q.Name != "" {
		where = append(where, fmt.Sprintf("name like '%%%s%%'", q.Name))
	}
	if q.StudentNo != "" {
		where = append(where, fmt.Sprintf("student_no like '%%%s%%'", q.StudentNo))
	}
	if q.Type != 0 {
		where = append(where, "type = ?")
		args = append(args, q.Type)
	}
	countArgs := args
	args = append(args, offset, q.PageNum)
	whereField := ""
	if len(where) > 0 {
		whereField = " where " + strings.Join(where, " and ")
	}
	_, err = o.Raw(fmt.Sprintf(sqlTpl, whereField), args...).QueryRows(&users)
	if err != nil {
		return
	}
	err = o.Raw(fmt.Sprintf(countSql, whereField), countArgs...).QueryRow(&total)
	return
}

func updateToken(uid int, token string) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update gpt_project.users set token = ? where uid = ?", token, uid).Exec()
	return
}

func DelUser(uid int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("delete from gpt_project.users where uid = ?", uid).Exec()
	return
}

// NewToken 分配新的token
func NewToken(uid int) (token string, err error) {
	token = newToken()
	oldToken, err := getTokenByUID(uid)
	if err != nil {
		if err != redis.Nil {
			return
		}
		err = nil
	} else {
		if err = RedisCli().Del(makeTokenUIDKey(oldToken)).Err(); err != nil {
			return
		}
		if err = RedisCli().Del(oldToken).Err(); err != nil {
			// 删除历史聊天记录
			return
		}
	}
	// 保存新token
	if err = saveToken(RedisCli(), uid, token, defaultTokenExpireSec); err != nil {
		return
	}
	return
}

// GetUIDByToken 根据token获取对应的用户
func GetUIDByToken(token string) (uid int, err error) {
	key := makeTokenUIDKey(token)
	rt, err := RedisCli().Get(key).Result()
	if err != nil {
		return
	}
	return strconv.Atoi(rt)
}

// DelTokenByUID 删除uid对应的token
func DelTokenByUID(uid int) (err error) {
	// 如果旧的token存在,则删除
	token, err := getTokenByUID(uid)
	if err != nil {
		if err != redis.Nil {
			return
		}
		return nil
	}
	return DelToken(token, uid)
}

// DelToken 删除token
func DelToken(token string, uid int) (err error) {
	if err = RedisCli().Del(makeUIDTokenKey(uid)).Err(); err != nil {
		return
	}
	if err = RedisCli().Del(makeTokenUIDKey(token)).Err(); err != nil {
		return
	}
	return
}

func GetBasicInfo() (b BasicInfo, err error) {
	o := orm.NewOrm()
	sqlTpl := `select login_title, index_bottom_tips
				from gpt_project.basic_info
				order by id desc
				limit 1`
	err = o.Raw(sqlTpl).QueryRow(&b)
	if err != nil {
		return
	}
	return
}

func getTokenByUID(uid int) (token string, err error) {
	key := makeUIDTokenKey(uid)
	return RedisCli().Get(key).Result()
}

func newToken() string {
	newid := strings.ToUpper(util.UUID().String())
	return strings.ReplaceAll(newid, "-", "")
}

func saveToken(redisCli *redis.Client, uid int, token string, timeoutSec int) (err error) {
	expire := time.Duration(timeoutSec) * time.Second
	value := strconv.Itoa(uid)
	// token->uid的映射关系
	if err = redisCli.Set(makeTokenUIDKey(token), value, expire).Err(); err != nil {
		return
	}
	// uid->token的映射关系
	if err = redisCli.Set(makeUIDTokenKey(uid), token, expire).Err(); err != nil {
		return
	}
	err = updateToken(uid, token)
	return
}

func makeTokenUIDKey(token string) string {
	return fmt.Sprint("user_token:", token)
}

func makeUIDTokenKey(uid int) string {
	return fmt.Sprint("uid_token:", uid)
}
