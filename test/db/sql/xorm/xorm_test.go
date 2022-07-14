package xorm

import (
	"testing"

	"xorm.io/xorm"
)

type User struct {
	Id   int64
	Name string `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
}

func TestXORM(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:13306)/test")
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.CreateTables(new(User)); err != nil {
		t.Fatal(err)
	}
	engine.Query("select * from users")
}
