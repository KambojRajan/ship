package temp

type ObjectType int

const (
	BLOB ObjectType = iota
	TREE
	COMMIT
)

func (o ObjectType) String() string {
	switch o {
	case BLOB:
		return "blob"
	case TREE:
		return "tree"
	case COMMIT:
		return "commit"
	default:
		panic("unknown object type")
	}
}
