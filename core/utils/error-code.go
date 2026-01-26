package utils

// Error message constants for the Ship VCS
const (
	// Init command errors
	ErrFailedToGetWorkingDir      = "failed to get working directory: %w"
	ErrFailedToAccessPath         = "failed to access path: %w"
	ErrPathNotDirectory           = "path is not a directory: %s"
	ErrFailedToCreateObjectsDir   = "failed to create objects directory: %w"
	ErrFailedToCreateRefsHeadsDir = "failed to create refs/heads directory: %w"
	ErrFailedToCreateRefsTagsDir  = "failed to create refs/tags directory: %w"
	ErrFailedToCreateIndexFile    = "failed to create index file: %w"
	ErrFailedToCreateHEADFile     = "failed to create HEAD file: %w"

	// Cat-file command errors
	ErrInvalidObjectFormat = "invalid Object Format"
	ErrInvalidObjectHeader = "invalid Object Header"
	ErrTreeNotImplemented  = "to be impl"
	ErrUnknownType         = "unknown type %s"
)
