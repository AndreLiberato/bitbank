#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 2 || $# -gt 3 ]]; then
  echo "uso: $0 <source-branch> <base-branch> [target-sha]" >&2
  exit 1
fi

source_branch="$1"
base_branch="$2"
target_sha="${3:-}"

if [[ "$base_branch" != "production" ]]; then
  echo "branch de destino invalida: $base_branch (deve ser production)" >&2
  exit 1
fi

if [[ -n "$target_sha" ]]; then
  while IFS= read -r tag; do
    [[ "$tag" == rel-* ]] || continue
    echo "$tag"
    exit 0
  done < <(git tag --points-at "$target_sha")
fi

max_major=-1
max_minor=-1

while IFS= read -r tag; do
  [[ -z "$tag" ]] && continue
  if [[ "$tag" =~ ^rel-([0-9]+)\.([0-9]+)(\.[0-9]+)?$ ]]; then
    major="${BASH_REMATCH[1]}"
    minor="${BASH_REMATCH[2]}"
    if (( major > max_major || (major == max_major && minor > max_minor) )); then
      max_major="$major"
      max_minor="$minor"
    fi
  fi
done < <(git tag --list 'rel-*')

if (( max_major < 0 )); then
  max_major=1
  max_minor=0
fi

case "$source_branch" in
  stabilization/rc-*|origin/stabilization/rc-*)
    if [[ "$source_branch" =~ stabilization/rc-([0-9]+)\.([0-9]+)$ ]]; then
      base_major="${BASH_REMATCH[1]}"
      base_minor="${BASH_REMATCH[2]}"
    else
      echo "branch de estabilizacao de origem invalida: $source_branch" >&2
      exit 1
    fi

    expected_major="$max_major"
    expected_minor=$((max_minor + 1))
    
    if git rev-parse "rel-${base_major}.${base_minor}" >/dev/null 2>&1; then
      echo "a tag rel-${base_major}.${base_minor} ja existe" >&2
      exit 1
    fi

    if (( base_major != expected_major || base_minor != expected_minor )); then
      echo "a branch de origem $source_branch nao representa o proximo release esperado (rel-${expected_major}.${expected_minor})" >&2
      exit 1
    fi

    echo "rel-${base_major}.${base_minor}"
    ;;

  hotfix*|origin/hotfix*)
    release_line="rel-${max_major}.${max_minor}"
    highest_patch=-1

    while IFS= read -r tag; do
      [[ -z "$tag" ]] && continue
      if [[ "$tag" == "$release_line" ]]; then
        highest_patch=0
        continue
      fi
      if [[ "$tag" =~ ^rel-${max_major}\.${max_minor}\.([0-9]+)$ ]]; then
        patch="${BASH_REMATCH[1]}"
        if (( patch > highest_patch )); then
          highest_patch="$patch"
        fi
      fi
    done < <(git tag --list "${release_line}*")

    if (( highest_patch < 0 )); then
      highest_patch=0
    fi

    echo "${release_line}.$((highest_patch + 1))"
    ;;

  *)
    echo "branch de origem nao suportada para producao: $source_branch (use stabilization/rc-* ou hotfix*)" >&2
    exit 1
    ;;
esac
