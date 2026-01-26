package utils

const (
	UnknownTypeError = "Unknown type %s"
)

const (
	RootShipDir   = ".ship"
	RootObjectDir = ".ship/objects"
	RootIndexPath = ".ship/index"
	RootHEADPath  = ".ship/HEAD"
	MainHeadPath  = "refs/heads/main"
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

const (
	Reset = "\033[0m"

	ShipBlue = "\033[34m"       // staged
	Sand     = "\033[33m"       // unstaged (safe)
	Sand256  = "\033[38;5;223m" // prettier sand
)

const (
	Staged   = "staged"
	Unstaged = "unstaged"
)
