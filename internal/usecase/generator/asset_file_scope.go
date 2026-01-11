package generator

// AssetFileScopeStrategy defines which files to select from an agent definition
type AssetFileScopeStrategy interface {
	SelectFiles(def AgentDef) []string
}

// CodingFileScope handles identifying coding-related files
type CodingFileScope struct{}

func (s *CodingFileScope) SelectFiles(def AgentDef) []string {
	return def.Files.Coding
}

// DocumentationFileScope handles identifying documentation-related files
type DocumentationFileScope struct{}

func (s *DocumentationFileScope) SelectFiles(def AgentDef) []string {
	return def.Files.Documentation
}

// ResolveFileScopes maps a configuration directive to a list of scope strategies
func ResolveFileScopes(directive string) []AssetFileScopeStrategy {
	switch directive {
	case "all":
		return []AssetFileScopeStrategy{&CodingFileScope{}, &DocumentationFileScope{}}
	case "coding", "coding-only":
		return []AssetFileScopeStrategy{&CodingFileScope{}}
	case "documentation", "documentation-only":
		return []AssetFileScopeStrategy{&DocumentationFileScope{}}
	case "none", "":
		return []AssetFileScopeStrategy{}
	default:
		// Default behavior for unknown directives matches "all"
		return []AssetFileScopeStrategy{&CodingFileScope{}, &DocumentationFileScope{}}
	}
}
