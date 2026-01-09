---
id: avoid-unrequested-scope-expansion
title: Avoid Unrequested Scope Expansion
description: Guidelines for staying focused on the immediate problem and avoiding implementing unrelated features.
status: adopted
keywords: [best-practice, scope-creep, requirement-analysis, YAGNI]
---

# Avoid Unrequested Scope Expansion

## Context
When assigned a specific task (e.g., "Fix a bug in feature X"), there is often a temptation to "improve" surrounding code or add related features (e.g., "Add a new filter capability") that were not requested.

## The Principle
**"Solve the problem in front of you, not the problem you imagine might exist."**

1.  **Strict Adherence**: If the task is to fix Feature A, do not implement Feature B, even if it seems helpful or "easy to add while I'm here."
2.  **Cost of Unrequested Features**:
    *   **Distraction**: It shifts focus away from the core issue.
    *   **Maintenance**: New code increases the surface area for bugs and tests without a verified user need.
    *   **Review Overhead**: The reviewer has to review unrequested logic, confusing the PR intent.

## Implementation Guideline
*   **Verify Requirements**: If you believe a new feature is necessary to solve the assigned task, **propose it explicitly** and wait for approval before implementing.
*   **YAGNI (You Aren't Gonna Need It)**: Do not implement features "just in case" they are useful later.
*   **Revert Unjustified Changes**: If you have added logic not traceable to the original requirement, revert it.
