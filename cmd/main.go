package main

import (
	"fmt"
	"github.com/eapache/queue"
)

type (
	From struct {
		Source string
		Kind   string
	}
	graph map[From]*adj

	adj struct {
		to targets
	}
	targets map[string]struct{}
)

func newAdj() *adj {
	return &adj{to: make(map[string]struct{})}
}

func (t targets) add(tar string) {
	t[tar] = struct{}{}
}

var g graph = make(map[From]*adj)

const (
	LocatedIn = "located_in"
	IsPartOf  = "is_part_of"
	Supply    = "supply"
)

func init() {
	g[From{Source: "A", Kind: LocatedIn}] = newAdj()
	g[From{Source: "A", Kind: LocatedIn}].to.add("B")
	g[From{Source: "A", Kind: LocatedIn}].to.add("C")

	g[From{Source: "A", Kind: IsPartOf}] = newAdj()
	g[From{Source: "A", Kind: IsPartOf}].to.add("E")

	g[From{Source: "B", Kind: IsPartOf}] = newAdj()
	g[From{Source: "B", Kind: IsPartOf}].to.add("D")

	g[From{Source: "B", Kind: Supply}] = newAdj()
	g[From{Source: "B", Kind: Supply}].to.add("D")

	g[From{Source: "B", Kind: Supply}] = newAdj()
	g[From{Source: "B", Kind: Supply}].to.add("E")
}

func main() {
	tar := bfs("A", 1, []string{LocatedIn, IsPartOf})
	for t := range tar {
		fmt.Println(t)
	}
}

func bfs(s string, maxSteps int, overs []string) targets {
	depth := 0
	var visited = make(map[string]struct{})
	q := queue.New()
	q.Add(s)
	visited[s] = struct{}{}
	for q.Length() != 0 {
		if depth >= maxSteps {
			break
		}
		depth++
		size := q.Length()

		for i := 0; i < size; i++ {
			v := q.Peek().(string)
			q.Remove()
			for _, over := range overs {
				if g[From{Source: v, Kind: over}] == nil {
					continue
				}
				for next := range g[From{Source: v, Kind: over}].to {
					if _, ok := visited[next]; !ok {
						q.Add(next)
						visited[next] = struct{}{}
					}
				}
			}
		}
	}

	return visited
}
