package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
)

var rt *redis.Client
var gptClient *openai.Client

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

	config := openai.DefaultConfig("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")
	proxyUrl, err := url.Parse("http://127.0.0.1:12333")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

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
