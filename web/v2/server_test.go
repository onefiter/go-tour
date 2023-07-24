package v2

import "testing"

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello World!"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello, User!"))
	})

	s.Start(":8081")
}
