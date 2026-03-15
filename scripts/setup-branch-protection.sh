#!/usr/bin/env bash
#
# setup-branch-protection.sh
#
# Configures branch protection and repository merge settings for
# mycargus/riteway-golang.
#
# Prerequisites:
#   - gh CLI installed and authenticated (gh auth status)
#   - You must be an admin on the repository
#
# Run: bash scripts/setup-branch-protection.sh
#
# To verify results afterward:
#   gh api repos/mycargus/riteway-golang/branches/main/protection | jq .

set -euo pipefail

echo "==> Configuring mycargus/riteway-golang..."

# -- Repo-level merge settings --
#
# allow_squash_merge:          Contributors can squash all PR commits into one commit on main.
# allow_rebase_merge:          Disabled — squash only, one consistent merge path for contributors.
# allow_merge_commit:          Disabled — prevents non-linear history on main.
# allow_update_branch:         Enables the "Update branch" button on PRs so contributors can
#                              sync with main without using the CLI (required since we enforce
#                              branches be up-to-date before merge).
# delete_branch_on_merge:      Automatically deletes PR branches after merge.
# squash_merge_commit_title:   Uses the PR title as the squash commit subject line.
# squash_merge_commit_message: Uses the PR body as the squash commit message body.
gh api \
  --method PATCH \
  repos/mycargus/riteway-golang \
  --field allow_squash_merge=true \
  --field allow_rebase_merge=false \
  --field allow_merge_commit=false \
  --field allow_update_branch=true \
  --field delete_branch_on_merge=true \
  --field squash_merge_commit_title=PR_TITLE \
  --field squash_merge_commit_message=PR_BODY \
  > /dev/null

echo "    repo merge settings OK"

# -- Branch protection for main --
#
# required_status_checks.strict:    PR branch must be up-to-date with main before merge.
#                                   Prevents "works on my machine" merges in a trunk-based workflow.
# required_status_checks.contexts:  Every CI job must pass. Names must match the `name:` fields
#                                   in .github/workflows/ci.yml exactly (including matrix variants).
# enforce_admins:                   false — repo admins (you) can bypass rules when needed.
# required_pull_request_reviews:
#   required_approving_review_count: At least one approval required before merge.
#   dismiss_stale_reviews:           Approval is revoked if new commits are pushed,
#                                    forcing re-review of the updated code.
#   require_code_owner_reviews:      false — no CODEOWNERS file is configured.
# restrictions:                     null — no restriction on who can push to main
#                                    (PRs are enforced by required_pull_request_reviews instead).
# required_conversation_resolution: All review comments must be resolved before merge.
# allow_force_pushes:               false — protects commit history on main.
# allow_deletions:                  false — prevents accidental deletion of main.
gh api \
  --method PUT \
  repos/mycargus/riteway-golang/branches/main/protection \
  --input - <<'EOF'
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "Lint",
      "Test (Go 1.21)",
      "Test (Go 1.22)",
      "Test (Go 1.23)",
      "Test (Go 1.24)",
      "Script tests"
    ]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": false
  },
  "restrictions": null,
  "required_conversation_resolution": true,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF

echo "    branch protection OK"
echo ""
echo "Done. Verify with:"
echo "  gh api repos/mycargus/riteway-golang/branches/main/protection | jq ."
