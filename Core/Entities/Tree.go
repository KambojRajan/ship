package entities

import "fmt"

type Node struct {
	Mode uint32
	Name string
	Hash [20]byte
}

type Tree struct {
	Nodes []Node
}

func (n Node) ModeString() string {
	return fmt.Sprintf("%o", n.Mode)
}
