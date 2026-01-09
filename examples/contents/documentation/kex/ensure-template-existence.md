---
id: ensure-template-existence
title: Ensure Template Existence for Generators
description: Guidelines for ensuring template files exist physically when using fs.WalkDir for generation strategies.
status: adopted
keywords: [generator, templates, golang, walkdir, troubleshooting]
---

# Ensure Template Existence for Generators

## Context
When implementing file generators that rely on iterating over a template directory (e.g., using `fs.WalkDir`), a common pitfall is to implement dynamic generation logic for a file but forget to place a placeholder for that file in the template source.

## The Problem
File system walkers typically only visit files that physically exist in the source directory. If you implement logic like:

```go
if filepath.Base(path) == "DYNAMIC_FILE.md" {
    // Generate content dynamically
}
```

But "DYNAMIC_FILE.md" does not exist in the source `templates/` folder, the walker will never trigger this condition.

## The Solution
Always ensure that every file you intend to generate—even if its content is wholly dynamic—has a corresponding physical file (placeholder) in the template directory.

### Example
If `AGENTS.md` is generated dynamically:
1.  **Code**: Implement the generation logic in the walker loop.
2.  **Template**: Create an empty or comment-only `assets/templates/AGENTS.md`.

```markdown
<!-- Placeholder for dynamic generation -->
```

This ensures the walker visits the path and triggers your generation logic.
