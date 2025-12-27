---
id: handle-errors-explicitly
title: Explicit Error Handling
description: >
  Never ignore errors. Handle them or return them.
keywords:
  - go
  - error-handling
  - robustness
---

## Summary
Do not use `_` to ignore errors. Always check `if err != nil`.

## Guidance
- Use automated tools like `errcheck` to find ignored errors.
- Wrap errors with context when bubbling up: `fmt.Errorf("do something: %w", err)`.

## Rationale
Silent failures are the hardest bugs to debug.
