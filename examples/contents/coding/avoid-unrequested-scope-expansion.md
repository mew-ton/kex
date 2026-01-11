---
title: Avoid Unrequested Scope Expansion
description: Guidelines for staying focused on the immediate problem and avoiding implementing unrelated features.
keywords:
  - best-practice
  - scope-creep
  - requirement-analysis
  - YAGNI
---

## Summary
Solve the problem in front of you, not the problem you imagine might exist. Avoid implementing features or improvements that were not explicitly requested or approved.

## Rationale
- **Distraction**: Shifts focus away from the core issue, delaying the actual fix.
- **Maintenance Cost**: New features increase surface area for bugs and ownership costs.
- **Review Complexity**: Confuses the intent of the PR and wastes reviewer time on unrequested logic.

## Guidance
1.  **Strict Adherence**: If the task is to fix Feature A, do not touch Feature B.
2.  **Propose First**: If a specific feature seems necessary to solve the task, propose it explicitly and wait for approval.
3.  **YAGNI**: Do not implement code "just in case" it's useful later.
4.  **Revert**: If you find yourself adding logic not traceable to a requirement, revert it.

## Examples

### Bad
Task: "Fix keyword matching collision."
Action: "While I'm here, I'll also add title search because it seems easy." -> **Rejected**.

### Good
Task: "Fix keyword matching collision."
Action: "I noticed title search could be useful. I will open a separate issue or discussion to propose it, but for now I will strictly fix the collision bug." -> **Approved**.
