package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
)

func mysqlDB() *sql.DB {
	DB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	return DB
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
		if err := rows.Scan(&user.Name, &user.Age); err != nil {
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
