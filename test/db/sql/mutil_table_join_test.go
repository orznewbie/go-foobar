package sql

import (
	"context"
	"fmt"
	"testing"
)

// Object  <---extends---  Space  <---extends---  Room
const (
	CreateObjectTable = `
		CREATE TABLE IF NOT EXISTS object(
			uid BIGINT UNSIGNED AUTO_INCREMENT,
			create_time INT,
			model VARCHAR(40),
			PRIMARY KEY (uid)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`

	CreateSpaceTable = `
		CREATE TABLE IF NOT EXISTS space(
			uid BIGINT UNSIGNED,
			height INT,
			CONSTRAINT space_object_fk FOREIGN KEY (uid) REFERENCES object (uid) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`

	CreateRoomTable = `
		CREATE TABLE IF NOT EXISTS room(
			uid BIGINT UNSIGNED,
			capacity INT,
			CONSTRAINT room_object_fk FOREIGN KEY (uid) REFERENCES object (uid) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
)

func init() {
	once.Do(func() {
		mysqldb = sqlDB("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	})
	_, err := mysqldb.ExecContext(context.Background(), CreateObjectTable)
	if err != nil {
		panic(err)
	}
	_, err = mysqldb.ExecContext(context.Background(), CreateSpaceTable)
	if err != nil {
		panic(err)
	}
	_, err = mysqldb.ExecContext(context.Background(), CreateRoomTable)
	if err != nil {
		panic(err)
	}
}

func TestAddSpace(t *testing.T) {
	tx, err := mysqldb.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(context.Background(), `INSERT INTO object (create_time, model) VALUES (123456, 'space')`)
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.ExecContext(context.Background(), fmt.Sprintf(`INSERT INTO space (uid, height) VALUES (%d, 120)`, id))

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestAddRoom(t *testing.T) {
	tx, err := mysqldb.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(context.Background(), `INSERT INTO object (create_time, model) VALUES (199999, 'room')`)
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.ExecContext(context.Background(), fmt.Sprintf(`INSERT INTO space (uid, height) VALUES (%d, 800)`, id))
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.ExecContext(context.Background(), fmt.Sprintf(`INSERT INTO room (uid, capacity) VALUES (%d, 50)`, id))
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestJoinQuerySpace(t *testing.T) {
	//q := `select * from space s left join object o on o.uid=s.uid where height>100;`
	//rows, err := mysqldb.QueryContext(context.Background(), q)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer rows.Close()
	//cols, err := rows.Columns()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//var vs []Vertex
	//
	//for rows.Next() {
	//	var values = make([]interface{}, len(cols))
	//	for i := range values {
	//		values[i] = new(Value)
	//	}
	//	if err := rows.Scan(values...); err != nil {
	//		t.Fatal(err)
	//	}
	//	var v = Vertex{
	//		Uid:        "",
	//		CreateTime: 0,
	//		Attributes: make(map[string]Value),
	//	}
	//	for i := range values {
	//		if cols[i] == "uid" {
	//			//v.Uid = values[i]
	//		} else if cols[i] == "create_time" {
	//
	//		} else {
	//			v.Attributes[cols[i]] = *(values[i].(*Value))
	//		}
	//	}
	//	vs = append(vs, v)
	//}
	//
	//log.Info(vs)
}
