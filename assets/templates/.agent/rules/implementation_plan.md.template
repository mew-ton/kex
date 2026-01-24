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

1.  **Analyze Task**: Determine the **Scope** for the task (e.g., `["coding", "go"]`).
2.  **Search**: Execute `search_documents` with `exactScopeMatch: true` for the identified scope.
3.  **Record**: You **MUST** list **EVERY** guideline returned by the search in the "Relevant Knowledge" section.
    *   **Constraint**: Do not filter, select, or exclude any results found by the scope search. Capture all of them.

**(This section is mandatory. Plans without it are considered incomplete.)**
