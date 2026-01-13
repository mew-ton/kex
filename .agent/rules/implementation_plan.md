---
trigger: always_on
---
# Implementation Plan (Ref: Kex)

**Core Rule**: When creating an `implementation_plan.md`, you MUST include a "Relevant Knowledge" section.

## Relevant Knowledge Section

**Requirement**: Before design/coding, you must search Kex for **guidelines** relevant to your task and list them here.

### Format

```markdown
## Relevant Knowledge
([Search Query used])

- [Guideline Name]: <1-line summary of constraint>
```

### Process

1.  **Analyze Task**: "What guidelines might apply?" (e.g., "go", "testing", "naming").
2.  **Search**: Use `search_documents`.
3.  **Record**: Copy the file path and specific rule into the plan.

**(This section is mandatory. Plans without it are considered incomplete.)**
