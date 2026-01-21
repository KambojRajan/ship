package entities

import "fmt"

type Node struct {
	Mode uint32
	Name string
	Hash [20]byte
}

func (n Node) modeString() string {
	return fmt.Sprintf("%o", n.Mode)
}
