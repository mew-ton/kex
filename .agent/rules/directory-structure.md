---
trigger: always_on
---
# Directory Structure & Context Rules

To avoid confusion between Kex's *source code* and Kex's *runtime configuration* (Self-Hosting), understand the following directory roles:

## Root Directories

- **`.agent/` (or `.antigravity/`)**:
    - **Purpose**: Contains the **Active Rules** and configuration for **THIS** repository (Kex itself).
    - **Usage**: The AI Agent reads these files to understand how to develop Kex.
    - **Note**: Rules here (like Self-Hosting guidelines) apply *immediately* to the current session.

- **`assets/`**:
    - **Purpose**: Contains **Static Resources** and **Templates** embedded into the Kex binary.
    - **Usage**: Source code for `kex init` and `kex update`.
    - **Critical**: Modifying files here affects *future users* of Kex, but does NOT affect the current agent's rules unless explicitly copied to `.agent/`.
    - **Do NOT** add repo-specific rules (like "Self-Hosting") here unless they should apply to *everyone* using Kex.

- **`bin/`**:
    - **Purpose**: Output directory for `make build`.
    - **Usage**: The active `kex` command used by the agent (`mcp-kex`) runs from here.
    - **Warning**: Overwriting this file destroys the tool currently supporting the agent.

## Workflow Implication

When asked to "Update Rules":
1.  **For Kex Developers (Us)**: Update `.agent/rules/`.
2.  **For Kex Users (Everyone)**: Update `assets/templates/`.

Always verify which context (Local vs. Global) the request targets.
