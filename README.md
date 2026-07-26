# gastown-mod

**Status: archived.** A January 2026 experiment fork of [gastownhall/gastown](https://github.com/gastownhall/gastown), kept as work history. Active Gas City orchestration work moved to [gascity](https://github.com/sjarmak/gascity).

## What this is

Gas Town is a Go multi-agent orchestrator for Claude Code: work is tracked in convoys, slung to worker agents (polecats), and persisted in a git-backed issue ledger ([beads](https://github.com/sjarmak/beads)). This fork carries exactly one commit on top of upstream: `2c42fe71` (`gt-9x8.3`), which added pluggable runtime support so agent sessions could launch on backends other than Claude Code, with OpenHands as the first alternative.

## What the commit did

Upstream at the time could swap the launch command via `RuntimeConfig`, but every backend got identical treatment. The commit introduced provider-aware adapters in a new `internal/llm` package:

- `adapter.go` defines an `Adapter` interface (`BuildLaunchCommand`) plus `LaunchOptions` carrying the session's role, actor, rig, and runtime settings.
- `registry.go` resolves a provider name to a registered adapter; `claude_adapter.go`, `openhands_adapter.go`, and `shell_adapter.go` are the three implementations.
- The OpenHands adapter maps Gas Town session identity to `OPENHANDS_AGENT_ROLE`, `OPENHANDS_AGENT_NAME`, `OPENHANDS_RIG_NAME`, and `OPENHANDS_MODEL` environment variables, and defaults the launch command to `openhands --exp`.
- `internal/config` gained `Provider` and `Model` fields on `RuntimeConfig`, plus a per-role `runtime_overrides` map keyed by `GT_ROLE`, so a rig could run different roles on different backends.

The change added 688 lines across 10 source files, including unit tests in `internal/llm/adapter_test.go` and `internal/config/loader_test.go` (`go test ./internal/llm ./internal/config`).

## Status

The experiment ended here; the commit was never merged upstream and the fork remains one commit ahead of `gastownhall/gastown`. At archive time the branch sat 5,694 commits behind upstream, so it does not build a current Gas Town and is not maintained. The repo stays up because the commit documents a working pattern for multi-runtime agent launches in Gas Town's architecture.

## Related

- [gascity](https://github.com/sjarmak/gascity), fork of gastownhall/gascity, where current multi-agent orchestration work happens.
- [beads](https://github.com/sjarmak/beads), fork of the git-backed issue tracker Gas Town builds on.
- [gastownhall/gastown](https://github.com/gastownhall/gastown), the upstream this fork diverged from.
