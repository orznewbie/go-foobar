package mariadb

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx/types"
)

type User struct {
	Id     int64          `db:"id"`
	Name   string         `db:"usr_name"`
	Facets types.JSONText `db:"facets"`
}

func TestBatchUpdateJSON(t *testing.T) {
	//JSON := `{"k":"v","num":135}`
	//byt, _ := json.Marshal(JSON)
	//byt = bytes.TrimPrefix(byt, []byte{'{'})
	//byt = bytes.TrimSuffix(byt, []byte{'}'})
	//byt = bytes.ReplaceAll()

	_, err := mariadb.ExecContext(context.Background(), `REPLACE INTO user(id,usr_name,facets) 
		VALUES (1996,'hhl',JSON_SET(facets,'$.k',100))`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	var user User
	err := mariadb.GetContext(context.Background(), &user, "SELECT id, JSON_EXTRACT(facets,'$.weight', '$.height') facets from user WHERE id=1000")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}
