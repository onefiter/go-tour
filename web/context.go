package web

import (
	"net/http"
	"net/url"
)

type Context struct {
	Req *http.Request
	//Resp 原生的ResponseWriter，当你直接使用Resp的时候
	// 那么相当于你绕开了 RespStatusCode 和 RespData。
	// 响应数据直接被发送到前端，其它中间件将无法修改响应
	// 其实我们也可以考虑将这个做成私有的
	Resp http.Response

	// 缓存的响应部分
	// 这部分数据会在最后刷新
	RespStatusCode int

	// RespData []byte
	RespData []byte

	PathParams map[string]string
	// 命中的路由
	MatchedRoute string

	// 缓存的数据
	cacheQueryValues url.Values

	// 页面渲染的引擎
	//tplEngine TemplateEngine

	// 用户可以自由决定在这里存储什么，
	// 主要用于解决在不同 Middleware 之间数据传递的问题
	// 但是要注意
	// 1. UserValues 在初始状态的时候总是 nil，你需要自己手动初始化
	UserValues map[string]any
}
