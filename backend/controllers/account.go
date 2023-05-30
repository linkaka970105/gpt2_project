package controllers

import (
	"encoding/json"
	beegoctx "github.com/astaxie/beego/context"
	"gpt2_project/backend/models"
	"net/http"
)

var rolesMap = map[int]string{
	1: "admin",
	2: "teacher",
	3: "student",
}

type AccountController struct {
	BaseController
}

type loginParams struct {
	Email    string `form:"email" json:"email" valid:"Required"`
	Password string `form:"password" json:"password" valid:"Required"`
}

type editUserParams struct {
	Uid       int    `form:"uid" json:"uid"` // uid,修改时需传2
	Name      string `form:"name" json:"name"`
	StudentNo string `form:"student_no" json:"student_no"` // 学号或者工号
	Password  string `form:"password" json:"password"`     // 默认密码

	Sex      int `form:"sex" json:"sex"`             // 性别，1=男，2=女
	Type     int `form:"type" json:"type"`           // 用户类型，1=管理员，2=老师，3=学生
	JobTitle int `form:"job_title" json:"job_title"` // 职务,1=学院导师，2=学院主任，3=学院院长
	School   int `form:"school" json:"school"`       // 学院，1=数科院
	Subject  int `form:"subject" json:"subject"`     // 专业，1=金融分析
	Class    int `form:"class" json:"class"`         // class id

	IsChangePwd int `form:"is_change_pwd" json:"is_change_pwd"` // is_change_pwd 是否改密码
}

type listUserParams struct {
	models.ListUserQuery
}

type delUserParams struct {
	Uid int `form:"uid" valid:"Required"`
}

type loginInfo struct {
	Token string `json:"token"`
}

func (c *AccountController) Login() {
	params := loginParams{}
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
	// 获取账号数据，验证password
	u, err := models.GetUserInfoByEmail(params.Email, params.Password)
	if err != nil {
		if err == models.PwdNotMatchErr {
			c.echoErrMsg(PwdNotMatch, RespMsg[PwdNotMatch])
		} else {
			c.echoErr(err)
		}
		return
	}
	// 生成新token，更新db，redis
	token, err := models.NewToken(u.Uid)
	if err != nil {
		c.echoErr(err)
		return
	}

	info := loginInfo{
		Token: token,
	}
	c.echoJSON(map[string]interface{}{"info": info})
}

func (c *AccountController) Logout() {
	// 删除token
	uid := c.getUid()
	err := models.DelTokenByUID(uid)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}

func (c *AccountController) UserInfo() {
	// 删除token
	uid := c.getUid()
	u, err := models.GetUserInfoByUid(uid)
	if err != nil {
		c.echoErr(err)
		return
	}
	var isSelected int
	if u.Uid == 3 {
		// 学生只能选题一次
		var ts []models.Topic
		ts, err = models.QueryTopicByUid(u.Uid)
		if err != nil {
			c.echoErr(err)
			return
		}
		if len(ts) > 0 {
			isSelected = 1
		}
	}
	info := struct {
		Roles         []string `json:"roles"`
		Introduction  string   `json:"introduction"`
		Avatar        string   `json:"avatar"`
		Name          string   `json:"name"`
		PhoneNum      string   `json:"phone_num"`
		Email         string   `json:"email"`
		TopicSelected int      `json:"topic_selected"`
	}{
		Roles:         []string{rolesMap[u.Uid]},
		Introduction:  "hello, i am from xxxx",
		Avatar:        "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Name:          u.Email,
		PhoneNum:      "",
		Email:         u.Email,
		TopicSelected: isSelected,
	}
	c.echoJSON(map[string]interface{}{"data": info})
}

func (c *AccountController) EditUser() {
	params := editUserParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	if params.IsChangePwd == 1 && params.Password != "" {
		params.Uid = c.getUid()
	}
	err := models.InsertOrUpdateUser(models.Users{
		Uid:      params.Uid,
		Password: params.Password,
	})
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}

func (c *AccountController) ListUsers() {
	params := listUserParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	if params.PageNum == 0 {
		params.PageNum = 4
	}
	if params.Page == 0 {
		params.Page = 1
	}
	users, total, err := models.ListUsers(params.ListUserQuery)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{"users": users, "total": total})
}

func (c *AccountController) DelUser() {
	params := delUserParams{}
	if ok, err := c.parseAndValid(&params); !ok || err != nil {
		return
	}
	err := models.DelUser(params.Uid)
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{})
}

func (c *AccountController) BasicInfo() {
	basicInfo, err := models.GetBasicInfo()
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{
		"login_title":       basicInfo.LoginTitle,
		"index_bottom_tips": basicInfo.IndexBottomTips,
	})
}

func CheckAuthorization(ctx *beegoctx.Context) {
	routerPath := ctx.Request.URL.Path
	if isExcluedRouterPath(routerPath) {
		return
	}
	token := ctx.Request.Header.Get("X-Token")
	if token == "" {
		token = ctx.Request.URL.Query().Get("token")
	}
	// check token
	uid, _ := models.GetUIDByToken(token)
	if uid == 0 {
		ctx.Abort(http.StatusUnauthorized, "401")
		return
	}
	ctx.Input.SetData("uid", uid)
}

var excludeRouterPathes = []string{
	"/api/account/login",
	"/api/upload",
	"/api/basic_info",
}

func isExcluedRouterPath(path string) bool {
	for _, routerPath := range excludeRouterPathes {
		if routerPath == path {
			return true
		}
	}
	return false
}
