package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"gpt2_project/backend/controllers"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "AccessToken", "Authentication", "X-Token", "x-token"},
		AllowCredentials: true,
	}))
	ns := beego.NewNamespace("/api",
		beego.NSBefore(controllers.CheckAuthorization),
		beego.NSNamespace("/account",
			beego.NSRouter("/login", &controllers.AccountController{}, "post:Login"),
			beego.NSRouter("/logout", &controllers.AccountController{}, "post:Logout"),
			beego.NSRouter("/user_info", &controllers.AccountController{}, "get:UserInfo"),
		),
		beego.NSNamespace("/message",
			// 1. 创建消息
			// 2. 查看消息
			// 3. 删除消息
			beego.NSRouter("/create", &controllers.MessageController{}, "post:CreateMsg"),
			beego.NSRouter("/list", &controllers.MessageController{}, "get:ListMsg"),
			beego.NSRouter("/delete", &controllers.MessageController{}, "get:DelMsg"),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/edit", &controllers.AccountController{}, "post:EditUser"),
			beego.NSRouter("/list", &controllers.AccountController{}, "get:ListUsers"),
			beego.NSRouter("/delete", &controllers.AccountController{}, "get:DelUser"),
		),
		beego.NSNamespace("/topic",
			// 1.创建选题
			// 2.调整选题状态
			// 4.删除选题
			// 5.选题列表
			beego.NSRouter("/edit", &controllers.TopicController{}, "post:EditTopic"),
			beego.NSRouter("/list", &controllers.TopicController{}, "get:ListTopic"),
			beego.NSRouter("/status", &controllers.TopicController{}, "post:EditTopic"),
			beego.NSRouter("/delete", &controllers.TopicController{}, "get:DelTopic"),
		),
		beego.NSNamespace("/article",
			//1.论文记录列表
			//2.创建论文记录
			beego.NSRouter("/edit", &controllers.ArticleController{}, "post:EditArticle"),
			beego.NSRouter("/list", &controllers.ArticleController{}, "get:ListArticle"),
		),
		beego.NSRouter("/upload", &controllers.UploadController{}, "post:Upload"),
		beego.NSNamespace("/experiment",
			beego.NSRouter("/", &controllers.ExperimentController{}, "get:Experiment"),
			beego.NSRouter("/answer", &controllers.ExperimentController{}, "post:ExperimentReply"),
		),
		beego.NSNamespace("/questionnaire",
			beego.NSRouter("/answer", &controllers.ExperimentController{}, "post:QuestionnaireReply"),
		),
		beego.NSNamespace("/chat",
			beego.NSRouter("/", &controllers.ChatController{}, "post:Chat"),
			beego.NSRouter("/event_stream", &controllers.ChatController{}, "get:ChatStream"),
		),
	)
	beego.AddNamespace(ns)

}
