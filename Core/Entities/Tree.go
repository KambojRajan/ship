package entities

type Node struct {
	Name string
	Hash string
	Type string
}

type Tree struct {
	Nodes []Node
	Hash  string
}
