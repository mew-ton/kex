---
id: write-concise-content
title: Write concise content
description: >
  Keep guideline text brief and direct. Code examples can be long if necessary for context.
keywords:
  - content
  - writing
  - style
  - conciseness
  - concise
  - example
---

## Summary
Keep explanations brief and direct on "what" and "why". Use code examples to show "how".

## Rationale
- **AI Efficiency**: Concise text consumes fewer tokens and reduces hallucination risk.
- **Human Efficiency**: Engineers prefer reading code over long prose.

## Guidance
1.  **Text**: Avoid fluff. Use bullet points and numbered lists.
2.  **Code Examples**: 
    - Do **not** truncate essential context if it harms understanding.
    - It is acceptable for code examples to be long if they demonstrate a complex pattern or a necessary "before/after" transformation.
    - Always provide a **Bad** (anti-pattern) and **Good** (corrected) example.

## Examples

### Bad (Text)
"In this section, we will discuss the importance of keeping your functions short. It is generally believed by many experts in the field of software engineering that..." (Too verbose)

### Good (Text)
"Keep functions short. Long functions are hard to test and maintain." (Direct)
