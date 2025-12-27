# Document Librarian MCP

Final Specification v1.0
Reference Implementation: kex

## 1. Purpose

Document Librarian MCP is a local MCP (Model Context Protocol) that provides
structured access to project-specific documents such as:

- Coding guidelines
- Readable Code principles
- Design patterns
- Architecture and design conventions

It acts as a librarian, not a judge.

The system does not enforce rules, assign severity, or make decisions.
Instead, it enables humans, AI coding agents, and CI workflows to
locate and reference the correct documents at the moment a decision is needed.

## 2. Core Principles

- Documents are queried, not memorized
- Knowledge lives inside the repository
- AI is a scribe, humans are decision makers
- Structure is strict, meaning is human-owned
- No central aggregation
- No enforcement, no policy engine

## 3. Non-Goals

This project explicitly does not:

- Lint code or analyze ASTs
- Assign or manage rule severity
- Enforce compliance
- Aggregate public best practices
- Quote or redistribute books or paid sources
- Act as a documentation site or portal
- Let AI decide which rules matter

## 4. Distribution Model

- Local CLI
- Single binary (Go)
- Executed from the repository root
- MCP server starts only when needed (non-resident)

## 5. CLI Interface (Reference Implementation: kex)

### 5.1 `kex init`

Initializes the repository for use with Document Librarian MCP.

- Creates configuration files
- Creates the document root directory (default: contents/)
- Optionally scaffolds example structure
- Idempotent (safe to run multiple times)

### 5.2 `kex start`

Starts the MCP server for document access.

- Loads documents
- Builds the keyword index
- Exposes MCP endpoints
- Fails to start if structure or schema is invalid

### 5.3 `kex check`

Validator for CI and local verification.

- Validates document schema
- Validates path and ID consistency
- Fails if any document has status: draft
- Returns exit codes for CI integration

## 6. Document Format

### 6.1 General Rules

- One document per file
- Markdown (UTF-8)
- Filename must match id.md
- Body content must be English (AI-first)

### 6.2 Required Frontmatter

```yaml
---
id: avoid-magic-numbers
title: Avoid magic numbers
description: >
  Do not use unexplained numeric literals directly in code.
keywords:
  - readability
  - magic-number
  - maintainability
---
```

#### Required Fields

| Field	| Description |
| --- | --- |
| id	| Unique within the repository; must match filename |
| title	| Short, human- and AI-readable guideline |
| description	| Summary returned by MCP |
| keywords	| Primary index keys |

#### Optional Fields

```yaml
sources:
  - name: Biome
    url: https://biomejs.dev/linter/rules/no-explicit-any/
  - name: ESLint
    url: https://eslint.org/docs/latest/rules/no-explicit-any
status: draft | adopted
```

- `sources`: List of provenance references. Must link directly to the **rule definition**, not the tool homepage.
- If `status: draft` exists → omit document from index
- If omitted, status is treated as `adopted`

No severity field exists by design.

## 6.3 Body (Recommended Structure)

```markdown
## Summary
## Rationale
## Guidance
## Examples
```

- Content correctness is not validated
- Body text is not indexed

## 7. Scope and Directory Structure

### 7.1 Root Directory

- Default: contents/
- Configurable (single directory only)

### 7.2 Universal Rules (Special Handling)

```plaintext
contents/universal/
  <no-language-guideline>.md
  <language>/
    <language-specific-guideline>.md
```

#### Semantics

- `universal/*.md`
  Language- and environment-independent principles
  (Readable Code, design principles, patterns)
