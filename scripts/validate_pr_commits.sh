#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "uso: $0 <numero-do-pull-request>" >&2
  exit 1
fi

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
  echo "GITHUB_TOKEN nao configurado para validar os commits do pull request." >&2
  exit 1
fi

if [[ -z "${GITHUB_REPOSITORY:-}" ]]; then
  echo "GITHUB_REPOSITORY nao configurado." >&2
  exit 1
fi

pr_number="$1"
api_url="https://api.github.com/repos/${GITHUB_REPOSITORY}/pulls/${pr_number}/commits?per_page=100"
response_file="$(mktemp)"
messages_dir="$(mktemp -d)"
trap 'rm -f "$response_file"; rm -rf "$messages_dir"' EXIT

http_code="$(
  curl \
    --silent \
    --show-error \
    --location \
    --output "$response_file" \
    --write-out "%{http_code}" \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    "$api_url"
)"

if [[ "$http_code" != "200" ]]; then
  echo "nao foi possivel obter os commits do pull request #${pr_number}." >&2
  exit 1
fi

python3 - "$response_file" "$messages_dir" <<'PY'
import json
import pathlib
import sys

response_path = pathlib.Path(sys.argv[1])
messages_dir = pathlib.Path(sys.argv[2])
commits = json.loads(response_path.read_text())

if not commits:
    print("Nenhum commit encontrado no pull request.", file=sys.stderr)
    sys.exit(1)

for index, commit in enumerate(commits, start=1):
    message = commit["commit"]["message"]
    (messages_dir / f"{index}.txt").write_text(message, encoding="utf-8")
PY

for message_file in "$messages_dir"/*.txt; do
  scripts/validate_commit_message.sh "$message_file"
done
