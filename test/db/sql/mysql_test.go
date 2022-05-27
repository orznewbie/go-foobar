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
		mysqldb = sqlDB("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	})
	_, err := mysqldb.ExecContext(context.Background(), CreateUserTable)
	if err != nil {
		panic(err)
	}
}

func sqlDB(driver, dns string) *sqlx.DB {
	db, err := sqlx.Connect(driver, dns)
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
	db.SetConnMaxLifetime(30 * time.Minute)
	return db
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

	// 单条插入更新直接使用tx.Exec
	_, err = tx.ExecContext(context.Background(), "INSERT INTO user (`name`, `age`) VALUES(?, ?)", "临时用户", 100)

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}

func TestTxBatchAdd(t *testing.T) {
	tx, err := mysqldb.Begin()
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
		_, err := stmt.ExecContext(context.Background(), "用户"+strconv.Itoa(i), i*10)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}
