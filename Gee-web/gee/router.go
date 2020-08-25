package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node //各种请求方式的Trie树根节点
	handlers map[string]HandlerFunc//根据访问方法返回处理器
}

//newRouter 构造函数
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}
//解析模式，将完整路径信息拆分为字符切片
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") //用“/”分割为字符切片

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item) //添加
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
//addRoute 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)//完整路径拆分

	key := method + "-" + pattern//方法-路由
	_, ok := r.roots[method]//方法对应的路由
	if !ok {
		r.roots[method] = &node{}//没有对应的路径，就添加
	}
	r.roots[method].insert(pattern, parts, 0)//有对应路径节点，就插入
	r.handlers[key] = handler//存储 方法-url:处理器
}
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)//拆分后的路径
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {//method对应路由不存在
		return nil, nil
	}

	n := root.search(searchParts, 0)//检索节点是否匹配

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
//得到所有的路由信息
func (r *router)getRoutes(method string)[]*node {
	root,ok:=r.roots[method]
	if !ok {
		return nil
	}
	nodes:=make([]*node,0)
	root.travel(&nodes)
	return nodes
}
//handle 
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {//没有找到对应路由规则
		key := c.Method + "-" + n.pattern//方法-url
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])//添加新的路由规则
	} else {//存在
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
