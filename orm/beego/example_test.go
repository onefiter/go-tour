package beego

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type User struct {
	ID   int    `orm:"column(id)"`
	Name string `orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(User))

	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "beego.db")

}

func TestCURD(t *testing.T) {
	orm.RunSyncdb("default", false, true)

	o := orm.NewOrm()
	user := new(User)
	user.Name = "Mike"

	o.Insert(user)
}
