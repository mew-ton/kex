---
id: clarify-search-intent
title: Clarify Search Intent (Scope vs Keyword)
description: Guidelines for designing search interfaces that distinguish between container scopes and tag keywords.
status: adopted
keywords: [search, design, indexing, scope, keyword]
---

# Clarify Search Intent (Scope vs Keyword)

## Context
When implementing search functionality for a knowledge base, it is crucial to distinguish between "finding a file within a container" (Scope) and "finding a file matched by a tag" (Keyword).

## The Principle
*   **Scope Match**: Matches the directory structure or hierarchy. A file is *in* a scope.
*   **Keyword Match**: Matches the content or explicit tags. A file *is about* a keyword.
*   **Broad Search**: Combining both (OR logic) is useful for general exploration.
*   **Strict Filtering**: Refactoring tools often need **Strict Scope Matching** ("Give me all domain rules", regardless of whether they mention "domain" in the text).

## Implementation Guideline
Do not rely on implicit text matching for strict scoping.
1.  **Separate Indices**: Maintain a clean `ScopeIndex` separate from the broad text `Index`.
2.  **Explicit Flags**: Provide users (and agents) with explicit control (e.g., `exactScopeMatch: true`) when they need to filter strictly by location/container.

## Anti-Pattern
Merging Scope names into the general text index without a way to filter them out later leads to "noisy" results where a document is found because it *mentions* a folder name, not because it resides in it.
