---

description: "Task list for Focus Message Routing implementation"
---

# Tasks: Focus Message Routing

**Input**: Design documents from `/specs/002-focus-message-routing/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/api.go.txt

**Tests**: Required by Constitution Principle IV (Test-Required). Tests are written alongside implementation within each user story phase.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- Single Go module at repository root: `focus.go`, `focus_test.go`
- Test helpers updated in `focus_test.go` (mockFocusable, mockBoundedFocusable)
- No new source files

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: Modify the Focusable interface and update test helpers so all user stories can proceed

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T001 Modify Focusable interface to embed tea.Model in focus.go: add `tea.Model` embedding so Focusable requires Init, Update, and View methods. Update the doc comment to reflect the breaking change.
- [x] T002 [P] Update test helpers in focus_test.go: add Init() tea.Cmd, Update(tea.Msg) (tea.Model, tea.Cmd), and View() string methods to mockFocusable (with messages []tea.Msg field for recording routed messages). mockBoundedFocusable inherits via embedding. Verify existing focus tests still pass.

**Checkpoint**: Focusable embeds tea.Model. All existing tests compile and pass. User story implementation can begin.

---

## Phase 2: User Story 1 - Automatic Message Routing (Priority: P1) MVP

**Goal**: FocusManager automatically forwards non-focus messages to the currently focused component via tea.Model.Update, returning commands and retaining updated state.

**Independent Test**: Create a FocusManager with three mock components, send a non-focus key message, verify only the focused component's messages slice contains the message and its updated state is retained.

### Implementation for User Story 1

- [x] T003 [US1] Implement message routing in FocusManager.Update() in focus.go: add a routeMessage(msg) tea.Cmd helper that calls item.Update(msg) on the focused item, asserts the returned tea.Model back to Focusable, stores it in the items slice, and returns the command. Modify the default path at the end of Update to call routeMessage instead of returning nil. Guard against focusIndex == -1.
- [x] T004 [US1] Write table-driven tests for automatic message routing in focus_test.go: (a) non-focus key routed only to focused item (check messages slice), (b) command from focused item's Update is returned by FocusManager, (c) updated tea.Model is retained (state change persists after routing), (d) unfocused items do NOT receive the message
- [x] T005 [US1] Write edge case tests in focus_test.go: (a) no routing when focusIndex is -1 (use FocusIndex(-1) first), (b) no routing when items list is empty, (c) routed message causing focused item to return a navigation command (e.g., PushMsg) — command is returned unmodified by FocusManager

**Checkpoint**: Non-focus messages automatically route to the focused component via tea.Model.Update. MVP complete.

---

## Phase 3: User Story 2 - Keyboard Focus Cycling Consumed (Priority: P2)

**Goal**: Tab/Shift+Tab are consumed by FocusManager (not forwarded to any component). Mouse clicks that target a bounded component move focus first, then forward the click to the newly focused component.

**Independent Test**: Create a FocusManager with two mock components, send Tab, verify focus moved but neither component's messages slice contains the Tab. Then test mouse click on bounded component: verify focus moved AND click appears in the newly focused component's messages.

### Implementation for User Story 2

- [x] T006 [US2] Modify mouse click handling in FocusManager.Update() in focus.go: (a) when click moves focus to a different bounded item, call setFocus then routeMessage, combine commands with tea.Batch; (b) when click targets the already-focused bounded item, call routeMessage instead of returning nil; (c) remove the early return for non-press mouse events so they fall through to default routing
- [x] T007 [P] [US2] Write tests for Tab/Shift+Tab consumption in focus_test.go: send Tab and Shift+Tab to FocusManager with mock items, verify focus changed but NO item's messages slice contains the Tab/Shift+Tab message
- [x] T008 [P] [US2] Write tests for mouse click forwarding in focus_test.go: (a) click on bounded item at different index moves focus then forwards click (check messages and combined commands), (b) click on already-focused bounded item forwards click without focus change, (c) click outside all bounds routes to focused item via default path

**Checkpoint**: Focus cycling keys consumed. Mouse clicks forwarded after focus change. All commands combined.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Final validation

- [x] T009 Run `go vet ./...` and `go test ./...` to verify zero warnings, all tests pass, and all exported APIs have coverage
- [x] T010 Validate quickstart.md code patterns match the actual API (Focusable interface with tea.Model embedding, FocusManager.Update routing behavior, adapter pattern)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: No dependencies — start immediately. BLOCKS all user stories.
- **US1 (Phase 2)**: Depends on Foundational — core routing logic
- **US2 (Phase 3)**: Depends on US1 — modifies mouse handling to use routeMessage from US1
- **Polish (Phase 4)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (P1)**: After Foundational — independent, MVP
- **US2 (P2)**: After US1 — extends routing to mouse click path

### Within Each User Story

- Interface changes before implementation
- Core logic before edge cases
- Implementation before tests (Constitution allows alongside/after)
- All tests passing before moving to next story

### Parallel Opportunities

- T001, T002 (foundational — different files: focus.go vs focus_test.go)
- T007, T008 (US2 tests — different test concerns, same file but independent)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Foundational (T001-T002)
2. Complete Phase 2: User Story 1 (T003-T005)
3. **STOP and VALIDATE**: `go test ./...` passes, non-focus messages route to focused item
4. FocusManager immediately eliminates dispatch boilerplate

### Incremental Delivery

1. Foundational → Focusable embeds tea.Model, test helpers updated
2. US1 → Automatic message routing works → MVP
3. US2 → Mouse click forwarding after focus change → full routing semantics
4. Polish → validation → release-ready

---

## Notes

- [P] tasks = different files or independent test concerns, no dependencies
- [Story] label maps task to specific user story for traceability
- Primary files modified: focus.go and focus_test.go
- Breaking change: Focusable now embeds tea.Model (spec permits this)
- Existing mockFocusable/mockBoundedFocusable gain tea.Model methods via T002
- Existing focus tests should pass after T002 (routing returns nil cmd from mock)
- Constitution Principle IV requires tests for all public APIs
- Constitution Principle II: no goroutines or I/O outside tea.Cmd
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
