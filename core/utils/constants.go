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

var RunPurgeCommands = []string{"commit", "add", "agent"}

const (
	TRACE_SHORT = "Trace internal execution of a ship command"
	TRACE_LONG  = `Trace records and prints every internal pipeline step for the given command.
			SUPPORTED COMMANDS
			  commit <message>          trace the commit pipeline
			  add    <path> [paths...]  trace the staging pipeline
			  status [path]             trace the status pipeline
			
			OUTPUT FORMATS
			  --format text   (default) human-readable colour output + summary footer
			  --format json   NDJSON – one JSON object per step; pipe to jq for filtering
			
			OPENTELEMETRY
			  --otel                    export spans to an OTel-compatible JSON file
			  --otel-output <file>      destination file (default: stderr)
			
			EXAMPLES
			  ship trace commit "feat: add auth"
			  ship trace add src/ README.md
			  ship trace status
			  ship trace commit "fix: bug" --format json | jq 'select(.status=="error")'
			  ship trace commit "feat: observability" --otel --otel-output spans.json`
)
