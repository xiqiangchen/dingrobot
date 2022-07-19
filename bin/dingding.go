package main

import (
	"cn.oskey/dingding"
	"io/ioutil"
	"log"
	"os"
)

type Robot struct {
	Webhook	string `json:"webhook"`
	Secret	string `json:"secret"`
}

func main() {

	urlFile := os.Args[1]

	byts, err := ioutil.ReadFile(urlFile)
	if err != nil {
		log.Fatal(err)
	}

/*
	var robot = dingding.NewRobot(
		conf.Webhook,
		conf.Secret)*/
	var robot = dingding.NewRobot(
	"https://oapi.dingtalk.com/robot/send?access_token=9ad389a96f39f8bd8669e29f7ea7c9e335d2671a41d4bfe060ca0c9aeafb90f0",
		"SECce9802d485f757363a6ad8ae2796abc4514a83c7e00dc8fe537c20b6b7a443c3")
	robot.SendTextMessage("人群包上传通知:\n" + string(byts), nil, false)
	//robot.SendMarkdownMessage("## 测试标题", "- 第一点xx  \n- 第二点oo", nil, false)
}
