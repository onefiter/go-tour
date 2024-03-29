package v3

import (
	"fmt"
	"strings"
)

// 全静态路由匹配

type router struct {
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
	trees map[string]*node
}

func newRouter() router {
	return router{trees: map[string]*node{}}
}

func (r *router) addRoute(method string, path string, handler HandleFunc) {
	if path == "" {
		panic("web: 路由是空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}

	root, ok := r.trees[method]

	// 这是一个全新的 HTTP 方法，创建根节点
	if !ok {
		// 创建根节点
		// 此时根节点node.children, node.handler都为nil
		root = &node{path: "/"}
		r.trees[method] = root
	}

	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handler
		return
	}

	// 以 "/" 分割 path 如 /v1/user/info ,分割后===> [v1,user,info]
	segs := strings.Split(path[1:], "/")

	// 开始一段段处理
	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(s)
	}

	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}

	root.handler = handler
}

// findRoute 查找对应的节点
// 注意，返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return root, true
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, s := range segs {
		root, ok = root.childOf(s)
		if !ok {
			return nil, false
		}
	}
	return root, true
}

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 通配符匹配
// 这是不回溯匹配
type node struct {
	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符 * 表达节点，任意匹配
	starChild *node
}

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return n.starChild, n.starChild != nil
	}
	res, ok := n.children[path]
	if !ok {
		return n.starChild, n.starChild != nil
	}
	return res, ok
}

// childOrCreate 查找子节点，如果子节点不存在就创建一个
// 并且将子节点放回去了 children 中
func (n *node) childOrCreate(path string) *node {

	if path == "*" {
		if n.starChild == nil {
			n.starChild = &node{path: "*"}
		}
		return n.starChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}
