---
id: no-use-any
title: Avoid using the any type
description: >
  Enforce strict type safety by avoiding 'any'.
keywords:
  - typescript
  - type-safety
  - lint-rule
  - no-any
sources:
  - name: Biome (noExplicitAny)
    url: https://biomejs.dev/linter/rules/no-explicit-any/
  - name: TypeScript (noImplicitAny)
    url: https://www.typescriptlang.org/tsconfig/#noImplicitAny
---

## Summary
Do not use `any`. This includes both explicit `any` usage and implicit `any` where type inference fails. Use `unknown`, Generics, or specific types instead.

## Rationale
Using `any` defeats the purpose of TypeScript.
- **Explicit Any**: Silences the compiler intentionally, hiding potential bugs.
- **Implicit Any**: Occurs when TypeScript cannot infer a type, falling back to dynamic typing which is unsafe.
