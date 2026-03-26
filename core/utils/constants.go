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

	ShipBlue  = "\033[34m"       // staged
	Sand      = "\033[33m"       // unstaged (safe)
	Sand256   = "\033[38;5;223m" // prettier sand
	ShipRed   = "\033[31m"       // deleted lines in diff
	ShipGreen = "\033[32m"       // added lines in diff
)

const (
	Staged   = "staged"
	Unstaged = "unstaged"
)

const (
	GitFileModeRegular    uint32 = 0o100644
	GitFileModeExecutable uint32 = 0o100755
)

var RunPurgeCommands = []string{"commit", "add"}
