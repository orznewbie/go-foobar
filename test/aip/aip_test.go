package aip

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/orznewbie/gotmpl/pkg/dqlx"
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
	//d, err := filtering.NewDeclarations(
	//	filtering.DeclareStandardFunctions(),
	//	filtering.DeclareIdent("a", filtering.TypeBool),
	//	filtering.DeclareIdent("b", filtering.TypeString),
	//	filtering.DeclareIdent("c", filtering.TypeInt),
	//	//filtering.DeclareIdent("c", filtering.TypeFloat),
	//	//filtering.DeclareIdent("d", filtering.TypeString),
	//	//filtering.DeclareFunction("regex", filtering.NewFunctionOverload("regex_string", filtering.TypeBool, filtering.TypeString, filtering.TypeString)),
	//)
	//if err != nil {
	//	t.Fatal(err)
	//}
	req := &user_v1.LiseUsersRequest{
		Filter: `a?12 = 1.1700`,
	}

	var parser filtering.Parser
	parser.Init(req.Filter)
	filter, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}

	//callExpr := filter.GetExpr().GetExprKind().(*expr.Expr_CallExpr)

	//str := fmt.Sprintf("res: %v", getValue(callExpr.CallExpr.GetArgs()[1].GetConstExpr()))
	//fmt.Println(str)
	fmt.Println(filter.Expr)

	//dqlxFilter, err := dfs(filter.CheckedExpr.GetExpr())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//q, a, err := dqlxFilter.Dql()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("\nafter:")
	//fmt.Println(q)
	//fmt.Println(a)
}

func getValue(kind *expr.Constant) string {
	switch kindCast := kind.ConstantKind.(type) {
	case *expr.Constant_StringValue:
		return kindCast.StringValue
	case *expr.Constant_Int64Value:
		return strconv.FormatInt(kindCast.Int64Value, 10)
	case *expr.Constant_DoubleValue:
		return strconv.FormatFloat(kindCast.DoubleValue, 'E', -1, 64)

	}
	return ""
}

func dfs(ex *expr.Expr) (dqlx.Filter, error) {
	callExpr, ok := ex.GetExprKind().(*expr.Expr_CallExpr)
	if !ok {
		return nil, fmt.Errorf("cannot parse filter: %v", ex)
	}
	var (
		fn   = callExpr.CallExpr.GetFunction()
		args = callExpr.CallExpr.GetArgs()
	)
	if fn != filtering.FunctionAnd && fn != filtering.FunctionOr {
		switch fn {
		case filtering.FunctionEquals:
			return dqlx.Eq(args[0].GetIdentExpr().GetName(), args[1].GetConstExpr().GetStringValue()), nil
		case filtering.FunctionGreaterThan:
			return dqlx.Gt(args[0].GetIdentExpr().GetName(), args[1].GetConstExpr().GetInt64Value()), nil
		case filtering.FunctionLessThan:
			return dqlx.Gt(args[0].GetIdentExpr().GetName(), args[1].GetConstExpr().GetInt64Value()), nil
		}
	}

	if len(args) != 2 {
		return nil, fmt.Errorf("internal parse error")
	}
	left, err := dfs(args[0])
	if err != nil {
		return nil, err
	}
	right, err := dfs(args[1])
	if err != nil {
		return nil, err
	}
	if fn == filtering.FunctionAnd {
		return dqlx.And(left, right), nil
	}
	return dqlx.Or(left, right), nil
}