- universal/<language>/*.md
  Language-specific rules that apply universally
  (language semantics, not runtime-specific)

#### Constraints

- Maximum depth: 2 levels
- <language> must be a known language identifier
- No platform or framework under universal

### 7.3 General Structure

```plaintext
<root>/<domain>/<platform>/<technology>/<id>.md
```

- Maximum depth: 3 levels
- Language names appear only as technology
- Structure represents context hints, not strict taxonomy

#### Layer Roles

|Layer|Meaning|Examples|
|---|---|---|
|domain|Application context|frontend, backend|
|platform|Execution environment|web, nodejs|
|technology|Language or framework|javascript, typescript, react


## 8. Keyword Index

- Primary retrieval mechanism
- Indexed fields:
  - `keywords`
  - `id`
  - `title`
  - `description`
  - Path-derived scope
- Relevance is advisory, not authoritative

## 9. MCP Responsibilities

Document Librarian MCP provides:

- Keyword-based document lookup
- Scope-aware relevance hints
- Access to document content
- Document location references

It does not:

- Decide applicability
- Rank importance
- Enforce compliance

All interpretation belongs to the caller.

## 10. Relationship with Linter and CI

- Linters:
  - Decide which rules are enabled
  - Decide enforcement severity
- Document Librarian MCP:
  - Explains what a rule means
  - Provides rationale and guidance

These systems are complementary and independent.

## 11. AI Usage Model (Normative)

#### AI must use Document Librarian MCP when:

- Writing or modifying code
- Making structural or design decisions
- Explaining linter errors
- Justifying architectural choices

#### AI behavior rules:

1. Query MCP when a decision requires project context
2. Treat returned documents as project assumptions
3. If MCP returns nothing, assume no explicit guideline exists
4. Do not substitute general best practices
5. Do not assign importance or enforcement strength

AI is a consumer of documents, not their authority.

## 12. AI Configuration Guidelines

To ensure consistent behavior, AI agents must be configured with explicit roles and system instructions. 
This distinction prevents "hallucinations" of rules and ensures that the AI acts as a consumer of documentation rather than an inventor of it.

### 12.1 Role Separation

We define two distinct AI roles. You must decide which role the AI agent is playing before providing instructions.

#### A. Consumer AI (Coding Agent)
- **Goal**: Write or modify code to meet user requirements.
- **Permission**: READ-ONLY access to documents.
- **Behavior**:
    1.  **Search first**: Queries MCP before writing code.
    2.  **Compliance**: Follows found documents as project constraints.
    3.  **No Invention**: If no document is found, acts based on general knowledge but DOES NOT invent project rules.

#### B. Maintainer AI (Librarian Agent)
- **Goal**: Assist humans in writing or updating guideline documents.
- **Permission**: WRITE access to `contents/` directory.
- **Behavior**:
    1.  **Content Focus**: Focuses on writing clear, concise, and helpful guidelines.
    2.  **Taxonomy**: Places new files in the correct `universal/` or `domain/` structure.
    3.  **Delegation**: Relies on `kex check` for strict schema validation.

### 12.2 Configuration for Consumer AI (System Prompt)

Add this instruction to the Coding Agent's system prompt (e.g., `.cursorrules` or `claude_desktop_config.json` description):

```markdown
# Tool Usage Guidelines (Document Librarian MCP)
You have access to a tool named `kex` that indexes project guidelines.

**When to use:**
- ALWAYS query `kex` before writing new code or refactoring.
- Query when facing architectural decisions (e.g., "naming convention", "folder structure").

**Rules:**
1. **Fact Check**: Treat `kex` results as the source of truth for this project.
2. **Silence**: If `kex` returns nothing, assume no specific project rule exists.
3. **Do Not Judge**: Do not criticize the rules; just apply them.
```

### 12.3 Configuration for Maintainer AI (System Prompt)

Add this instruction when the AI is tasked with creating or editing documentation:

```markdown
# Role
You are a Librarian AI responsible for assisting in writing project documentation.

**Goal:**
Write clear, rationale-backed guidelines.

**Workflow:**
1. **Draft**: Create the document with required Frontmatter.
2. **Verify**: Use `kex check` to validate the schema and paths.
3. **Refine**: If `kex check` fails, fix the issues.
```

## 13. Final Definition

Document Librarian MCP is a local MCP that organizes and serves
project-specific documents, acting as a librarian rather than a judge.

The reference implementation, kex, provides keyword-indexed access
for coding agents and CI workflows.

## 14. Status

- Design: Final
- Scope: Locked
- Next phase: Implementation

## 15. Sample Documents and Code Usage

This section demonstrates how Document Librarian MCP is used in practice,
including a real guideline document and how an AI coding agent consumes it.

### 14.1 Sample Guideline Document

#### Path

```plaintext
contents/universal/avoid-magic-numbers.md
```

#### Document

````markdown
---
id: avoid-magic-numbers
title: Avoid magic numbers
description: >
  Do not use unexplained numeric literals directly in code.
keywords:
  - readability
  - magic-number
  - maintainability
---

## Summary
Avoid embedding raw numeric literals whose meaning is not immediately clear.

## Rationale
Magic numbers make code harder to understand and maintain.
They obscure intent and increase the risk of errors when requirements change.

## Guidance
Replace magic numbers with named constants or enums that clearly express intent.

## Examples

### Bad

```ts
function isAdult(age: number): boolean {
  return age >= 20;
}
```

### Good

```ts
const ADULT_AGE_THRESHOLD = 20;

function isAdult(age: number): boolean {
  return age >= ADULT_AGE_THRESHOLD;
}
```
````

This document:

- Is language-agnostic
- Lives under `universal/`
- Provides rationale and guidance without enforcing severity

### 14.2 Language-Specific Universal Rule

#### Path

```txt
contents/universal/typescript/avoid-any.md
```

#### Document

````markdown
---
id: avoid-any
title: Avoid using the any type
description: >
  Avoid the any type to preserve type safety in TypeScript.
keywords:
  - typescript
  - type-safety
  - any
---

## Summary
The any type disables type checking and should be avoided.

## Rationale
Using any bypasses TypeScript’s type system and hides potential bugs.

## Guidance
Use unknown, generics, or proper union types instead of any.

## Examples

### Bad
```ts
function parse(input: any): any {
  return JSON.parse(input);
}
```

### Good
```ts
function parse(input: string): unknown {
  return JSON.parse(input);
}
```
````

## 14.3 AI Usage Example (Normative)

### Scenario: Code Generation with Project Context

#### Situation

An AI coding agent is about to write TypeScript code
that includes numeric thresholds and type parsing.

#### Step 1: AI Identifies Decision Keywords

```text
keywords:
- readability
- magic-number
- typescript
- type-safety
```

#### Step 2: AI Queries Document Librarian MCP

```json
{
  "query": {
    "keywords": [
      "readability",
      "magic-number",
      "typescript",
      "type-safety"
    ],
    "context": {
      "language": "typescript"
    }
  }
}
  ```

#### Step 3: MCP Returns Relevant Documents

```json
{
  "documents": [
    {
      "id": "avoid-magic-numbers",
      "title": "Avoid magic numbers",
      "description": "Do not use unexplained numeric literals directly in code.",
      "path": "contents/universal/avoid-magic-numbers.md"
    },
    {
      "id": "avoid-any",
      "title": "Avoid using the any type",
      "description": "Avoid the any type to preserve type safety in TypeScript.",
      "path": "contents/universal/typescript/avoid-any.md"
    }
  ]
}
```

#### Step 4: AI Generates Code Using Documents as Constraints

```ts
const ADULT_AGE_THRESHOLD = 20;

function parseUserAge(input: string): unknown {
  const age = Number.parseInt(input, 10);
  return age >= ADULT_AGE_THRESHOLD;
}
```

#### Important Notes

- The AI does not decide whether rules are mandatory
- The AI does not assign severity
- The documents act as project assumptions, not laws

## 15.4 Linter Integration Example (Complementary)

```json
// biome.json
{
  "linter": {
    "rules": {
      "suspicious": {
        "noExplicitAny": "error"
      }
    }
  }
}
```

- The linter enforces noExplicitAny
- Document Librarian MCP explains why
- Responsibility is clearly separated

## 16. What This Example Demonstrates

- Documents are small, focused, and human-owned
- Keywords are the only retrieval contract
- AI uses MCP dynamically, not as memory
- Enforcement is external and optional
- The system scales without centralization

## 17. Final Confirmation

At this point:

- The concept is explicit
- The AI interaction model is concrete
- The document format is proven
- The scope rules are exercised
- The separation of concerns is visible in code

This specification is ready for implementation and handoff.