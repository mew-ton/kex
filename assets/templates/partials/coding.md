## Design & Implementation Phase

**Core Rule**: Before proposing changes or writing code, understand existing design and coding standards.

1.  **Search for guidelines** using `search_documents`.
    *   **Keywords**: "architecture", "design", "coding", "style", "naming", "testing", "strategy".
    *   **Critical:** If `kex` tools are unavailable or fail, **STOP** and report this to the user.
2.  **Read relevant docs** using `read_document`.
3.  **Validate your plan** (or existing code) against principles found.

## Refactoring & Code Audit Phase

**Core Rule**: When asked to review code, find bugs, or check for guideline violations, you MUST NOT rely solely on your internal training.

1.  **Identify the context** (Language, Framework, Component, Function).
2.  **Search for specific rules** using `search_documents`.
    *   **Keywords**: "anti-pattern", "best practice", "forbidden", "required", "error handling", [Context Specific Keywords].
    *   **Tip**: When focusing on a specific layer or component (e.g., `domain`, `infrastructure`), use `search_documents` with `exactScopeMatch: true` and the component name as the keyword to see all rules for that scope.
3.  **Verify the code** against the *explicit* rules retrieved from Kex.
    *   *Do not assume* standard conventions apply if Kex has specific overrides.


