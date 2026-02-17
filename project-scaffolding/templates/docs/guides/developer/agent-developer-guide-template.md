# Agent Development Guide - Agent System

## 🎯 Purpose

This guide explains how to create, modify, and maintain AI agents for this project.

**Main Entry Point:** To start any development session or interact with the assistant, you **MUST** load the `agents/init.yaml` agent.

---

## 🚀 Agent Capabilities & Workflows

Our AI agents are designed to assist you across various development tasks and provide structured interaction:

### 1. Starting a Work Session
The `init.yaml` agent acts as your main menu. From there, you can choose to:
- **Continue or start a development task (Sprints and Tasks):** This will typically lead you to `load-session.yaml`, which manages existing sessions or helps you start new development workflows based on project tasks. This is where most active development work begins.
- **Get project information (documentation, architecture, agents):** This option allows you to query the AI about project specifics, load context agents like `architecture.yaml` or `tech-stack.yaml`, or review agent documentation.

### 2. Context Loading
Agents are crucial for context management. `load-session.yaml` and `load-project-context.yaml` ensure that `generic-rules.yaml` and project-specific contexts are loaded dynamically as needed, providing the AI with the relevant information for your task.

### 3. Modular Development
Agents facilitate the Standard Module Development Workflow (SMDW) by guiding you through phases like scaffolding, domain-driven design, backend and frontend implementation, and final integration.

### 4. Code Quality & Standards Enforcement
Agents leverage `generic-rules.yaml` and `code-standards.yaml` to provide guidelines on coding conventions, testing philosophy, and quality gates, helping ensure adherence to project standards.

---

## 📁 System Structure

```
agents/
├── generic-rules.yaml          # Universal rules
├── init.yaml                   # Entry point
├── load-project-context.yaml   # Context loader
├── load-session.yaml           # Session manager
├── end-session.yaml            # Session closer (newly added to agent list)
├── align-context-agents.yaml   # Context alignment with docs
└── project/
    ├── project-context.yaml    # Project context
    ├── sprint-registry.yaml    # Sprint registry
    └── context/                # Modular contexts
        ├── README.md
        ├── architecture.yaml
        ├── bounded-contexts.yaml
        ├── code-standards.yaml
        └── tech-stack.yaml
```

### Agent Levels

