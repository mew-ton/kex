---
id: avoid-unrequested-scope-expansion
title: Avoid Unrequested Scope Expansion
description: Guidelines for staying focused on the immediate problem and avoiding implementing unrelated features.
status: adopted
keywords: [best-practice, scope-creep, requirement-analysis, YAGNI]
---

# Avoid Unrequested Scope Expansion

## Context
When assigned a task to fix a specific issue (e.g., "Fix keyword matching collision"), there is a temptation to "improve" other aspects of the system (e.g., "Add title searching for better relevance") without explicit request or approval.

## The Principle
**"Solve the problem in front of you, not the problem you imagine might exist."**

1.  **Strict Adherence**: If the task is to fix a bug in Feature A, do not implement Feature B, even if it seems helpful.
2.  **Cost of Unrequested Features**:
    *   **Distraction**: It shifts focus away from the core issue (e.g., debugging why the original logic failed).
    *   **Maintenance**: New code ("Title Search") increases the surface area for bugs and tests.
    *   **Review Overhead**: The reviewer has to review unrequested logic, confusing the PR intent.

## Implementation Guideline
*   **Verify Requirements**: If you think a new feature (like Title/Description search) is needed to solve the issue, **ask first**. do not implement it unilaterally.
*   **YAGNI (You Aren't Gonna Need It)**: Do not implement broad search capabilities if the requirement is narrow precision filtering.
*   **Revert Unjustified Changes**: If you realize you have added logic not traceable to the original requirement, be prepared to revert it immediately.
