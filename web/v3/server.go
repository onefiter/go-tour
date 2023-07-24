package v3

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	// Handler 需要实现ServeHTTP
	http.Handler

	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// AddRoute 注册一个路由
	// method 是 HTTP 方法
	// path 是路径，必须以 / 为开头
	addRoute(method string, path string, handler HandleFunc)

	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
}

// NewHTTPServer 创建一个HTTPServer
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{router: newRouter()}
}

// ServeHTTP实现http.Handler接口interface中的ServeHTTP方法
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	s.serve(ctx)

}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) serve(ctx *Context) {
	n, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)

	if !ok || n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}

	n.handler(ctx)

}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) POST(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}
