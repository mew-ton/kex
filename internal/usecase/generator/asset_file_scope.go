package generator

// AssetFileScopeStrategy defines which files to select from an agent definition
type AssetFileScopeStrategy interface {
	SelectFiles(def AgentDef) []string
}

// CommonFileScope handles identifying common-related files
type CommonFileScope struct{}

func (s *CommonFileScope) SelectFiles(def AgentDef) []string {
	return def.Files.Common
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
		case "documentation":
			strategies = append(strategies, &CommonFileScope{})
			strategies = append(strategies, &DocumentationFileScope{})
		case "all":
			return []AssetFileScopeStrategy{&CommonFileScope{}, &DocumentationFileScope{}}
		}
	}
	return strategies
}
