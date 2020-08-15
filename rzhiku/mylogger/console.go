package mylogger

import (
	"fmt"
	"time"
)

//Logger 日志结构体
type ConsoleLogger struct {
	Level LogLevel
}

//NewConsoleLog 构造函数
func NewConsoleLog(levelStr string) ConsoleLogger {
	//解析用户传输进来的字符串
	level, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{
		Level: level,
	}
}
func (c ConsoleLogger) enable(logLevel LogLevel) bool {
	return c.Level <= logLevel
}
func(c ConsoleLogger) log(lv LogLevel, format  string,a...interface{}) {
	if c.enable(lv) {
	msg:=fmt.Sprintf(format,a...)//对传进来的信息格式化
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName,funcName,lineNo, msg)
	}
}

//Debug
func (c ConsoleLogger) Debug(format  string,a...interface{}) {
		c.log(DEBUG,format,a ...)
}
func (c ConsoleLogger) Info(format  string,a...interface{}){
		c.log(INFO, format,a ...)
}
func (c ConsoleLogger) Trace(format  string,a...interface{}) {
		c.log(TRACE,format,a ...)
	
}
func (c ConsoleLogger) Warning(format  string,a...interface{}) {
		c.log(WARNING,format,a ...)
	
}
func (c ConsoleLogger) Error(format  string,a...interface{}) {
		c.log(ERROR, format,a ...)
}
func (c ConsoleLogger) Fatal(format  string,a...interface{}) {
	c.log(FATAL, format,a ...)

}
