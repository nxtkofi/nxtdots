# Purpose

You bootstrap technical documentation by creating `/docs` structure and analyzing the codebase with multiple Haiku explore agents. Generate high-level, maintainable documentation based on actual code.

## Process

**Step 1: Pre-flight checks**
- Check if `/docs` exists (ask before overwriting)
- Identify project type (check package.json, requirements.txt, etc.)
- Assess codebase size to adjust analysis depth

**Step 2: Create structure**
```
/docs
├── README.md
├── architecture/ (overview.md, data-flows.md, decisions/)
├── guides/ (getting-started.md, common-patterns.md)
├── reference/ (tech-stack.md)
└── diagrams/
```

**Step 3: Launch 6 parallel Haiku explore agents**

Use Task tool to spawn these agents simultaneously:

1. **Architecture Hunter** - Identify architectural pattern (MVC, microservices, feature-based, etc.) and directory organization philosophy
2. **Utilities Scout** - Find shared utilities/hooks/services in /utils, /lib, /hooks, /common
3. **Tech Stack Auditor** - List frameworks, libraries, build tools from package files
4. **Data Flow Tracer** - Trace key flows (auth, API requests, state management)
5. **Integration Mapper** - Identify external services (APIs, databases, auth providers)
6. **Patterns Analyst** - Find coding patterns and conventions

**Step 4: Synthesize findings**
Read all 6 agent outputs and synthesize into coherent, high-level documentation.

**Step 5: Generate documentation files**

Create with real content (no placeholders):
- **README.md** - Navigation hub
- **overview.md** - Directory structure + purposes, architecture pattern (NO file counts/listings)
- **getting-started.md** - Installation only
- **common-patterns.md** - Patterns conceptually, 1-2 example files max
- **data-flows.md** - High-level flows, simple diagrams (4-8 boxes)
- **tech-stack.md** - Frameworks and tools
- **0001-initial-architecture.md** - Baseline ADR

## Principles

**Patterns over files** - Concepts, not listings  
**High-level only** - No implementations or type signatures  
**Natural wikilinks** - Inline contextual: `uses [[data-flows#auth|JWT]]`  
**Simple diagrams** - 4-8 boxes maximum

## Anti-Patterns (Never Do This)

❌ List files/components individually  
❌ Include directory trees (`tree` output)  
❌ Document type signatures  
❌ Add file counts (e.g., "9 files")  
❌ Forced "Related documentation" sections  
❌ Detailed implementation code in docs

## Quality Check

Good docs: won't break when files rename, focus on concepts, natural prose with contextual links, maintainable over time.
