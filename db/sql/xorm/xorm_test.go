package xorm

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx/types"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type User struct {
	Id     int64
	Name   string         `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
	Facets types.JSONText `xorm:"json"`
}

func TestXORM(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:13306)/test")
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Sync(new(User)); err != nil {
		t.Fatal(err)
	}
	user := &User{
		Id:     1000,
		Name:   "ttt",
		Facets: types.JSONText(`{"height":1.80,"weight":150}`),
	}
	affected, err := engine.Insert(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(affected)
}

func TestUpdate(t *testing.T) {

}

func TestQuery(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:13306)/test")
	if err != nil {
		t.Fatal(err)
	}
	var user User
	_, err = engine.Table("user").Where("id=1000").Get(&user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}
