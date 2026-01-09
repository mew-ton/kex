---
id: prevent-dead-code-asset-mismatch
title: Prevent Dead Code and Asset Mismatch
description: Guidelines for maintaining synchronization between application logic and static assets.
status: adopted
keywords: [best-practice, maintenance, cleanup, assets, dead-code]
---

# Prevent Dead Code and Asset Mismatch

## Context
Codebases often contain logic coupled with static assets (templates, config files, images). A disconnect between the two leads to technical debt and bugs.

## The Principle
**"Do not leave dead code, unnecessary logic, or unused assets."**

1.  **Logic without Assets is Dead Code**:
    If you write code to process a file (e.g., specific handling for `AGENTS.md` in a generator), but that file does not exist, the code is unreachable/dead. It clutters the codebase and misleads developers.
    *   *Action*: If the asset is needed, ensure it exists (even as a placeholder). If not, delete the logic.

2.  **Assets without Logic are Bloat**:
    Files sitting in your repository that are never read, processed, or deployed should be removed.
    *   *Action*: Delete unused files to reduce noise.

## Implementation Guideline
*   **When Deleting Assets**: Always search the codebase for references (filenames, paths). Remove the corresponding handling logic.
*   **When Adding Logic**: Verify the target asset exists. If the logic depends on iteration (e.g., `fs.WalkDir`), ensure the asset is present in the walked directory to trigger the code.
*   **Testing**: Add tests that fail if a required asset is missing, rather than failing silently (which hides the dead code).
