# ADR-012: Signed Commits Required

**Context:** Commits from contributors need to be verifiable. Unsigned commits
can be trivially authored with any name and email, making it impossible to
distinguish legitimate contributions from spoofed ones. GitHub displays a
"Verified" badge on signed commits, giving reviewers a clear trust signal.

**Decision:** All commits must be cryptographically signed — either via GPG or
SSH. Contributors configure signing locally; setup instructions are in
`RELEASING.md`.

**Consequences:**
- Every commit carries a verifiable identity tied to a key the author controls.
- GitHub shows "Verified" on all commits in the repository.
- Contributors must configure signing once before their first commit.
