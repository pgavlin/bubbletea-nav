---

description: "Task list for Screen Navigation Library implementation"
---

# Tasks: Screen Navigation Library

**Input**: Design documents from `/specs/001-screen-nav/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/api.go

**Tests**: Required by Constitution Principle IV (Test-Required). Tests are written alongside implementation within each user story phase.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Single Go module at repository root: `nav.go`, `focus.go`, `lifecycle.go`
- Tests alongside source: `nav_test.go`, `focus_test.go`, `lifecycle_test.go`
- Examples in `_examples/basic/main.go`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize Go module and project structure

- [x] T001 Initialize Go module with `go mod init`, add `github.com/charmbracelet/bubbletea` v1.3.x dependency via `go get` in go.mod
- [x] T002 Create doc.go with `package nav` declaration and package-level doc comment describing the library purpose

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define all interfaces and message types that all user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 [P] Define Screen interface (Init, Update, View) and navigation message types (PushMsg, PopMsg, ReplaceMsg) with Push(), Pop(), Replace() command functions in nav.go
- [x] T004 [P] Define Focusable interface (Focus, Blur, Focused), Bounded interface (Bounds), and FocusChangedMsg type in focus.go
- [x] T005 [P] Define LifecycleScreen interface (extends Screen with Appeared, Disappeared) and lifecycle message types (ScreenAppearedMsg, ScreenDisappearedMsg) in lifecycle.go

**Checkpoint**: All interfaces and message types defined. User story implementation can now begin.

---

## Phase 3: User Story 1 - Screen Stack Navigation (Priority: P1) MVP

**Goal**: Implement the navigation stack with push/pop/replace operations and full state preservation for inactive screens.

**Independent Test**: Create two mock screens, push the second, verify it is active, pop it, verify the first is restored with original state.

### Implementation for User Story 1

- [x] T006 [US1] Implement NewStack(root Screen) constructor that validates non-nil root and initializes the screens slice in nav.go
- [x] T007 [US1] Implement Stack.Init() that delegates to the root screen's Init() method in nav.go
- [x] T008 [US1] Implement Stack.View() that delegates to the active (top) screen's View() method in nav.go
- [x] T009 [US1] Implement PushMsg handling in Stack.Update(): append new screen to screens slice, call new screen's Init(), return updated stack in nav.go
- [x] T010 [US1] Implement PopMsg handling in Stack.Update(): remove top screen if depth > 1, restore previous screen as active; ignore pop when only root remains in nav.go
- [x] T011 [US1] Implement ReplaceMsg handling in Stack.Update(): replace top screen with new screen, call new screen's Init() in nav.go
- [x] T012 [US1] Implement message forwarding in Stack.Update(): forward non-navigation messages to active screen's Update(), propagate returned Screen and Cmd in nav.go
- [x] T013 [US1] Implement Stack.String() returning "Stack[N screens, active: N-1]" and Stack.Depth() returning len(screens) in nav.go
- [x] T014 [P] [US1] Write table-driven tests for NewStack (nil panic), Init delegation, View delegation, Depth, and String in nav_test.go
- [x] T015 [P] [US1] Write table-driven tests for push (new screen active), pop (previous restored with state), pop-on-root (ignored), replace (swaps top), and message forwarding in nav_test.go

**Checkpoint**: Navigation stack fully functional with push/pop/replace. MVP complete.

---

## Phase 4: User Story 2 - Component Focus Navigation (Priority: P2)

**Goal**: Implement the FocusManager for Tab/Shift+Tab sequential focus cycling and mouse-click focus targeting within a screen.

**Independent Test**: Create a FocusManager with three mock focusable items, verify Tab cycles forward with wrap, Shift+Tab cycles backward with wrap, and only the focused item reports Focused() == true.

### Implementation for User Story 2

- [x] T016 [US2] Implement NewFocusManager(items ...Focusable) constructor that stores items and calls Focus() on the first item (or sets focusIndex to -1 if empty) in focus.go
- [x] T017 [US2] Implement FocusManager.Update() for tea.KeyTab (advance focusIndex, wrap cyclically, call Blur on old and Focus on new) and tea.KeyShiftTab (reverse direction) in focus.go
- [x] T018 [US2] Implement FocusManager.Update() for tea.MouseMsg: iterate items implementing Bounded, check if click (X,Y) is within Bounds(), move focus if hit; pass through non-matching clicks without changing focus in focus.go
- [x] T019 [US2] Implement FocusManager.SetItems(items ...Focusable) that replaces the items list and resets focus to first item (or -1 if empty) in focus.go
- [x] T020 [US2] Implement FocusManager.FocusedIndex() returning current index (-1 if none) and FocusManager.FocusIndex(index int) that sets focus to a specific index with Blur/Focus calls in focus.go
- [x] T021 [US2] Implement FocusManager.String() returning "FocusManager[N items, focused: M]" in focus.go
- [x] T022 [P] [US2] Write table-driven tests for Tab forward cycling, Shift+Tab backward cycling, cyclic wrap in both directions, and empty-list no-op in focus_test.go
- [x] T023 [P] [US2] Write table-driven tests for mouse click targeting (hit within bounds moves focus), click on non-bounded item (no change), click outside all bounds (no change), and FocusChangedMsg emission in focus_test.go
- [x] T024 [US2] Write tests for SetItems (resets focus), FocusedIndex, FocusIndex (valid index, -1 blurs all, out-of-range), and dynamic list changes in focus_test.go

**Checkpoint**: FocusManager fully functional with keyboard and mouse focus cycling.

---

## Phase 5: User Story 3 - Screen Lifecycle Events (Priority: P3)

**Goal**: Add lifecycle event dispatching (Appeared/Disappeared) to Stack operations and implement re-entrant stack modification prevention.

**Independent Test**: Create two mock LifecycleScreen implementations that record events, push and pop them, verify the correct sequence of Appeared/Disappeared calls.

### Implementation for User Story 3

- [x] T025 [US3] Add lifecycle dispatching to Stack push in nav.go: after push, call Disappeared() on old top screen and Appeared() on new top screen (via LifecycleScreen type assertion; skip if screen does not implement it)
- [x] T026 [US3] Add lifecycle dispatching to Stack pop in nav.go: call Disappeared() on popped screen and Appeared() on revealed screen (via LifecycleScreen type assertion)
- [x] T027 [US3] Add lifecycle dispatching to Stack replace in nav.go: call Disappeared() on old screen and Appeared() on new screen (via LifecycleScreen type assertion)
- [x] T028 [US3] Implement re-entrant stack modification prevention in nav.go: add inLifecycle bool and pendingOps []tea.Msg fields to Stack; queue navigation messages received during lifecycle processing and replay after completion
- [x] T029 [P] [US3] Write table-driven tests for lifecycle events: push fires Disappeared+Appeared, pop fires Disappeared+Appeared, replace fires Disappeared+Appeared, non-LifecycleScreen screens are silently skipped in lifecycle_test.go
- [x] T030 [P] [US3] Write test for re-entrant stack modification: push during Appeared handler is queued and executed after lifecycle completes in lifecycle_test.go

**Checkpoint**: Lifecycle events fire correctly for all stack operations. Re-entrant modifications are safely queued.

---

## Phase 6: User Story 4 - Combined Stack and Focus (Priority: P4)

**Goal**: Verify that per-screen focus state is preserved independently across push/pop cycles.

**Independent Test**: Create two screens each with a FocusManager and focusable items, push Screen B, change focus within B, pop B, verify Screen A's focus is unchanged.

### Implementation for User Story 4

- [x] T031 [P] [US4] Write integration test: Screen A has focus on item 2, push Screen B, pop Screen B, verify Screen A's FocusManager still reports FocusedIndex() == 1 in nav_test.go
- [x] T032 [P] [US4] Write integration test: push Screen B, move focus to item 3, pop Screen B, push a new Screen B instance, verify new Screen B has default focus (index 0) in nav_test.go

**Checkpoint**: Stack and FocusManager compose correctly. Per-screen focus state is independent and preserved.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, and final validation

- [x] T033 [P] Add Go doc comments to all exported symbols (interfaces, structs, functions, methods, message types) in nav.go, focus.go, lifecycle.go per Constitution Principle VI
- [x] T034 [P] Create basic example application (2-screen push/pop with 2 focusable inputs per screen) in _examples/basic/main.go following quickstart.md patterns
- [x] T035 Run `go vet ./...` and `go test ./...` to verify zero warnings, all tests pass, and all exported APIs have coverage
- [x] T036 Validate quickstart.md code snippets compile and match the actual API

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Foundational — no other story dependencies
- **US2 (Phase 4)**: Depends on Foundational — no dependency on US1; **can run in parallel with US1**
- **US3 (Phase 5)**: Depends on US1 (lifecycle dispatching modifies Stack push/pop/replace)
- **US4 (Phase 6)**: Depends on US1 + US2 (integration of both)
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (P1)**: After Foundational — independent, MVP
- **US2 (P2)**: After Foundational — independent of US1, can run in parallel
- **US3 (P3)**: After US1 — adds lifecycle dispatching to existing Stack operations
- **US4 (P4)**: After US1 + US2 — integration verification only

### Within Each User Story

- Interfaces/types before implementation
- Core logic before edge cases
- Implementation before tests (Constitution allows alongside/after)
- All tests passing before moving to next story

### Parallel Opportunities

- T003, T004, T005 (foundational interfaces — different files)
- T014, T015 (US1 tests — different test concerns)
- US1 and US2 can run in parallel after Foundational phase
- T022, T023 (US2 tests — different test concerns)
- T029, T030 (US3 tests — different test concerns)
- T031, T032 (US4 integration tests — different scenarios)
- T033, T034 (polish — doc comments vs example app)

---

## Parallel Example: After Foundational Phase

```bash
# US1 and US2 can start simultaneously:
# Stream 1 (US1 - Stack):
Task: "T006 [US1] Implement NewStack constructor in nav.go"
Task: "T007 [US1] Implement Stack.Init() in nav.go"
# ...through T015

