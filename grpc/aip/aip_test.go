package aip

import (
	"fmt"
	"strconv"
	"testing"

	"go.einride.tech/aip/fieldmask"
	"go.einride.tech/aip/filtering"
	"go.einride.tech/aip/ordering"
	"go.einride.tech/aip/pagination"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	userpb "github.com/orznewbie/go-foobar/api/user/v1"
	"github.com/orznewbie/go-foobar/pkg/dqlx"
)

func TestParseOrderBy(t *testing.T) {
	order := "test:Computer::name asc, test:Computer::orderNum desc"
	orderBy, err := ordering.ParseOrderBy(&userpb.LiseUsersRequest{OrderBy: order})
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range orderBy.Fields {
		fmt.Println(field)
	}
}

func TestParsePagination(t *testing.T) {
	req1 := &userpb.LiseUsersRequest{
		PageSize: 100,
		Skip:     2,
	}
	page1, err := pagination.ParsePageToken(req1)
	if err != nil {
		t.Fatal(err)
	}
	nextPageToken1 := page1.Next(req1).String()

	req2 := &userpb.LiseUsersRequest{
		PageToken: nextPageToken1,
		PageSize:  200,
		Skip:      3,
	}
	page2, err := pagination.ParsePageToken(req2)
	if err != nil {
		t.Fatal(err)
	}
	nextPageToken2 := page2.Next(req2).String()

	req3 := &userpb.LiseUsersRequest{
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

type Projection []string

func TestFieldMask(t *testing.T) {
	var masked userpb.User
	fieldmask.Update(
		&fieldmaskpb.FieldMask{Paths: []string{"id", "name"}},
		&masked,
		&userpb.User{
			Id:   1,
			Name: "188",
			Age:  100,
			Role: nil,
		},
	)
	fmt.Println(masked.Id, masked.Name, masked.Age)
}

func TestParseFilter(t *testing.T) {
	req := &userpb.LiseUsersRequest{
		Filter: `test:FooBar::a>10 AND test:Computer::b<100`,
	}

	var parser filtering.Parser
	parser.Init(req.Filter)
	res, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.String())

	//filtering.Walk(func(currExpr, parentExpr *expr.Expr) bool {
	//	switch c := currExpr.ExprKind.(type) {
	//	case *expr.Expr_IdentExpr:
	//		c.IdentExpr.Name = "prefix:" + c.IdentExpr.GetName()
	//		return false
	//	default:
	//		return true
	//	}
	//}, filter.GetExpr())
	//
	//fmt.Println(filter.Expr)
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
