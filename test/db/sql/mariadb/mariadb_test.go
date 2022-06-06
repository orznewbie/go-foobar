package mariadb

import (
	"context"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/orznewbie/gotmpl/pkg/log"
	testsql "github.com/orznewbie/gotmpl/test/db/sql"
	"testing"
)

const (
	CreateVertexTable = `
		CREATE TABLE IF NOT EXISTS vertex(
			uid BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			create_time INT,
 			attributes BLOB,
			PRIMARY KEY (uid)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
)

type (
	Value json.RawMessage
	// 数据库的结果映射tag用db:
	Vertex struct {
		Uid        string `db:"uid"`
		CreateTime int    `db:"create_time"`
		Attributes Value  `db:"attributes"`
	}
)

var mariadb *sqlx.DB

func init() {
	mariadb = testsql.NewDB("mysql", "root:123456@tcp(127.0.0.1:13306)/test")
	_, err := mariadb.ExecContext(context.Background(), CreateVertexTable)
	if err != nil {
		panic(err)
	}
}

func TestPing(t *testing.T) {
	if err := mariadb.Ping(); err != nil {
		t.Fatal(err)
	}
	log.Info("pong!")
}

func TestDynamicQuery(t *testing.T) {
	q := `SELECT uid, COLUMN_JSON(attributes) AS attributes FROM vertex 
		WHERE COLUMN_GET(attributes, 'price' AS INTEGER)>50;`

	var vs []Vertex
	if err := mariadb.SelectContext(context.Background(), &vs, q); err != nil {
		t.Fatal(err)
	}

	for _, v := range vs {
		log.Info(v.Uid, v.CreateTime, string(v.Attributes))
	}
}

func TestDynamicAdd(t *testing.T) {
	q := `INSERT INTO vertex (create_time, attributes) 
		VALUES (?, COLUMN_CREATE('color', ?, 'price', ?));`
	q1 := `INSERT INTO vertex (create_time, attributes) 
		VALUES (:create_time, COLUMN_CREATE('color', :color, 'weight', :weight, 'deleted', :deleted));`

	ctx := context.Background()

	tx := mariadb.MustBeginTx(ctx, nil)
	defer tx.Rollback()

	_, err := tx.ExecContext(ctx, q, 1234, "white", 500)
	if err != nil {
		t.Fatal(err)
	}

	var args = []map[string]interface{}{
		{"create_time": 1234, "color": "black", "weight": 1.97, "deleted": false},
	}
	_, err = tx.NamedExecContext(ctx, q1, args)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestDynamicUpdate(t *testing.T) {
	q := `UPDATE vertex
		SET attributes=COLUMN_ADD(attributes, 'color', :color) 
		WHERE COLUMN_GET(attributes, 'price' AS INTEGER)>:price;`

	ctx := context.Background()

	tx := mariadb.MustBeginTx(ctx, nil)
	defer tx.Rollback()

	var args = []map[string]interface{}{
		{"color": "red", "price": 300},
	}
	_, err := tx.NamedExecContext(ctx, q, args)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestDynamicAddArray(t *testing.T) {
	q := `INSERT INTO vertex (create_time, attributes) 
		VALUES (:create_time, COLUMN_CREATE('dtdl:test:implements', :implements));`

	ctx := context.Background()

	tx := mariadb.MustBeginTx(ctx, nil)
	defer tx.Rollback()

	byt, _ := json.Marshal([]string{"space", "room"})
	var args = []map[string]interface{}{
		{"create_time": 666, "dtdl:test:implements": byt},
	}
	_, err := tx.NamedExecContext(ctx, q, args)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
