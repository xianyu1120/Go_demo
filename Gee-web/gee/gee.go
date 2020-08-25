package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

//HandlerFunc 定义请求处理程序
type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string        //前缀
		middlewares []HandlerFunc //中间件
		parent      *RouterGroup  //支持嵌套
		engine      *Engine       //所有组共享一个Engine实例
	}
	//Engine 实现ServeHttp的接口
	Engine struct {
		*RouterGroup
		//路由映射表，key由请求方法和静态路由地址构成GET-/
		router *router//路由规则
		groups []*RouterGroup //存储所有分组
		htmlTemplates *template.Template //用于渲染
		funcMap template.FuncMap
	}
)
//New gee.Engine的构造函数
//返回一个存有路由规则的信息接口
func New() *Engine {
	engine:=&Engine{router: newRouter()}
	engine.RouterGroup=&RouterGroup{engine: engine}
	engine.groups=[]*RouterGroup{engine.RouterGroup}
	return engine
}

//Group 定义组以创建新的RouterGroup；记住所有组共享同一个Engine实例
func(group *RouterGroup)Group(prefix string)*RouterGroup{
	engine:=group.engine
	newGroup:=&RouterGroup{
		prefix: group.prefix+prefix,
		parent: group,
		engine: engine,
	}

	engine.groups=append(engine.groups,newGroup)
	return newGroup
}
//addRoute 添加路由规则
//method 访问方法；
//pattern 模式
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern:=group.prefix+comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//GET 定义添加GET请求的方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)//添加GET方法对应的url以及处理器
	
}

//POST 定义添加POST请求的方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
//create 静态处理器
func (group *RouterGroup) createStaticHandler(relativePath string ,fs http.FileSystem)HandlerFunc {
	absolutePath:=path.Join(group.prefix,relativePath)
	fileServer:=http.StripPrefix(absolutePath,http.FileServer(fs))
	return func (c *Context)  {
		file:=c.Param("filepath")
		if _,err:=fs.Open(file);err != nil {
			c.Status(http.StatusNotFound)
			return 
		}
		fileServer.ServeHTTP(c.Write,c.Req)
	}
}
// 提供静态文件
func (group *RouterGroup)Static(relativePath string,root string)  {
	handler:=group.createStaticHandler(relativePath,http.Dir(root))
	urlPattern:=path.Join(relativePath,"/*filepath")
	group.GET(urlPattern,handler)
}
// for custom render function
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap= funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//Run 定义启动http server的方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
// Use 定义使用添加中间件到组
func (group *RouterGroup)Use (middlewares ...HandlerFunc){
	group.middlewares=append(group.middlewares,middlewares...)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path,group.prefix) {
			middlewares=append(middlewares,group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers=middlewares
	c.engine=engine
	engine.router.handle(c)
}
