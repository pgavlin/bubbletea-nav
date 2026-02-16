# Tasks: Message-Based Lifecycle Notifications

**Input**: Design documents from `/specs/003-lifecycle-msg-routing/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story. US1 and US2 share the core handleNav rewrite (placed in US1 since it's P1), then each story validates its specific behavior. US3 removes the old interface after both are working.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: Update the test mock to support message-based lifecycle delivery. MUST complete before any user story work.

- [x] T001 Rewrite `lifecycleScreen` test mock in lifecycle_test.go: remove `Appeared()`/`Disappeared()` methods, record lifecycle events in `Update` via type switch on `ScreenAppearedMsg` and `ScreenDisappearedMsg`. Keep `onAppear tea.Cmd` field to return commands from appeared handling. Add `onDisappear tea.Cmd` field to return commands from disappeared handling. Keep the `events []string` recording pattern.

**Checkpoint**: Test mock ready — user story implementation can begin.

---

## Phase 2: User Story 1 - Update Screen State on Appear (Priority: P1)

**Goal**: Screens receive `ScreenAppearedMsg` through their `Update` method on push, pop-reveal, and replace, allowing state changes that persist.

**Independent Test**: Push screen A, push screen B, pop screen B. Verify screen A receives `ScreenAppearedMsg` via `Update` and any state change persists in subsequent `View` calls.

### Implementation for User Story 1

- [x] T002 [US1] Rewrite `handleNav` in nav.go to deliver lifecycle messages via `screen.Update()` instead of `dispatchAppeared`/`dispatchDisappeared`. For each operation (push, pop, replace): call `screen.Update(ScreenDisappearedMsg{})` on the outgoing screen and `screen.Update(ScreenAppearedMsg{})` on the incoming screen. Store updated screens back into the stack slice for screens that remain. Collect all returned commands. Preserve disappear-before-appear ordering per FR-006/007/008. Keep `inLifecycle` guard and `pendingOps` mechanism.
- [x] T003 [US1] Rewrite `TestLifecyclePush` in lifecycle_test.go to verify: root receives `ScreenDisappearedMsg` (recorded as "disappeared"), detail receives `ScreenAppearedMsg` (recorded as "appeared"), both via `Update`.
- [x] T004 [P] [US1] Rewrite `TestLifecyclePop` in lifecycle_test.go to verify: popped screen receives "disappeared", revealed screen receives "appeared", both via `Update`.
- [x] T005 [P] [US1] Rewrite `TestLifecycleReplace` in lifecycle_test.go to verify: old screen receives "disappeared", replacement receives "appeared", root unaffected.
- [x] T006 [US1] Rewrite `TestLifecycleAppearedReturnsCmd` in lifecycle_test.go to verify: command returned from `Update` in response to `ScreenAppearedMsg` is included in the batched command output.
- [x] T007 [US1] Rewrite `TestLifecycleFullSequence` in lifecycle_test.go for message-based delivery: push A, push B, pop B, pop A — verify the same event sequence (root: [disappeared, appeared], A: [appeared, disappeared, appeared, disappeared], B: [appeared, disappeared]).
- [x] T008 [US1] Rewrite `TestLifecycleReentrantPush` in lifecycle_test.go: screen A's `Update` returns a `PushMsg` command in response to `ScreenAppearedMsg` — verify B is eventually pushed (depth 3).
- [x] T009 [US1] Add new `TestLifecycleAppearedStatePersists` in lifecycle_test.go: create a screen whose `Update` sets a `visited` flag when it receives `ScreenAppearedMsg`, changing its `View` output. Push it onto the stack. Verify `View` reflects the state change.

**Checkpoint**: ScreenAppearedMsg delivery works for all operations with state persistence. All appeared-side tests pass.

---

## Phase 3: User Story 2 - Clean Up Screen State on Disappear (Priority: P2)

**Goal**: Screens receive `ScreenDisappearedMsg` through their `Update` method, can update state (persisted for screens remaining in the stack), and can return commands.

**Independent Test**: Push screen B on top of screen A. Verify screen A receives `ScreenDisappearedMsg` via `Update`, state changes persist, and returned commands are collected.

### Implementation for User Story 2

- [x] T010 [US2] Add new `TestLifecycleDisappearedReturnsCmd` in lifecycle_test.go: create a screen whose `Update` returns a command when it receives `ScreenDisappearedMsg`. Push a new screen on top. Verify the command is included in the batched output.
- [x] T011 [US2] Add new `TestLifecycleDisappearedStatePersists` in lifecycle_test.go: create a screen whose `Update` sets a `hidden` flag when it receives `ScreenDisappearedMsg`, changing its `View` output. Push screen B on top. Pop screen B. Verify the previously-hidden screen's `View` reflects the state change from the disappeared handler.

**Checkpoint**: ScreenDisappearedMsg delivery works with state persistence and command return. All disappeared-side tests pass.

---

## Phase 4: User Story 3 - Remove LifecycleScreen Interface (Priority: P3)

**Goal**: Remove the `LifecycleScreen` interface from the public API. All lifecycle behavior works exclusively through messages.

**Independent Test**: Verify `LifecycleScreen` no longer exists as an exported type. Verify screens that don't handle lifecycle messages still work without errors.

### Implementation for User Story 3

- [x] T012 [US3] Remove `LifecycleScreen` interface from lifecycle.go. Keep `ScreenAppearedMsg` and `ScreenDisappearedMsg` type definitions. Update the doc comment on the remaining types to explain they are delivered via `Screen.Update`.
- [x] T013 [US3] Remove `dispatchAppeared` and `dispatchDisappeared` helper methods from nav.go (should be unused after T002).
- [x] T014 [US3] Rewrite `TestLifecycleNonLifecycleScreenSkipped` in lifecycle_test.go: rename to `TestLifecycleScreenIgnoresUnhandledMessages`. Verify that a `mockScreen` (which doesn't handle lifecycle messages in its `Update`) still works correctly through push/pop without panic.

**Checkpoint**: LifecycleScreen interface removed. Public API matches contract in contracts/api.go.txt.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all files and examples.

- [x] T015 Run `go test ./...` to verify all tests pass across all packages
- [x] T016 Run `go vet ./...` to verify no vet warnings
- [x] T017 Run `go build ./examples/...` to verify all examples compile

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: No dependencies — start immediately
- **US1 (Phase 2)**: Depends on T001 (test mock rewrite). T002 is the core code change; T003-T009 depend on T002.
- **US2 (Phase 3)**: Depends on T002 (handleNav rewrite from US1). T010-T011 are independent of US1 tests.
- **US3 (Phase 4)**: Depends on T002 (handleNav no longer calls dispatch helpers). T012-T014 can be parallelized.
- **Polish (Phase 5)**: Depends on all previous phases.

### User Story Dependencies

- **US1 (P1)**: Depends on Foundational only. Contains the core handleNav rewrite that US2 and US3 also depend on.
- **US2 (P2)**: Depends on T002 from US1 (the handleNav rewrite enables disappeared delivery).
- **US3 (P3)**: Depends on T002 from US1 (dispatch helpers become unused after the rewrite).

### Within Each User Story

- T002 (handleNav rewrite) must complete before any lifecycle test rewrites
- Test rewrites (T003-T008) can partially parallelize where marked [P]
- New tests (T009, T010, T011) depend on T002 but not on each other

### Parallel Opportunities

- T004 and T005 can run in parallel (different test functions, same file)
- T010 and T011 can run in parallel (different test functions)
- T012, T013, and T014 can run in parallel (different files/functions)
- T015, T016, T017 can run in parallel (read-only verification)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001: Rewrite test mock
2. Complete T002: Rewrite handleNav (core change)
3. Complete T003-T009: Validate appeared delivery
4. **STOP and VALIDATE**: `go test ./...` — all appeared tests pass

### Incremental Delivery

1. T001 → Foundational ready
2. T002-T009 → US1 complete (appeared delivery works with state persistence)
3. T010-T011 → US2 complete (disappeared delivery validated)
4. T012-T014 → US3 complete (old interface removed)
5. T015-T017 → Full verification

---

## Notes

- T002 is the pivotal task — it rewrites `handleNav` to use `screen.Update()` for both appeared and disappeared, enabling all three user stories.
- The `inLifecycle` guard and `pendingOps` queue are preserved unchanged (per research.md R3).
- For screens being removed (popped, replaced-away), `Update` is still called but the returned screen is discarded — only commands are collected.
- All existing lifecycle test names are preserved (rewritten, not renamed) except T014 which renames for clarity.
