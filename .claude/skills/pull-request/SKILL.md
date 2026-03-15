---
name: pull-request
description: Use when the user wants to open a pull request. Reads commits on the current branch, derives a title and body, and opens a PR with gh pr create.
version: 0.1.0
disable-model-invocation: false
user-invocable: true
---

# Open a Pull Request

Read the commits on this branch and open a pull request with a title and body
derived from those commits.

## Important

The PR title and body become the squash commit message on `main`. Write them
as you would a commit message: title is the subject line, body is the context
(what changed and why).

## Process

### Step 1: Read the branch commits

```sh
git log main..HEAD --format="%s%n%n%b"
```

If that is empty (branch is up to date with main), stop and tell the user there
are no new commits to open a PR for.

### Step 2: Derive the title and body

**Title:**
- Summarize all commits in one imperative-mood subject line (≤72 characters)
- If there is only one commit, use its subject line directly
- If there are multiple commits, synthesize a single summary

**Body:** Follow the template and rules in [template.md](template.md).

### Step 3: Show the title and body to the user and ask for confirmation

Present the title and body you derived, then ask:

**Proceed? (yes/no)**

Do not open the PR until the user confirms with yes.

### Step 4: Open the PR

```sh
gh pr create --title "<title>" --body "<body>"
```
