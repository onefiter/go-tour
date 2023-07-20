package beego

import "github.com/beego/beego/v2/server/web"

type UserController struct {
	web.Controller
}

type User struct {
	Name string
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("你好，go-tour")
}

func (c *UserController) CreateUser() {
	u := &User{}
	err := c.Ctx.BindJSON(u)

	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	_ = c.Ctx.JSONResp(u)
}