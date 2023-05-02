package controllers

import "log"

type UploadController struct {
	BaseController
}

func (c *UploadController) Upload() {
	f, h, err := c.GetFile("file")
	if err != nil {
		log.Fatal("getfile err ", err)
	}
	defer f.Close()
	err = c.SaveToFile("file", "static/upload/"+h.Filename) // 保存位置在 static/upload, 没有文件夹要先创建
	if err != nil {
		c.echoErr(err)
		return
	}
	c.echoJSON(map[string]interface{}{"url": "http://192.168.0.110:8080/static/upload/" + h.Filename})
}
