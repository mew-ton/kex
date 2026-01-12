---
title: Prevent Dead Code and Asset Mismatch
description: When modifying code or assets, prevent dead code and ensure synchronization between application logic and static resources.
keywords:
  - best-practice
  - maintenance
  - cleanup
  - assets
  - dead-code
---

## Summary
Do not leave dead code, unnecessary logic, or unused assets. Ensure perfect synchronization between application logic and the static assets (templates, images, etc.) it processes.

## Rationale
- Use of nonexistent assets leads to phantom logic that confuses developers.
- Unused assets bloat the repository and reduce signal-to-noise ratio.
- Coupled logic (like generators) requires both the code and the target asset to functional correctly.

## Guidance
1.  **Logic without Assets is Dead Code**: If you write code to process a file (e.g., specific handling for `AGENTS.md`), ensure that file exists. If the asset is removed, remove the logic.
2.  **Assets without Logic are Bloat**: Delete files that are never read, processed, or deployed.
3.  **Search References**: When deleting assets, search for path strings to clean up handling logic.
4.  **Fail Fast**: Prefer tests that fail when an expected asset is missing over logic that silently skips it.

## Examples

### Bad
Writing a generator loop that checks for `AGENTS.md` when no such file exists in the templates folder. The code runs but does nothing, confusing future maintainers.

### Good
If logic exists to process `AGENTS.md`, a placeholder `assets/templates/AGENTS.md` exists to verify that logic is active and tested.
