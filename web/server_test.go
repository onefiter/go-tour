package web

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server
	// 自定义handler就是我们自定义的服务器
	// http.ListenAndServe(":8080", "自定义handler")
	http.ListenAndServe(":8080", s)

	http.ListenAndServeTLS(":443", "", "", s)

	s.Start(":8000")
}
