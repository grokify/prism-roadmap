# Merging Modular Requirement Files

The `splan merge` command provides a powerful way to manage large, complex requirements documents (like PRDs, MRDs, etc.) by breaking them into smaller, more manageable modular files. This document explains when and how to use this feature.

## When to Use `merge`

Managing a single, large JSON file for a complete requirements document can be cumbersome. It often leads to merge conflicts, difficulty in collaboration, and challenges for AI assistants working with large files.

The `merge` command is ideal when you want to:

- **Promote Collaboration**: Different teams or individuals can work on separate sections of the document simultaneously (e.g., Product Managers on user stories, Engineers on architecture).
- **Improve Readability**: Smaller, focused files are easier to read and maintain.
- **Reduce Merge Conflicts**: Changes are isolated to specific files, minimizing the chance of conflicting edits in a monolithic file.
- **Enable Reusability**: Common sections, like personas, can potentially be reused across different documents.
- **Support AI Assistants**: Large PRD files can exceed AI context windows or make updates error-prone. Modular files allow AI assistants to focus on specific sections.

The workflow involves two key activities:

1. **Splitting**: A large document is broken down into logical, self-contained JSON files.
2. **Merging**: The modular files are combined back into a single, complete document for validation, rendering, or distribution.

## File Size Guidelines for AI Assistants

AI assistants have context window limits and work most reliably with smaller, focused files:

| File Size | Recommendation |
|-----------|----------------|
| < 500 lines | Single file is fine |
| 500-1000 lines | Consider splitting if frequent updates needed |
| > 1000 lines | Split into modular files |
| > 2000 lines | Strongly recommend splitting |

**Signs you should split:**

- AI assistant makes errors when updating the file
- AI assistant "loses" content during updates
- Frequent merge conflicts in version control
- Difficulty finding specific sections to update

## Naming Conventions

### Topic-Based Naming (Recommended)

Use descriptive topic names based on PRD sections:

```
prd-metadata.json
prd-executive-summary.json
prd-objectives.json
prd-personas.json
prd-user-stories.json
prd-requirements.json
prd-roadmap.json
prd-architecture.json
prd-risks.json
```

**Advantages:**

- Easy to locate specific content
- Self-documenting file structure
- Supports targeted updates by AI assistants
- Clear ownership (e.g., "update prd-personas.json")

### Standard Topic Names

| Topic | Contents | Typical Size |
|-------|----------|--------------|
| `metadata` | ID, title, version, authors, status | Small |
| `executive-summary` | Problem, solution, outcomes | Small |
| `objectives` | OKRs, business objectives, goals | Small-Medium |
| `personas` | User personas with goals/pain points | Medium |
| `user-stories` | User stories with acceptance criteria | Medium-Large |
| `requirements` | Functional and non-functional requirements | Large |
| `roadmap` | Phases, deliverables, milestones | Medium |
| `architecture` | Technical architecture, services, APIs | Medium-Large |
| `risks` | Risk assessment and mitigations | Small-Medium |
| `market` | Market analysis, alternatives | Medium |
| `solution` | Solution options and rationale | Medium |
| `decisions` | Decision records | Small-Medium |
| `success-metrics` | North star, supporting, guardrail metrics | Small |

### Hybrid Naming for Large Topics

When a single topic grows too large, use topic + number:

```
prd-requirements-001-auth.json
prd-requirements-002-billing.json
prd-requirements-003-reporting.json
```

Or by category:

```
prd-requirements-functional.json
prd-requirements-nonfunctional.json
```

### Numbered Naming (Simple Alternative)

For simpler cases or when topics aren't clearly defined:

```
prd-part-001.json
prd-part-002.json
prd-part-003.json
```

**Note:** Numbered naming is simpler but makes it harder to locate specific content.

## How to Use `merge`

The command recursively merges two or more JSON files into a single output file. The merge is "deep," meaning nested objects are merged intelligently. If the same key exists in multiple files, the value from the *last* file specified in the command will take precedence.

### Syntax

```bash
splan merge -o <output-file.json> <input-file-1.json> <input-file-2.json> ...
```

- `-o, --output <file>`: Specifies the path for the final merged JSON file.
- `<input-file-n.json>`: A space-separated list of the modular JSON files to merge. The order is important.

### Example

Using the example files located in the `examples/merge/` directory, we can construct a complete Product Requirements Document.

The modular files are:

- `prd-metadata.json`
- `prd-objectives.json`
- `prd-personas.json`
- `prd-user-stories.json`
- `prd-requirements.json`
- `prd-roadmap.json`
- `prd-architecture.json`

