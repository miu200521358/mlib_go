#!/usr/bin/env bash
set -euo pipefail

dry_run=0
for arg in "$@"; do
  case "$arg" in
    --dry-run|-n)
      dry_run=1
      ;;
    *)
      echo "Usage: $0 [--dry-run|-n]" >&2
      exit 2
      ;;
  esac
done

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$repo_root" ]; then
  echo "git repository not found." >&2
  exit 1
fi

echo "Repository root: $repo_root"

instruction="$(printf '// \u6307\u793A: miu200521358')"

mapfile -d '' go_files < <(find "$repo_root" -type f -name '*.go' \
  -not -path '*/.git/*' \
  -not -path '*/vendor/*' \
  -not -path '*/node_modules/*' \
  -not -path '*/.cache/*' \
  -print0)

if [ "${#go_files[@]}" -eq 0 ]; then
  echo "No .go files detected."
  exit 0
fi

targets=()
for file in "${go_files[@]}"; do
  if grep -Fq "$instruction" "$file"; then
    continue
  fi
  rel="${file#"$repo_root"/}"
  targets+=("$rel")
done

if [ "${#targets[@]}" -eq 0 ]; then
  echo "All .go files already have instruction."
  exit 0
fi

echo "Missing instruction header:"
for rel in "${targets[@]}"; do
  echo "  $rel"
done

if [ "$dry_run" -eq 1 ]; then
  exit 0
fi

for rel in "${targets[@]}"; do
  file="$repo_root/$rel"
  python - "$file" "$instruction" <<'PY'
import re
import sys
from pathlib import Path

path = Path(sys.argv[1])
instruction = sys.argv[2]

content = path.read_text(encoding="utf-8")
eol = "\r\n" if "\r\n" in content else "\n"
has_trailing_newline = content.endswith(eol)
lines = re.split(r"\r?\n", content)

insert_index = 0
saw_build_tag = False
while insert_index < len(lines):
    line = lines[insert_index]
    if re.match(r"^\s*//go:build\b", line):
        saw_build_tag = True
        insert_index += 1
        continue
    if re.match(r"^\s*// \+build\b", line):
        saw_build_tag = True
        insert_index += 1
        continue
    if saw_build_tag and re.match(r"^\s*$", line):
        insert_index += 1
        continue
    break

new_lines = lines[:insert_index] + [instruction] + lines[insert_index:]
new_content = eol.join(new_lines)
if has_trailing_newline and not new_content.endswith(eol):
    new_content += eol

with open(path, "w", encoding="utf-8", newline="") as f:
    f.write(new_content)
PY
  echo "Updated: $rel"
done
