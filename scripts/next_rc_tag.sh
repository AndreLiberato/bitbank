#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 2 || $# -gt 3 ]]; then
  echo "uso: $0 <source-branch> <base-branch> [target-sha]" >&2
  exit 1
fi

source_branch="$1"
base_branch="$2"
target_sha="${3:-}"

if [[ ! "$base_branch" =~ ^stabilization/rc-([0-9]+)\.([0-9]+)$ ]]; then
  echo "branch de estabilizacao invalida: $base_branch" >&2
  exit 1
fi

base_major="${BASH_REMATCH[1]}"
base_minor="${BASH_REMATCH[2]}"
release_line="rc-${base_major}.${base_minor}"

if [[ -n "$target_sha" ]]; then
  while IFS= read -r tag; do
    [[ "$tag" == rc-* ]] || continue
    echo "$tag"
    exit 0
  done < <(git tag --points-at "$target_sha")
fi

max_major=-1
max_minor=-1

while IFS= read -r tag; do
  [[ -z "$tag" ]] && continue
  if [[ "$tag" =~ ^rc-([0-9]+)\.([0-9]+)(\.[0-9]+)?$ ]]; then
    major="${BASH_REMATCH[1]}"
    minor="${BASH_REMATCH[2]}"
    if (( major > max_major || (major == max_major && minor > max_minor) )); then
      max_major="$major"
      max_minor="$minor"
    fi
  fi
done < <(git tag --list 'rc-*')

case "$source_branch" in
  main)
    if git rev-parse "$release_line" >/dev/null 2>&1; then
      echo "a tag $release_line ja existe" >&2
      exit 1
    fi

    if (( max_major >= 0 )); then
      expected_major="$max_major"
      expected_minor=$((max_minor + 1))
      if (( base_major != expected_major || base_minor != expected_minor )); then
        echo "a branch $base_branch nao representa a proxima estabilizacao esperada (rc-${expected_major}.${expected_minor})" >&2
        exit 1
      fi
    fi

    echo "$release_line"
    ;;

  bugfix*)
    highest_patch=-1

    while IFS= read -r tag; do
      [[ -z "$tag" ]] && continue
      if [[ "$tag" == "$release_line" ]]; then
        highest_patch=0
        continue
      fi
      if [[ "$tag" =~ ^rc-${base_major}\.${base_minor}\.([0-9]+)$ ]]; then
        patch="${BASH_REMATCH[1]}"
        if (( patch > highest_patch )); then
          highest_patch="$patch"
        fi
      fi
    done < <(git tag --list "${release_line}*")

    if (( highest_patch < 0 )); then
      echo "nao existe uma tag base ${release_line} para iniciar bugfixes nesta estabilizacao" >&2
      exit 1
    fi

    echo "${release_line}.$((highest_patch + 1))"
    ;;

  *)
    echo "branch de origem nao suportada: $source_branch (use main ou bugfix*)" >&2
    exit 1
    ;;
esac
