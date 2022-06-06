package sql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/orznewbie/gotmpl/pkg/log"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	CreateUserTable = `
		CREATE TABLE IF NOT EXISTS user(
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(30),
			age INT,
			PRIMARY KEY (id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
)

var (
	once    sync.Once
	mysqldb *sqlx.DB
)

func init() {
	once.Do(func() {
		mysqldb = NewDB("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	})
	_, err := mysqldb.ExecContext(context.Background(), CreateUserTable)
	if err != nil {
		panic(err)
	}
}

func TestPing(t *testing.T) {
	if err := mysqldb.Ping(); err != nil {
		t.Fatal(err)
	}
	log.Info("pong!")
}

func TestMaxOpenConns(t *testing.T) {
	mysqldb.SetMaxOpenConns(1)

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go func(n int) {
			tx, err := mysqldb.Begin()
			if err != nil {
				panic(err)
			}
			defer tx.Rollback()

			log.Infof("begin tx %d", n)
			time.Sleep(2 * time.Second)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

type User struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	Age  int32  `db:"age"`
}

func TestTxQuery(t *testing.T) {
	tx, err := mysqldb.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback() // 开启后事务立即defer Rollback，不处理Rollback的error，这样即使事务提交也没问题

	var user User
	rows, err := tx.QueryContext(context.Background(), "SELECT * FROM user LIMIT 0,10")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Warnf("rows close error:%v", err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			log.Error(err)
		}
		log.Info(user)
	}

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}

func TestSqlxQuery(t *testing.T) {
	var users []User
	if err := mysqldb.SelectContext(context.Background(), &users, "SELECT age,name FROM user WHERE age>10;"); err != nil {
		t.Fatal(err)
	}
	log.Info(users)
}

func TestCommitNowAdd(t *testing.T) {
	// 直接使用db.Exec是默认将事务提交了
	_, err := mysqldb.ExecContext(context.Background(), "INSERT INTO user (`name`, `age`) VALUES(?, ?)", "CommitNow", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTxAdd(t *testing.T) {
	tx, err := mysqldb.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(context.Background(), "INSERT INTO user (`name`, `age`) VALUES(?, ?)", "临时用户", 100)

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}

func TestTxBatchAdd(t *testing.T) {
	ctx := context.Background()
	tx, err := mysqldb.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	var (
		//q = `INSERT INTO user (name, age) VALUES(:name, :age)`
		q = `INSERT INTO c (valuec) VALUES(:value)`
		args []map[string]interface{}
	)

	for i := 1; i <= 10000; i++ {
		args = append(args, map[string]interface{}{
			"value":"a"+strconv.Itoa(i),
			//"age":1,
		})
	}

	if _, err := tx.NamedExecContext(ctx, q, args); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}