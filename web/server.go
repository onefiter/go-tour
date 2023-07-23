package web

import "net/http"

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
}

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(add string) error

	// 增加路由注册功能
	addRoute(method string, path string, handler HandlerFunc)
}

func (h *HTTPServer) addRoute(method string, path string, handler HandlerFunc) {
	//TODO implement me
	panic("implement me")
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 你的Web框架代码
	panic("implement me")
}

func (h *HTTPServer) Start(add string) error {
	//TODO implement me
	panic("implement me")
}

type HandlerFunc func(ctx *Context)
