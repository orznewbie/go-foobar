package aip

import (
	"fmt"
	"testing"

	user_v1 "github.com/orznewbie/gotmpl/api/user/v1"

	"go.einride.tech/aip/ordering"
	"go.einride.tech/aip/pagination"
)

func TestParseOrderBy(t *testing.T) {
	orderBy, err := ordering.ParseOrderBy(&user_v1.LiseUsersRequest{OrderBy: "foo asc,bar desc"})
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range orderBy.Fields {
		fmt.Println(field)
	}
}

func TestParsePagination(t *testing.T) {
	token, err := pagination.ParsePageToken(&user_v1.LiseUsersRequest{
		PageToken: "",
		PageSize:  18,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token.Offset)
}
