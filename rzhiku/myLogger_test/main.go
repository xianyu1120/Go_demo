package main

import "github.com/demo/rzhiku/mylogger"

var log mylogger.Logger//声明一个全局的接口变量

//测试我们写的日志库

func main() {
	log=mylogger.NewConsoleLog("Info")//终端日志记录
	log := mylogger.NewFileLogger("Info", "./", "test.log", 100*1024*1024)//
	
	for {
		log.Debug("这是一个debug日志")
		log.Info("这是一个info日志")
		log.Warning("这是一个warning日志")
		id := 10010
		name := "理想"
		log.Error("这是一个error日志,id:%d,name:%s", id, name)
		log.Fatal("这是一个Fatal日志")
	}

}
