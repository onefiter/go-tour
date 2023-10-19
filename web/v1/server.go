package v1

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// AddRoute 注册一个路由
	// method 是 HTTP 方法
	// path 是路径，必须以 / 为开头
	AddRoute(method string, path string, handler HandleFunc)

	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
}

// 组合了http包server.go
// ServeHTTP实现http.Handler接口interface中的ServeHTTP方法
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	// ServeHTTP方法是作为http包和Web框架的关联点，需要在ServeHTTP内部，执行：
	// 1.构建起Web框架的上下文
	// 2.查找路由树，并执行命中的路由代码
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	// 查找路由，并执行代码
	s.serve(ctx)

}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) AddRoute(method string, path string, handler HandleFunc) {
	//TODO implement me
	panic("implement me")
}
func (s *HTTPServer) serve(ctx *Context) {

}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.AddRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) POST(path string, handler HandleFunc) {
	s.AddRoute(http.MethodPost, path, handler)
}