# Stream 2 (US2 - FocusManager, parallel with Stream 1):
Task: "T016 [US2] Implement NewFocusManager constructor in focus.go"
Task: "T017 [US2] Implement FocusManager.Update() for Tab/Shift+Tab in focus.go"
# ...through T024
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T005)
3. Complete Phase 3: User Story 1 (T006-T015)
4. **STOP and VALIDATE**: `go test ./...` passes, push/pop/replace work
5. Library is usable for single-focus-per-screen navigation

### Incremental Delivery

1. Setup + Foundational → interfaces and types ready
2. US1 → Stack navigation works → MVP
3. US2 → FocusManager works → screens can have multiple interactive components
4. US3 → Lifecycle events fire → screens can manage resources
5. US4 → Integration verified → full feature set
6. Polish → docs, examples, validation → release-ready

### Parallel Delivery

1. Complete Setup + Foundational together
2. Once Foundational is done:
   - Stream A: US1 (Stack in nav.go)
   - Stream B: US2 (FocusManager in focus.go)
3. After Stream A completes: US3 (lifecycle in nav.go)
4. After both streams complete: US4 (integration tests)
5. Polish phase

---

## Notes

- [P] tasks = different files or independent test concerns, no dependencies
- [Story] label maps task to specific user story for traceability
- All tasks include exact file paths from plan.md project structure
- Constitution Principle IV requires tests for all public APIs
- Constitution Principle II: no goroutines or I/O outside tea.Cmd
- Constitution Principle VI: String() on Stack and FocusManager for observability
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
