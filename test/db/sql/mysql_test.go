package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
)

func mysqlDB() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	return db
}

func TestConn(t *testing.T) {
	DB := mysqlDB()

	if err := DB.Ping(); err != nil {
		t.Fatal(err)
	}
	log.Info("connect to mysql successfully.")
}

type User struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
	Money int32 `json:"money"`
}

func TestQuery(t *testing.T) {
	DB := mysqlDB()

	var user User
	rows, err := DB.Query("SELECT * FROM user")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Name, &user.Age, &user.Money); err != nil {
			log.Error(err)
		}
		log.Info(user)
	}
}

func TestAdd(t *testing.T) {
	DB := mysqlDB()

	tx, err := DB.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO user (`name`, `age`) VALUES (?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	res, err := stmt.Exec("李四", 30)
	if err != nil {
		t.Fatal(err)
	}

	tx.Commit()

	log.Info(res.LastInsertId())
}

func TestSingleConn(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	row, err := db.Query("SELECT * FROM user where name='张三'")
	if err != nil {
		t.Fatal(err)
	}

	var user User

	//
	//for row.Next() {
	//	if err := row.Scan(&user.Name, &user.Age, &user.Money); err != nil {
	//		t.Fatal(err)
	//	}
	//	t.Log(user)
	//}

	row, err = db.Query("SELECT * FROM user where name='李四'")
	if err != nil {
		t.Fatal(err)
	}
	for row.Next() {
		if err := row.Scan(&user.Name, &user.Age, &user.Money); err != nil {
			t.Fatal(err)
		}
		t.Log(user)
	}
}