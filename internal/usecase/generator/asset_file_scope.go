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
func ResolveFileScopes(scopes []string) []AssetFileScopeStrategy {
	var strategies []AssetFileScopeStrategy

	for _, scope := range scopes {
		switch scope {
		case "coding":
			strategies = append(strategies, &CodingFileScope{})
		case "documentation":
			strategies = append(strategies, &DocumentationFileScope{})
		case "all":
			return []AssetFileScopeStrategy{&CodingFileScope{}, &DocumentationFileScope{}}
		}
	}
	return strategies
}
