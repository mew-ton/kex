---
type: indicator

title: Use BDD-style test naming
description: >
  When writing test descriptions, enforce BDD-style naming with "it should" or "it must".
keywords:
  - testing
  - naming
  - bdd
  - go
---

## Summary

When writing tests, especially with table-driven tests or `t.Run`, use a Behavior-Driven Development (BDD) style for naming test cases.

## Rationale

- **Readability**: The test output reads like a specification of the system.
- **Clarity**: Forces the writer to think about *behavior* rather than just "test case 1".
- **Consistency**: Standardizes how tests are described across the codebase.

## Guidance

Every `t.Run` description **MUST** start with `it should` or `it must`, followed by a clear description of the expected outcome or behavior.

## Examples

### Bad

```go
t.Run("invalid config", func(t *testing.T) { ... })
t.Run("test 1", func(t *testing.T) { ... })
t.Run("error check", func(t *testing.T) { ... })
```

### Good

```go
t.Run("it should return an error when the configuration is invalid", func(t *testing.T) { ... })
t.Run("it must fail if the required field 'id' is missing", func(t *testing.T) { ... })
```
