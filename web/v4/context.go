package v4

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
