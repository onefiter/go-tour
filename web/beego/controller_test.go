package beego

import (
	"github.com/beego/beego/v2/server/web"
	"testing"
)

func TestUserController(t *testing.T) {
	// Beego中独有的配置参数
	web.BConfig.CopyRequestBody = true
	c := &UserController{}
	// Get请求
	web.Router("/user", c, "get:GetUser")
	web.Run(":8081")
}
