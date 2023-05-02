package main

import (
	_ "gpt2_project/backend/models"
	_ "gpt2_project/backend/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}