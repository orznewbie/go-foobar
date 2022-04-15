package dqlx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuery(t *testing.T) {
	var1 := VarAs("V", Has("name")).Select("uid")

	test := Query(Uid("V", "0x123", "0x124")).
		Filter(
			Or(
				And(
					Lt("age", 20),
					Ge("height", 1.80),
				),
				Eq("gender", "female"),
			)).
		Recurse(10, true).
		First(3).Offset(5).
		Name("test").
		OrderAsc("asset").
		OrderDesc("age").
		Select("uid").
		Edge(
			"follow",
			Select(Alias("follow_count", Count("uid"))),
			Facets(Alias("follow_time", "time")),
			Facets(Eq("type", "favourite")),
		).
		Edge(
			"follow->subscribe",
			Select(Expand("_all_")),
			Facets(),
			Cursor{First: 10},
		)

	query, variables, err := ToDql(var1, test)
	fmt.Println(query)

	require.NoError(t, err)

	expected := `
	query Variable_Test($0:int, $1:int, $2:int, $3:float, $4:string, $5:int, $6:bool, $7:string, $8:int) {
		<V> AS var(func: has(<name>)) {
			<uid>
		}
	
		<test>(func: uid(<V>,0x123,0x124), orderasc: <asset>, orderdesc: <age>, first: $0, offset: $1) @filter(((lt(<age>,$2) AND ge(<height>,$3)) OR eq(<gender>,$4))) @recurse(depth: $5,loop: $6) {
			<uid>
			<follow> @facets(<follow_time> : <time>) @facets(eq(<type>,$7)) {
				<follow_count> : count(<uid>)
				<subscribe> (first: $8) @facets {
					expand(_all_)
				}
			}
		}
	}`
	require.Equal(t, Minify(expected), Minify(query))

	require.Equal(t, map[string]string{
		"$0": "3",
		"$1": "5",
		"$2": "20",
		"$3": "1.8",
		"$4": "female",
		"$5": "10",
		"$6": "true",
		"$7": "favourite",
		"$8": "10",
	}, variables)
}
