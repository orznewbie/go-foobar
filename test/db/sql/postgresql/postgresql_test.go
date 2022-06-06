package postgresql

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/orznewbie/gotmpl/pkg/log"
	testsql "github.com/orznewbie/gotmpl/test/db/sql"
	"strconv"
	"testing"
)

const (
	CreateUserTable = `
		CREATE TABLE IF NOT EXISTS userx (
			id SERIAL,
			name VARCHAR(30),
			age INT,
			schools VARCHAR(10)[],
			hobbies JSON,
			PRIMARY KEY (id)
		);
	`
)

var (
	pgdb *sqlx.DB
)

func init() {
	pgdb = testsql.NewDB("postgres",
		"host=127.0.0.1 port=5432 user=postgres password=123456 dbname=test sslmode=disable")
	if _, err := pgdb.ExecContext(context.Background(), CreateUserTable); err != nil {
		panic(err)
	}
}

func TestPing(t *testing.T) {
	if err := pgdb.Ping(); err != nil {
		t.Fatal(err)
	}
	log.Info("pong!")
}

type Userx struct {
	Id      string   `db:"id"`
	Name    string   `db:"name"`
	Age     int      `db:"age"`
	Schools pq.StringArray `db:"schools"`
	Hobbies string `db:"hobbies"`
}

func TestAdd(t *testing.T) {
	ctx := context.Background()
	tx, err := pgdb.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	var (
		q    = `INSERT INTO userx (name, age, schools, hobbies) VALUES(:name, :age, :schools, :hobbies)`
		args []map[string]interface{}
	)

	for i := 1; i <= 10; i++ {
		args = append(args, map[string]interface{}{
			"name":    "张三" + strconv.Itoa(i),
			"age":     i * 2,
			"schools": pq.Array([]string{"小学", "中学"}),
			"hobbies": `{"football":"11","computer":"bar"}`,
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

func TestNullValueAdd(t *testing.T) {
	ctx := context.Background()
	tx, err := pgdb.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	var (
		q    = `INSERT INTO userx (name, age, schools, hobbies) VALUES(:name, :age, :schools, :hobbies)`
		args []map[string]interface{}
	)

	args = append(args, map[string]interface{}{
		"name":    "张三",
		"age":     nil,
		"schools": pq.Array([]string{"小学", "中学"}),
		"hobbies": `{"football":"11","computer":"bar"}`,
	})

	if _, err := tx.NamedExecContext(ctx, q, args); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Warnf("tx commit error: %v", err)
		return
	}
}

func TestQuery(t *testing.T) {
	var users []Userx
	if err := pgdb.SelectContext(context.Background(), &users,
		"SELECT * FROM userx LIMIT 10 OFFSET 0"); err != nil {
		t.Fatal(err)
	}
	log.Info(users)
}