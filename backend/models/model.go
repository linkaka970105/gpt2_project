package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
)

var rt *redis.Client
var gptClient *openai.Client
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterDataBase("default", "mysql", "root:3702200733ljj@tcp(127.0.0.1:3306)/gpt_project?charset=utf8")
	//orm.RegisterDataBase("default", "mysql", "root:3702200733Ljj@@tcp(127.0.0.1:3306)/gpt_project?charset=utf8")
	//orm.RegisterModel(new(Users))
	//orm.RegisterModel()
	//orm.RunSyncdb("default", false, true)

	redisCli, err := NewRedis("127.0.0.1:6379", 0)
	if err != nil {
		panic(err)
	}
	rt = redisCli

	config := openai.DefaultConfig("fk193074-slhtDybpXZkkvJronLZ0sa0oX4azh7aD")
	config.BaseURL = "https://openai.api2d.net/v1"
	//proxyUrl, err := url.Parse("http://127.0.0.1:12333")
	//if err != nil {
	//	panic(err)
	//}
	//transport := &http.Transport{
	//	Proxy: http.ProxyURL(proxyUrl),
	//}
	//config.HTTPClient = &http.Client{
	//	Transport: transport,
	//}

	gptClient = openai.NewClientWithConfig(config)
}

// NewRedis new redis pool
func NewRedis(addr string, db int) (rt *redis.Client, err error) {
	cli := redis.NewClient(&redis.Options{
		Addr:         addr,
		MinIdleConns: 5,
		DB:           db,
	})
	if err = cli.Ping().Err(); err != nil {
		return
	}
	rt = cli
	return
}

func RedisCli() *redis.Client {
	return rt
}

func GptCli() *openai.Client {
	return gptClient
}
