package main

import (
	"github.com/astaxie/beego"
	_ "gpt2_project/backend/models"
	_ "gpt2_project/backend/routers"
)

func main() {
	beego.Run(":8089")
}
