package generator

import (
	"encoding/json"
	"fmt"
	"io/fs"
)

type AgentFiles struct {
	Coding        []string `json:"coding"`
	Documentation []string `json:"documentation"`
	Skills        []string `json:"skills"` // List of template patterns
}

type AgentDef struct {
	Name  string     `json:"name"`
	Files AgentFiles `json:"files"`
}

type Manifest struct {
	Kex      []string            `json:"kex"`
	AiAgents map[string]AgentDef `json:"ai-agents"`
}

// LoadManifest loads the manifest from the given filesystem.
// It expects manifest.json to be at the root of the provided fs.
func LoadManifest(fsys fs.FS) (*Manifest, error) {
	data, err := fs.ReadFile(fsys, "manifest.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest.json: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest.json: %w", err)
	}

	return &m, nil
}
