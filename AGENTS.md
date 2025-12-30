# AGENTS

- Architecture, naming, errors, and testing follow `mlib-guiding-principles` at `mlib_skills\skills\00_project\010_mlib_guiding_principles\references\001_guiding_principles.md`.
- Rebuild phases follow `mlib-rebuild-procedure` at `mlib_skills\skills\00_project\020_mlib_rebuild_procedure\references\001_rebuild_procedure.md`.
- Communication: outputs, comments, and responses are in Japanese.
- File encoding: save/update it as UTF-8.
- Comments/logs should be verbose enough to show intent and processing status.
- Public docs/logs should avoid drive-letter absolute paths; use relative or redacted paths.
- Internal-only files (non-build, non-test) live under `internal/`.
- Do not create empty Go packages; add packages only when needed.
- Local/CI test entrypoint: `internal/scripts/test.ps1`.
