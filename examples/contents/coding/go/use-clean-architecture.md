---
id: use-clean-architecture
title: Use Clean Architecture
description: >
  Organize Go projects using Clean Architecture layers to separate concerns and ensure testability.
keywords:
  - architecture
  - clean-architecture
  - go
  - layering
---

## Summary

Adopt Clean Architecture to organize code into concentric layers, ensuring that dependencies only point inward. This separates business logic from external concerns like the CLI, Database, or API.

## Rationale

-   **Separation of Concerns**: Business logic is independent of UI, Database, and Frameworks.
-   **Testability**: Business rules can be tested without external elements.
-   **Maintainability**: Changes in external libraries or tools do not affect the core logic.

## Guidance

Organize your project into the following directory structure:

### 1. Domain (`internal/domain`)
Contains **Entities** (Enterprise Business Rules) and **Repository Interfaces**.
-   **Dependencies**: None.
-   **Content**: Pure Go structs and interfaces.

### 2. Use Cases (`internal/usecase`)
Contains **Application Business Rules**.
-   **Dependencies**: `domain`.
-   **Content**: Interactors that orchestrate the flow of data to and from the domain entities.

### 3. Interfaces (`internal/interfaces`)
Contains **Interface Adapters**.
-   **Dependencies**: `usecase`.
-   **Content**: Entry points and converters.
    -   `cli`: CLI commands.
    -   `http` / `mcp`: Server handlers.

### 4. Infrastructure (`internal/infrastructure`)
Contains **Frameworks & Drivers**.
-   **Dependencies**: `domain` (implements interfaces defined there).
-   **Content**: Database implementations, File System access, Configuration loading.

## dependency Rule

Source code dependencies must point only **inward**, toward higher-level policies. Inner circles must know nothing about outer circles.
