package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
*涉及上下文context 封装Request和Response，简化代码量
func(w http.ResponseWriter, req *http.Request)
 变为func(c *gee.Context)
 提供查询query /postForm参数的功能
 封装了html/string/json函数
 
*/


//H map[string]interface{}的别名
type H map[string]interface{}

//Context 对*http.Request和http.ResponseWrite的封装
type Context struct{
	Write http.ResponseWriter
	Req *http.Request
	Path string//请求路径
	Method string//请求方法
	Params map[string]string//
	StatusCode int//状态码
	handlers []HandlerFunc
	index int
	engine *Engine
}
func newContext(w http.ResponseWriter,req *http.Request) *Context {
	return &Context{
		Write: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		index: -1,//记录执行到第几个中间件
	}
}
func (c *Context)Next(){
	c.index++
	s:=len(c.handlers)
	for ;c.index<s;c.index++{
		c.handlers[c.index](c)
	}
}
func (c *Context) Fail(code int,err string) {
	c.index=len(c.handlers)
	c.JSON(code,H{"message":err})
}

func(c *Context)Param(key string)string{
	value,_:=c.Params[key]
	return value
}
//PostForm 返回请求的参数
func (c *Context) PostForm(key string)string  {
	return c.Req.FormValue(key)//根据表单key取到对应方法
}
//Query 请求url参数是否相等
func (c *Context) Query(key string)string {
	return c.Req.URL.Query().Get(key)
}
//Status 设置响应状态码
func (c *Context) Status(code int){
	c.StatusCode=code
	c.Write.WriteHeader(code)
}
//SetHeader 设置报头类型和值
//string: "Content-Type","text/plain"
//json: "Content-Type","application/json"
//html: "Content-Type","text/html"
func (c *Context) SetHeader(key string,value string) {
	c.Write.Header().Set(key,value)
}
//String string类型时设置响应报头信息和状态码
func (c *Context) String(code int,format string,values ...interface{}) {
	c.SetHeader("Content-Type","text/plain")
	c.Status(code)
	c.Write.Write([]byte(fmt.Sprintf(format,values...)))
}
//JSON json类型时响应报头以及状态码
func (c *Context) JSON(code int,obj interface{}) {
	c.SetHeader("Content-Type","application/json")//设置报头类型
	c.Status(code)
	encoder:=json.NewEncoder(c.Write)//编码
	if err:=encoder.Encode(obj);err != nil {
		http.Error(c.Write,err.Error(),500)
	}
}
//Data 传入数据时返回响应状态码和数据
func (c *Context) Data(code int,data[]byte) {
	c.Status(code)
	c.Write.Write(data)
}
//HTML 返回状态码，文件信息和内容
func (c *Context) HTML(code int,name string,data interface{}) {
	c.SetHeader("Content-Type","text/html")
	c.Status(code)//状态码
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Write, name, data); err != nil {//渲染内容
		c.Fail(500, err.Error())
	}
}
