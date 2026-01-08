package fs_test

import (
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProvider is a mock implementation of DocumentProvider
type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Load() (*fs.IndexSchema, []error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).([]error)
	}
	return args.Get(0).(*fs.IndexSchema), args.Get(1).([]error)
}

func (m *MockProvider) FetchContent(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func TestCompositeProvider_Load(t *testing.T) {
	t.Run("should aggregate documents from multiple providers", func(t *testing.T) {
		// Mock Provider 1
		p1 := new(MockProvider)
		schema1 := &fs.IndexSchema{
			Documents: []*fs.DocumentSchema{
				{ID: "DOC_1", Path: "p1/doc1.md"},
			},
		}
		p1.On("Load").Return(schema1, []error{})

		// Mock Provider 2
		p2 := new(MockProvider)
		schema2 := &fs.IndexSchema{
			Documents: []*fs.DocumentSchema{
				{ID: "DOC_2", Path: "p2/doc2.md"},
			},
		}
		p2.On("Load").Return(schema2, []error{})

		composite := fs.NewCompositeProvider([]fs.DocumentProvider{p1, p2})
		schema, errs := composite.Load()

		assert.Empty(t, errs)
		assert.Len(t, schema.Documents, 2)
		// Order depends on implementation, but IDs should be unique
		ids := []string{schema.Documents[0].ID, schema.Documents[1].ID}
		assert.Contains(t, ids, "DOC_1")
		assert.Contains(t, ids, "DOC_2")
	})

	t.Run("should report error on ID collision", func(t *testing.T) {
		// Mock Provider 1
		p1 := new(MockProvider)
		schema1 := &fs.IndexSchema{
			Documents: []*fs.DocumentSchema{
				{ID: "DUPLICATE_ID", Path: "p1/doc1.md"},
			},
		}
		p1.On("Load").Return(schema1, []error{})

		// Mock Provider 2
		p2 := new(MockProvider)
		schema2 := &fs.IndexSchema{
			Documents: []*fs.DocumentSchema{
				{ID: "DUPLICATE_ID", Path: "p2/doc2.md"},
			},
		}
		p2.On("Load").Return(schema2, []error{})

		composite := fs.NewCompositeProvider([]fs.DocumentProvider{p1, p2})
		schema, errs := composite.Load()

		assert.NotEmpty(t, errs)
		assert.Contains(t, errs[0].Error(), "duplicate document ID 'DUPLICATE_ID'")
		// Should still return valid documents processed so far or excluding duplicates?
		// Implementation returns what it collected. First one "wins" in collection, second triggers error.
		assert.Len(t, schema.Documents, 1)
		assert.Equal(t, "DUPLICATE_ID", schema.Documents[0].ID)
	})
}

func TestCompositeProvider_FetchContent(t *testing.T) {
	t.Run("should fetch content from correct provider", func(t *testing.T) {
		p1 := new(MockProvider)
		p1.On("Load").Return(&fs.IndexSchema{
			Documents: []*fs.DocumentSchema{{ID: "DOC_1", Path: "p1/doc1.md"}},
		}, []error{})
		p1.On("FetchContent", "p1/doc1.md").Return("Content 1", nil)

		p2 := new(MockProvider)
		p2.On("Load").Return(&fs.IndexSchema{
			Documents: []*fs.DocumentSchema{{ID: "DOC_2", Path: "p2/doc2.md"}},
		}, []error{})
		p2.On("FetchContent", "p2/doc2.md").Return("Content 2", nil)

		composite := fs.NewCompositeProvider([]fs.DocumentProvider{p1, p2})
		composite.Load() // Populate map

		content, err := composite.FetchContent("p1/doc1.md")
		assert.NoError(t, err)
		assert.Equal(t, "Content 1", content)

		content, err = composite.FetchContent("p2/doc2.md")
		assert.NoError(t, err)
		assert.Equal(t, "Content 2", content)
	})
}
