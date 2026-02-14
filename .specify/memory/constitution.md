<!--
  Sync Impact Report
  ===================
  Version change: 1.0.1 → 1.1.0
  Modified principles:
    - II. Bubble Tea Idioms: distinguished top-level components
      (MUST implement tea.Model) from helper value types (MAY
      follow bubbles convention instead)
  Added sections: none
  Removed sections: none
  Templates requiring updates:
    - .specify/templates/plan-template.md ✅ no changes needed
    - .specify/templates/spec-template.md ✅ no changes needed
    - .specify/templates/tasks-template.md ✅ no changes needed
    - .specify/templates/checklist-template.md ✅ no changes needed
    - .specify/templates/agent-file-template.md ✅ no changes needed
  Follow-up TODOs: none
-->

# bubbletea-nav Constitution

## Core Principles

### I. Library-First

bubbletea-nav is a Go library, not a CLI tool or application.
Every feature MUST be importable as a standard Go package.

- All exported types, functions, and interfaces MUST be usable
  via `go get` and standard `import` without additional tooling.
- No `main` packages in the core module. Examples and demos MAY
  live in a separate `examples/` directory.
- Each package MUST have a clear, single responsibility.

**Rationale**: Consumers integrate bubbletea-nav into their own
Bubble Tea applications. A library-first design ensures minimal
coupling and maximum reuse.

### II. Bubble Tea Idioms

All navigation components MUST follow the Bubble Tea
Model-Update-View architecture.

- Top-level components (those passed to `tea.NewProgram`) MUST
  implement the `tea.Model` interface (`Init`, `Update`, `View`).
- Helper value types composed within screens MAY follow the
  bubbles convention (`Update` returning their own type) instead.
- State transitions MUST occur through `tea.Msg` and `tea.Cmd`
  — never through direct mutation of shared state.
- Components MUST NOT spawn goroutines or perform I/O outside
  of `tea.Cmd` functions.

**Rationale**: Consistency with Bubble Tea idioms means consumers
can drop bubbletea-nav components into existing applications
without learning a new programming model.

### III. Composable API

The public API MUST be small, focused, and composable.

- Prefer many small types over few large ones.
- Navigation components MUST compose with standard Bubble Tea
  patterns (embedding, message passing, command batching).
- Avoid God-objects: no single type should own the entire
  navigation lifecycle.
- Interfaces SHOULD be narrow (1-3 methods) when defined.

**Rationale**: Small composable primitives let consumers build
exactly the navigation patterns they need without fighting
the library's opinions.

### IV. Test-Required

All public APIs MUST have test coverage.

- Every exported function, method, and type MUST have at least
  one test exercising its primary behavior.
- Tests MAY be written alongside or after implementation (strict
  TDD is not required).
- Table-driven tests MUST be used when testing multiple input
  variations of the same function.
- Tests MUST be runnable via `go test ./...` with no external
  dependencies or environment setup.

**Rationale**: A navigation library is foundational infrastructure
for its consumers. Untested APIs erode trust and make upgrades
risky.

### V. Backward Compatibility

Public API changes MUST follow Go module semantic versioning.

- Removing or renaming an exported symbol is a breaking change
  and MUST trigger a major version bump.
- Adding new exported symbols or optional parameters is a minor
  change.
- Bug fixes and documentation changes are patch changes.
- Deprecated symbols MUST be annotated with `// Deprecated:`
  comments and retained for at least one minor release before
  removal.

**Rationale**: Consumers depend on stable APIs. Predictable
versioning lets them upgrade with confidence.

### VI. Observability

Navigation state MUST be inspectable for debugging and testing.

- All key state types MUST implement `fmt.Stringer` or provide
  a `String()` method with human-readable output.
- Navigation events (push, pop, replace) MUST be representable
  as `tea.Msg` types that consumers can log or inspect.
- Debug helpers (e.g., stack dump, route trace) SHOULD be
  provided in a separate `debug` sub-package.

**Rationale**: TUI debugging is harder than web or CLI debugging.
Making navigation state visible reduces development friction.

### VII. Simplicity

Start with the simplest solution that works. Apply YAGNI.

- Do not add features, abstractions, or configuration until a
  concrete use case demands them.
- Prefer explicit code over clever code. Three similar lines
  are better than a premature abstraction.
- Every new exported type or function MUST justify its existence
  with a concrete use case in the examples or tests.

**Rationale**: Complexity in a library compounds in every
application that depends on it. Keeping the surface area small
makes the library easier to learn, use, and maintain.

## Technical Constraints

- **Language**: Go (minimum version determined by `go.mod`).
- **Primary Dependency**: [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).
- **Build/Test**: `go build ./...` and `go test ./...` MUST pass
  with zero external service dependencies.
- **Linting**: `go vet ./...` MUST pass. Projects SHOULD also
  run `staticcheck` or `golangci-lint`.
- **Documentation**: All exported symbols MUST have Go doc
  comments. Package-level doc comments MUST include a usage
  example.

## Development Workflow

- **Branching**: Feature work MUST occur on a branch off `main`.
  Merge via pull request after review.
- **Commits**: Use conventional commit messages
  (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`).
- **Review gates**: PRs MUST pass `go test ./...` and `go vet`
  before merge. Test coverage MUST NOT decrease on public APIs.
- **Release**: Tag releases with `vMAJOR.MINOR.PATCH`. Update
  CHANGELOG before tagging.

## Governance

This constitution is the authoritative source for project
principles and development standards. It supersedes informal
conventions or ad-hoc decisions.

- **Amendments**: Any change to this constitution MUST be
  documented with a version bump, rationale, and updated
  `Last Amended` date. Amendments follow semantic versioning
  (see Principle V for version bump rules applied to this
  document).
- **Compliance**: All pull requests and code reviews MUST verify
  alignment with these principles. Violations MUST be flagged
  and resolved before merge.
- **Exceptions**: Temporary deviations MUST be justified in the
  PR description and tracked as tech debt with a follow-up
  issue.

**Version**: 1.1.0 | **Ratified**: 2026-02-14 | **Last Amended**: 2026-02-14
