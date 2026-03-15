# ADR-013: Branch Protection for main

**Context:** With the project open to external contributors, unreviewed or
untested code could land on `main` directly. Without guardrails, CI failures,
unresolved review comments, or stale approvals could all slip through.

**Decision:** The `main` branch is protected with the following rules, applied
via `scripts/setup-branch-protection.sh`:

- **Pull request required.** No direct pushes to `main` — all changes go
  through a PR.
- **One approving review required.** At least one reviewer (other than the PR
  author) must approve before merge.
- **Stale reviews dismissed.** Pushing new commits to a PR invalidates existing
  approvals, requiring re-review of the updated code.
- **All CI jobs must pass.** Every job in `.github/workflows/ci.yml` is a
  required status check: `Lint`, `Test (Go 1.21–1.24)`, and `Script tests`.
- **Branch must be up-to-date.** The PR branch must be current with `main`
  before merge. Contributors can use the "Update branch" button on the PR.
- **Conversations must be resolved.** All review comments must be addressed
  before merge.
- **Squash merge only.** PRs are merged as a single commit using the PR title
  and body. This keeps `git log` on `main` clean and readable.
- **Force pushes and deletions blocked.** `main` is append-only.
- **Admin bypass enabled.** The maintainer can override protections when
  necessary (e.g., recovering from a broken CI configuration).

**Consequences:**
- Contributors get a clear, consistent process: open a PR, pass CI, get a
  review, resolve comments, merge.
- `main` stays releasable at all times.
- The maintainer retains the ability to unblock themselves without removing
  protections entirely.
