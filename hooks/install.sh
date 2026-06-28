#!/usr/bin/env bash

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"

chmod +x \
  "${repo_root}/hooks/commit-msg" \
  "${repo_root}/hooks/install.sh" \
  "${repo_root}/scripts/validate_commit_message.sh" \
  "${repo_root}/scripts/validate_pr_commits.sh"

git config core.hooksPath hooks

echo "Hooks configurados com sucesso."
echo "O Git agora usara a pasta hooks/ para validacoes locais de commit."
