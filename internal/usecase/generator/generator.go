package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Generator struct {
	Templates embed.FS
}

func New(templates embed.FS) *Generator {
	return &Generator{Templates: templates}
}

func (g *Generator) Generate(cwd string) error {
	// Extract templates mirroring the structure in assets/templates
	err := fs.WalkDir(g.Templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(cwd, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(g.Templates, path)
		if err != nil {
			return err
		}

		// Special handling for AGENT.md to append guidelines if file exists
		if filepath.Base(relPath) == "AGENT.md" {
			return g.handleAgentMD(targetPath, data)
		}

		// Don't overwrite if exists
		if _, err := os.Stat(targetPath); err == nil {
			return nil
		}

		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		return fmt.Errorf("failed to extract templates: %w", err)
	}

	return nil
}

func (g *Generator) handleAgentMD(targetPath string, templateData []byte) error {
	// If file doesn't exist, create it
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return os.WriteFile(targetPath, templateData, 0644)
	}

	// Read existing file
	content, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}

	// Check if Kex guidelines are already present (simple check)
	if string(content) == string(templateData) {
		return nil
	}

	// Check for a characteristic string to avoid partial duplication
	// This is a heuristic; user might have modified it.
	// But simply appending if "Kex" isn't mentioned might be safer?
	// For now, let's append nicely with a separator if not exact match.
	// Wait, if I append blindly, I might duplicate.
	// Let's check for "Design Phase Guidelines" or similar unique headers from the template.
	// The template has "## 1. Design Phase Guidelines".
	// If that exists, we assume it's there.
	
	/*
		Proposed logic:
		If "Ref: Kex" or "Kex" guidelines seem missing, append.
		Let's look for "Kex" and "Design Phase Guidelines".
	*/
	
	// Actually, just checking if the specific "Project Guidelines (Ref: Kex)" header exists might be enough.
	// But users might change the header.
	// Let's just append if we don't find "Design Phase Guidelines" AND "Implementation Phase Guidelines".
	
	// Simplified merge: Just append with a note if it seems completely different.
	// But let's stick to the plan: "Append content if not present".
	
	// Let's just write a simple check for now.
	// If the file contains "Search for design documents", we assume it has the rules.
	
	/*
	   Actually, let's just append it with a newline if it's not the exact same content.
	   But that risks duplication.
	   Let's append ONLY IF "Search for design documents" is NOT present.
	*/
	
	searchStr := "Search for design documents"
	// simplified check
	for i := 0; i < len(content)-len(searchStr); i++ {
		if string(content[i:i+len(searchStr)]) == searchStr {
			return nil // Already present
		}
	}

	// Append
	f, err := os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n\n"); err != nil {
		return err
	}
	if _, err := f.Write(templateData); err != nil {
		return err
	}

	return nil
}
