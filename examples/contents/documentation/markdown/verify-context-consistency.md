---
id: verify-context-consistency
title: Verify Context Consistency
status: adopted
description: >
  When updating documents, always verify context consistency for both the single file and the overall multi-file context.
keywords:
  - markdown
  - context
  - consistency
---

## Summary

Use a holistic approach when updating documents. Ensure that changes in one file do not introduce contradictions or breaks in logic across related files.

## Rationale

- **Coherence**: Maintaining a consistent narrative prevents confusion for both human readers and AI agents.
- **Accuracy**: Isolated updates can lead to stale or conflicting information in other parts of the documentation.

## Guidance

1.  **Single File Check**: Verify the updated file for internal consistency.
2.  **Multi-File Check**: Verify the impact of the change on the overall project context.
    - Check for references to the updated content in other files.
    - Ensure definitions and terminology are consistent globally.
