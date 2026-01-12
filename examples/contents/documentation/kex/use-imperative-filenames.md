---

title: Use imperative filenames
description: When naming guideline files, create imperative filenames that act as actionable prompts.
keywords:
  - filename
  - convention
  - instruction
  - imperative
  - file
---

## Summary
Name guideline files using the imperative mood (verb-adjective-object) to act as direct instructions.

## Rationale
- **Actionability**: The filename itself tells the AI or human exactly what to do.
- **Discoverability**: Verbs like "use", "avoid", "handle" are strong search terms for actions.

## Guidance
1.  Start with a **Verb** (e.g., `Keep`, `Use`, `Avoid`, `Handle`).
2.  Follow with the **Object** or **Condition**.
3.  Ensure the `id` in frontmatter matches the filename (minus extension).

## Examples

### Bad
- `functions.md` (Noun only)
- `error_handling_guide.md` (Topic based)
- `about_naming.md` (Vague)

### Good
- `keep-functions-short.md`
- `handle-errors-explicitly.md`
- `use-consistent-naming.md`
