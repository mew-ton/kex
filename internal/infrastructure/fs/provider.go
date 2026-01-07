package fs

// DocumentProvider defines the strategy for loading and fetching documents
type DocumentProvider interface {
	// Load retrieves the index schema from the source
	Load() (*IndexSchema, []error)
	// FetchContent retrieves the raw content for a specific path
	FetchContent(path string) (string, error)
}
