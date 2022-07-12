package aip

import (
	"fmt"
	"testing"

	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"go.einride.tech/aip/filtering"

	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"go.einride.tech/aip/fieldmask"

	user_v1 "github.com/orznewbie/gotmpl/api/user/v1"

	"go.einride.tech/aip/ordering"
	"go.einride.tech/aip/pagination"
)

func TestParseOrderBy(t *testing.T) {
	orderBy, err := ordering.ParseOrderBy(&user_v1.LiseUsersRequest{OrderBy: "foo   asc  ,  bar desc"})
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range orderBy.Fields {
		fmt.Println(field)
	}
}

func TestParsePagination(t *testing.T) {
	req1 := &user_v1.LiseUsersRequest{
		PageSize: 100,
		Skip:     2,
	}
	page1, err := pagination.ParsePageToken(req1)
	if err != nil {
		t.Fatal(err)
	}
	nextPageToken1 := page1.Next(req1).String()

	req2 := &user_v1.LiseUsersRequest{
		PageToken: nextPageToken1,
		PageSize:  200,
		Skip:      3,
	}
	page2, err := pagination.ParsePageToken(req2)
	if err != nil {
		t.Fatal(err)
	}
	nextPageToken2 := page2.Next(req2).String()

	req3 := &user_v1.LiseUsersRequest{
		PageToken: nextPageToken2,
		PageSize:  50,
		Skip:      5,
	}
	page3, err := pagination.ParsePageToken(req3)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(page3.Offset)
}

func TestFieldMask(t *testing.T) {
	var masked user_v1.User
	fieldmask.Update(
		&fieldmaskpb.FieldMask{Paths: []string{"id", "name"}},
		&masked,
		&user_v1.User{
			Id:   1,
			Name: "188",
			Age:  100,
			Role: nil,
		},
	)
	fmt.Println(masked.Id, masked.Name, masked.Age)
}

func TestParseFilter(t *testing.T) {
	d, err := filtering.NewDeclarations(
		filtering.DeclareStandardFunctions(),
		filtering.DeclareIdent("a", filtering.TypeInt),
		filtering.DeclareIdent("b", filtering.TypeString),
		//filtering.DeclareIdent("c", filtering.TypeFloat),
		//filtering.DeclareIdent("d", filtering.TypeString),
		//filtering.DeclareFunction("regex", filtering.NewFunctionOverload("regex_string", filtering.TypeBool, filtering.TypeString, filtering.TypeString)),
	)
	if err != nil {
		t.Fatal(err)
	}
	req := &user_v1.LiseUsersRequest{
		Filter: `a > 10 AND b = "x"`,
	}
	filter, err := filtering.ParseFilter(req, d)
	if err != nil {
		t.Fatal(err)
	}

	macroDeclarations, err := filtering.NewDeclarations(
		filtering.DeclareStandardFunctions(),
		filtering.DeclareIdent("a", filtering.TypeInt),
		filtering.DeclareIdent("b", filtering.TypeString),
		filtering.DeclareFunction("gt", filtering.NewFunctionOverload("gt", filtering.TypeBool, filtering.TypeString, filtering.TypeInt)),
		filtering.DeclareFunction("eq", filtering.NewFunctionOverload("eq", filtering.TypeBool, filtering.TypeString, filtering.TypeString)),
		//filtering.DeclareIdent("c", filtering.TypeFloat),
		//filtering.DeclareIdent("d", filtering.TypeString),
		//filtering.DeclareFunction("regex", filtering.NewFunctionOverload("regex_string", filtering.TypeBool, filtering.TypeString, filtering.TypeString)),
	)
	if err != nil {
		t.Fatal(err)
	}

	macro, err := filtering.ApplyMacros(filter, macroDeclarations, func(cursor *filtering.Cursor) {
		callExpr := cursor.Expr().GetCallExpr()
		if callExpr == nil || callExpr.GetFunction() == filtering.FunctionAnd {
			return
		}
		if len(callExpr.Args) != 2 {
			return
		}
		switch callExpr.GetFunction() {
		case filtering.FunctionEquals:
			cursor.Replace(filtering.Function("eq", callExpr.GetArgs()...))
		case filtering.FunctionGreaterThan:
			cursor.Replace(filtering.Function("gt", callExpr.GetArgs()...))
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(macro.CheckedExpr.String())
}

func dfs(filter *expr.Expr) {
	//if filter.
}