To merge them into a single `prd-complete.json`, run the following command:

```bash
splan merge -o prd-complete.json \
  examples/merge/prd-metadata.json \
  examples/merge/prd-objectives.json \
  examples/merge/prd-personas.json \
  examples/merge/prd-user-stories.json \
  examples/merge/prd-requirements.json \
  examples/merge/prd-roadmap.json \
  examples/merge/prd-architecture.json
```

This will create `prd-complete.json` in your current directory, containing the contents of all seven files.

### Merge Order

The recommended merge order follows the logical structure of a PRD:

1. `metadata` - Document identity (always first)
2. `executive-summary` - High-level overview
3. `objectives` - Goals and success criteria
4. `personas` - Who we're building for
5. `user-stories` - What users need to do
6. `requirements` - Detailed specifications
7. `roadmap` - Delivery timeline
8. `architecture` - Technical design
9. `risks` - Risk assessment
10. Extended sections (market, solution, decisions, etc.)

**Important:** Later files override earlier files for duplicate keys. Place foundational content first.

---

## Instructions for AI Assistants

As an AI assistant, you can use this `merge` functionality to help users manage large structured documents. This section provides guidance for working with modular PRD files.

### When to Suggest Splitting

Proactively suggest splitting when:

- The PRD file exceeds ~1000 lines
- You encounter errors or lose content during updates
- The user reports issues with PRD management
- Multiple topics need frequent, independent updates

**Example response:** "This PRD is quite large (~1500 lines). To make updates more reliable, I recommend splitting it into modular files by topic. Would you like me to do that?"

### Task: Splitting a Large File

**Workflow:**

1. **Analyze the root keys** of the JSON document (e.g., `metadata`, `personas`, `roadmap`).
2. **For each root key**, create a new file named `prd-<key-name>.json`.
3. **Extract the corresponding object** and write it to the new file.
4. **Create an index file** listing all modular files (optional but helpful).
5. **Inform the user** that you have split the file and list the new modular files.

**Example split structure:**

```
myproduct/
├── prd-metadata.json
├── prd-executive-summary.json
├── prd-objectives.json
├── prd-personas.json
├── prd-user-stories.json
├── prd-requirements.json
├── prd-roadmap.json
├── prd-architecture.json
└── prd-complete.json  (merged output)
```

**Example Prompt:** "This PRD file is getting too large. Can you split it into modular files for me?"

### Task: Updating a Modular PRD

When working with modular files:

1. **Identify the relevant file** based on what needs to change.
2. **Read only that file** - don't read the complete merged PRD.
3. **Make targeted updates** to the specific file.
4. **Re-merge if needed** to update the complete document.

**Example:** To add a new persona, only read and update `prd-personas.json`.

### Task: Merging Modular Files

**Workflow:**

1. **Identify the modular files** the user wants to merge. Use `ls prd-*.json` or ask for clarification.
2. **Determine the correct merge order.** Follow the standard order: Metadata → Executive Summary → Objectives → Personas → User Stories → Requirements → Roadmap → Architecture → Risks → Extended sections.
3. **Construct and execute the `splan merge` command.**
4. **Validate the merged file** with `splan req prd validate`.
5. **Confirm the action** by stating that you have merged the files.

**Example command:**

```bash
splan merge -o prd-complete.json prd-metadata.json prd-objectives.json prd-personas.json prd-requirements.json prd-roadmap.json
```

**Example Prompt:** "Please merge our PRD files into a single document for final review."

### Task: Adding New Content

When adding substantial new content to a modular PRD:

1. **Determine if it fits an existing topic** - add to that file.
2. **If it's a new topic**, create a new modular file (e.g., `prd-market.json`).
3. **Update the merge command** to include the new file.

### Best Practices for AI Assistants

| Do | Don't |
|----|-------|
| Work with one modular file at a time | Try to update the entire merged PRD |
| Suggest splitting when files get large | Silently struggle with large files |
| Use topic-based naming | Use arbitrary numbering |
| Validate after merging | Skip validation |
| Keep modular files focused | Mix unrelated content in one file |

### Handling Topic Proliferation

If the number of topic files grows unwieldy (>15 files):

1. **Group related topics** into subdirectories:
   ```
   prd/
   ├── core/
   │   ├── metadata.json
   │   ├── objectives.json
   │   └── personas.json
   ├── requirements/
   │   ├── functional.json
   │   └── nonfunctional.json
   └── planning/
       ├── roadmap.json
       └── risks.json
   ```

2. **Use a manifest file** (`prd-manifest.json`) listing all files and merge order.

3. **Create a merge script** for consistent rebuilding.