1. **Universal** - Apply to any project (generic-rules.yaml)
2. **Project** - Specific to this project (project-context.yaml)
3. **Modular** - Specialized contexts (context/*)

---

## ✨ Creating a New Context Agent

### Step 1: Define Purpose

What information does the AI need that is not in existing agents?

**Example:** If you need to document external APIs → create `external-apis.yaml`

### Step 2: Create File

**Location:** `agents/project/context/descriptive-name.yaml`

**Base Template:**
```yaml
# Descriptive Agent Name
# Purpose: Brief description

---

metadata:
  filename: "descriptive-name.yaml"
  last_updated: "YYYY-MM-DD"
  purpose: "Detailed description"
  load_when: "When this context is needed"

# Specific content here
section_name:
  key: value
  items:
    - item1
    - item2
```

### Step 3: Register in README

Update `agents/project/context/README.md`:
- Add to table of available contexts
- Specify when to load
- Document purpose

### Step 4: Reference if Necessary

If it's automatic loading, add to `load-project-context.yaml`

---

## 🔧 Modifying an Existing Agent

### 1. Update Content

Edit the corresponding YAML file.

### 2. Update Metadata

**ALWAYS update:**
```yaml
metadata:
  last_updated: "2026-02-14"  # Current date (updated)
```

### 3. Maintain Consistency

If you change the structure:
- Update related documentation in `/docs`
- Verify references in other agents
- Update README if applicable

---

## ✨ Generating and Maintaining Context Agents from `docs/`

**Principle:** The `docs/` directory is the **single source of truth** for key architectural and module-specific information (e.g., ADRs, module specifications). Many context agents in `agents/project/context/` are **automatically generated or updated** from these `docs/` files.

**Mechanism:**
1.  **Source Files:** Architectural Decision Records (ADRs) in `docs/architecture/adrs/` and module specification documents in `docs/modules/` contain the detailed information.
2.  **Generation Script:** The `agents/scripts/doc_to_agent_context.py` script reads these source files, extracts structured data (e.g., Bounded Context definitions, Clean Architecture layers, testing strategies), and uses it to populate/update the corresponding agent context YAML files (e.g., `agents/project/context/bounded-contexts.yaml`, `agents/project/context/architecture.yaml`).

**Developer Workflow:**
-   **Primary Action:** When updating architectural decisions or module details, **ALWAYS edit the source Markdown files in `docs/` first.**
-   **Triggering Generation:** After updating the `docs/` files:
    -   Run `python agents/scripts/doc_to_agent_context.py` from the project root.
    -   Or, select the "Alinear agentes de contexto con la documentación (docs/)" option from `agents/init.yaml`.
-   **Avoid Manual Edits:** **DO NOT manually edit** the content of generated agent context files (e.g., `agents/project/context/bounded-contexts.yaml`, `agents/project/context/architecture.yaml`) directly. Any manual changes will be overwritten the next time the generation script is run.
-   **Extending Extraction:** If new information needs to be extracted or mapped to agent contexts, extend the `agents/scripts/doc_to_agent_context.py` script.

**Benefits:**
-   Ensures consistency between human-readable documentation (`docs/`) and AI-readable configurations (`agents/`).
-   Reduces maintenance overhead by centralizing information.
-   Prevents information duplication and potential inconsistencies.

---

## 📋 Mandatory Metadata

Every agent MUST have:

```yaml
metadata:
  filename: "file-name.yaml"
  last_updated: "YYYY-MM-DD"
  purpose: "What this agent does"
```

**Optional but recommended:**
```yaml
metadata:
  load_when: "When to load this agent"
  applies_to: "What part of the project it applies to"
  dependencies: ["other-needed-agents.yaml"]
```

---

## 🚫 Fundamental Principles

### 1. No Versioning in Names

❌ **BAD:**
- `architecture-v2.yaml`
- `code-standards_2024.yaml`
- `tech-stack-old.yaml`

✅ **GOOD:**
- `architecture.yaml`
- `code-standards.yaml`
- `tech-stack.yaml`

**Reason:** Only the current version exists. History is in Git.

### 2. Descriptive Names

- Use kebab-case: `descriptive-name.yaml`
- Be specific: `external-apis.yaml` better than `apis.yaml`
- Avoid obscure abbreviations

### 3. Clear and Concise Content

- Get straight to the point
- Use lists and tables
- Avoid redundant text
- Examples over lengthy explanations

---

## 🎨 Best Practices

### When to Create a New Agent

✅ **Create new if:**
- Information is not in existing agents
- Specific and isolated topic
- Loaded on demand (modular)

❌ **DO NOT create if:**
- Information fits into an existing agent
- It is temporary or changes frequently
- Better to document in `/docs`

### When to Modify an Existing One

- Incomplete or outdated information
- Changes in architecture/decisions
- Clarity improvements

### Sync with /docs

If you modify an agent:
- Verify that documentation in `/docs` is aligned
- Agents = information for AI
- Docs = information for humans
- Both must be consistent

---

## ✅ Checklist

Before committing changes:

- [ ] Metadata `last_updated` updated
- [ ] Content is clear and concise
- [ ] No versioning in filename
- [ ] README updated if new agent
- [ ] Documentation in `/docs` synchronized
- [ ] References in other agents verified

---

## 🆘 Frequently Asked Questions

**Should I create an agent or document in /docs?**
- Agent: Information the AI constantly needs
- Docs: Information for humans, tutorials, long explanations

**When to update last_updated?**
- Every time you modify the agent's content

**Can I have agents in other directories?**
- No. Fixed structure: `/agents` (universal) and `/agents/project` (project)

**What to do with obsolete agents?**
- Update content (no versioning)
- If no longer applicable: delete and document in commit

---

**For using agents:** See `AGENTS.md`  
**Full reference:** See `AGENTS.md` (or `user-guide.md` for users, if applicable)
