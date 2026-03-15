# Pull Request Body Template

```
## Problem
<one sentence explaining why this change is needed>

## Changes
- <what changed and why — add context, do not restate commit messages>
- <additional change if needed>

## Decisions
- <ADR title and link, e.g. [ADR-013: Branch protection](docs/decisions/013-branch-protection.md)>
← omit this section if no new or relevant ADRs

---
Closes #<issue number>  ← omit if no related issue
⚠️ Breaking change: <description>  ← omit if no breaking change
```

## Rules

- Omit the `Decisions` section if no ADRs are new or relevant to this PR
- Omit the `Closes` line if no related issue; do not guess an issue number
- Include the breaking change line only if a public API is removed, renamed, or
  its behavior changes in a way that requires callers to update their code
- Do not add sections for testing — CI is the test record
