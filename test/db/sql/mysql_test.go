package sql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/orznewbie/gotmpl/pkg/log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func mysqlDB() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	// 最大的连接数量maxOpen
	// numOpen是已打开的连接数量，所以numOpen <= maxOpen
	db.SetMaxOpenConns(5)
	// maxIdleCount是连接池中最大的可以重复使用的连接数量，默认是2
	// maxIdleCount <= maxOpen，如果maxIdleCount设置的比maxOpen还大，会自动调整maxIdleCount = maxOpen
	// freeConn []*driverConn为可重复使用的连接，所以len(freeConn) <= maxIdleOpen
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(30*time.Minute)
	return db
}

func TestConn(t *testing.T) {
	db := mysqlDB()
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	log.Info("pong!")
}

func TestMaxOpenConns(t *testing.T) {
	db := mysqlDB()
	db.SetMaxOpenConns(1)

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 1; i <= 3;i++ {
		go func(n int) {
			tx, err := db.Begin()
			if err != nil {
				panic(err)
			}
			defer tx.Rollback()

			log.Infof("open tx %d", n)
			time.Sleep(2 * time.Second)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

type User struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
	Money int32 `json:"money"`
}

func TestTxQuery(t *testing.T) {
	db := mysqlDB()

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()		// 开启后事务立即defer Rollback，不处理Rollback的error，这样即使事务提交也没问题

	var user User
	rows, err := tx.QueryContext(context.Background(), "SELECT * FROM user")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Warnf("rows close error:%v", err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&user.Name, &user.Age, &user.Money); err != nil {
			log.Error(err)
		}
		log.Info(user)
	}

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}

func TestCommitNowAdd(t *testing.T) {
	db := mysqlDB()
	// 直接使用db.Exec是默认将事务提交了
	_, err := db.ExecContext(context.Background(), "INSERT INTO user (`name`, `age`) VALUES(?, ?)", "CommitNow", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTxAdd(t *testing.T) {
	db := mysqlDB()

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	// 单条插入更新直接使用tx.Exec
	_, err = tx.ExecContext(context.Background(), "INSERT INTO user (`name`, `age`) VALUES(?, ?)", "复制人0", 0)

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}

func TestTxBatchAdd(t *testing.T) {
	db := mysqlDB()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	// Batch批量插入更新应该使用Prepare语句，减少数据库连接次数
	stmt, err := tx.PrepareContext(context.Background(), "INSERT INTO user (`name`, `age`) VALUES(?, ?)")
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 3; i++ {
		_, err := stmt.ExecContext(context.Background(), "复制人" + strconv.Itoa(i), i * 10000)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}