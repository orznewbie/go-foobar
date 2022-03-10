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
		Paginate(Cursor{
			First:  3,
			Offset: 5,
		}).
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

	query, variables, err := QueriesToDQL(var1, test)
	fmt.Println(query)

	require.NoError(t, err)

	expected := `
	query Variable_Test($0:string, $1:string, $2:int, $3:int, $4:int, $5:float, $6:string, $7:int, $8:bool, $9:string, $10:int) {
		<V> AS var(func: has(<name>)) {
			<uid>
		}

		<test>(func: uid(<V>,$0,$1), orderasc: <asset>, orderdesc: <age>, first: $2, offset: $3) @filter(((lt(<age>,$4) AND ge(<height>,$5)) OR eq(<gender>,$6))) @recurse(depth: $7,loop: $8) {
			<uid>
			<follow> @facets(<follow_time> : <time>) @facets(eq(<type>,$9)) {
				<follow_count> : count(<uid>)
				<subscribe> (first: $10) @facets {
					expand(_all_)
				}
			}
		}
	}`
	require.Equal(t, Minify(expected), Minify(query))

	require.Equal(t, map[string]string{
		"$0":  "0x123",
		"$1":  "0x124",
		"$2":  "3",
		"$3":  "5",
		"$4":  "20",
		"$5":  "1.8",
		"$6":  "female",
		"$7":  "10",
		"$8":  "true",
		"$9":  "favourite",
		"$10": "10",
	}, variables)
}
