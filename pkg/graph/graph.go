package graph

type Node struct {
}

type Edge struct {
}

type G interface {
	AddNode(node Node)
	AddEdge(edge Edge)
	CountNode() uint64
	CountEdge() uint64
	DeleteNode(uint642 uint64)
}
