#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "uso: $0 <arquivo-da-mensagem>" >&2
  exit 1
fi

message_file="$1"

if [[ ! -f "$message_file" ]]; then
  echo "arquivo de mensagem nao encontrado: $message_file" >&2
  exit 1
fi

subject_line="$(head -n 1 "$message_file" | tr -d '\r')"

if [[ ! "$subject_line" =~ ^#([0-9]+)[[:space:]][-–][[:space:]].+ ]]; then
  cat >&2 <<'EOF'
Mensagem de commit invalida.
Use o formato:
  #NUM_ISSUE - MENSAGEM
Exemplo:
  #33 - Corrige mensagem da tela principal
EOF
  exit 1
fi

issue_number="${BASH_REMATCH[1]}"

if [[ "${SKIP_ISSUE_LOOKUP:-0}" == "1" ]]; then
  exit 0
fi

if [[ -n "${GITHUB_REPOSITORY:-}" ]]; then
  repository="$GITHUB_REPOSITORY"
else
  remote_url="$(git remote get-url origin 2>/dev/null || true)"

  if [[ "$remote_url" =~ github\.com[:/]([^/]+/[^/.]+)(\.git)?$ ]]; then
    repository="${BASH_REMATCH[1]}"
  else
    echo "nao foi possivel identificar o repositorio do GitHub para validar a issue" >&2
    exit 1
  fi
fi

api_url="https://api.github.com/repos/${repository}/issues/${issue_number}"
response_file="$(mktemp)"
trap 'rm -f "$response_file"' EXIT

curl_args=(
  --silent
  --show-error
  --location
  --output "$response_file"
  --write-out "%{http_code}"
  -H "Accept: application/vnd.github+json"
)

if [[ -n "${GITHUB_TOKEN:-}" ]]; then
  curl_args+=(-H "Authorization: Bearer ${GITHUB_TOKEN}")
fi

http_code="$(curl "${curl_args[@]}" "$api_url" || true)"

if [[ "$http_code" != "200" ]]; then
  case "$http_code" in
    404)
      echo "a issue #${issue_number} nao foi encontrada no repositorio ${repository}." >&2
      ;;
    401|403)
      echo "nao foi possivel validar a issue #${issue_number} por falta de autorizacao ou limite da API do GitHub." >&2
      ;;
    000)
      echo "nao foi possivel validar a issue #${issue_number} porque a consulta ao GitHub falhou." >&2
      ;;
    *)
      echo "a validacao da issue #${issue_number} falhou com codigo HTTP ${http_code}." >&2
      ;;
  esac
  if [[ -z "${GITHUB_TOKEN:-}" ]]; then
    echo "se necessario, configure GITHUB_TOKEN para evitar falhas por limite da API." >&2
  fi
  exit 1
fi

if grep -q '"pull_request"' "$response_file"; then
  echo "o numero #${issue_number} corresponde a um pull request, nao a uma issue." >&2
  exit 1
fi
