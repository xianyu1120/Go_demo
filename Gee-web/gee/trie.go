package gee

import (
	"fmt"
	"strings"
)

/*实现功能：
*参数匹配`:`：如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
通配`*`。例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，
这种模式常用于静态服务器，能够递归地匹配子路径。
*/
//trie 树结构实现
type node struct{
	pattern string //待匹配路由规则，已有路由规则
	part string //当前节点对应的路径中的字符串
	children []*node//子节点索引
	isWild bool//是否精确匹配
}
func (n *node)String()string {
	return fmt.Sprintf("node{pattern=%s,part=%s,isWild=%t",n.pattern,n.part,n.isWild)
}
//第一个匹配成功的节点，用于开始插入节点继续遍历
//匹配子节点，根据子节点查找，匹配后返回
func (n *node) matchChild (part string)*node{
	for _, child := range n.children {//遍历节点
		if child.part==part||child.isWild {//如果节点相同或者精确匹配
			return child	//返回这个节点
		}	
	}
	return nil
}
//所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node{
	nodes:=make([]*node, 0)
	for _, child := range n.children {
		if child.part==part||child.isWild {
			nodes=append(nodes,child)//子节点存在时存储
		}
	}
	return nodes
}
//节点插入，pattern 已有节点，parts为按‘\'拆分后的各个节点，heigth为根节点的深度
//递归查找每一层节点，如果没有匹配到part的节点，就新建一个
func (n *node) insert(pattern string,parts[]string,height int) {
	if len(parts)==height {//路径没有：*等符号，拆分后长度不变；或者长度都为0说明为根路径
		//传递过来的就是完整路径，根路径或者不带*:的完整路径
		n.pattern=pattern//根下的子节点就是传递过来的节点
		return
	}
	//长度缩减，说明中间去掉的有*：所以将指定节点保留
	part:=parts[height]//指定节点0
	child:=n.matchChild(part)//查询该节点
	if child==nil {//为空，添加
		child=&node{part:part,isWild: part[0]==':'||part[0]=='*'}
		n.children=append(n.children,child)//添加该节点
	}
	child.insert(pattern,parts,height+1)//0+1，递归查找
}

//查询 height从0开始
func (n *node) search(parts[]string,height int)*node {
	if len(parts)==height||strings.HasPrefix(n.part,"*") {//路径为空或者为*
		if n.pattern=="" {//判断路由是否匹配
			return nil
		}
		return n
	}
	part:=parts[height]
	children:=n.matchChildren(part)
	for _, child := range children {
		result:=child.search(parts,height+1)
		if result!=nil {
			return result
		}
	}
	return nil
}
func(n *node) travel(list *([]*node))  {
	if n.pattern!=""{//子节点不为空
		*list=append(*list,n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
