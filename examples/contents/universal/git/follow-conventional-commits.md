---
id: follow-conventional-commits
title: Commit Messages (Conventional Commits)
description: >
  Follow Conventional Commits for git history.
keywords:
  - git
  - naming
  - history
---

## Summary
Use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages.

## Guidance
Format: `<type>(<scope>): <subject>`

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Changes that do not affect meaning of code (white-space, formatting, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries

## Rationale
Enables automated changelog generation and easier history parsing.
