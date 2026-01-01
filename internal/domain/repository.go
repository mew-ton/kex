package domain

// DocumentRepository defines the interface for accessing documents.
type DocumentRepository interface {
	Load() error
	GetAll() []*Document
	GetErrors() []error
	GetByID(id string) (*Document, bool)
	Search(keywords []string, scopes []string) []*Document
}
