package utils

const (
	UnknownTypeError = "Unknown type %s"
)

const (
	RootShipDir   = ".ship"
	RootObjectDir = ".ship/objects"
	RootIndexPath = ".ship/index"
	RootHEADPath  = ".ship/HEAD"
)

const (
	BLOB      = "blob"
	COMMIT    = "commit"
	TREE      = "tree"
	AUTHOR    = "author"
	COMMITTER = "committer"
	MESSAGE   = "message"
	PARENT    = "parent"
)

const (
	SEPARATOR = "/"
	NEWLINE   = '\n'
)

const (
	CatFileFormatPretty = "-p"
	CatFileFormatTree   = "-t"
	CatFileFormatCommit = "-c"
	CatFileContentSize  = "-s"
)
