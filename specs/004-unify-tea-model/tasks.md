# Tasks: Unify on tea.Model

**Input**: Design documents from `/specs/004-unify-tea-model/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story. US1 replaces Screen with tea.Model in the navigation stack. US2 replaces Focusable with tea.Model + FocusMsg/BlurMsg in the focus manager. US3 updates integration points that depend on both US1 and US2, plus documentation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: Add new message types that US2 depends on. No breaking changes.

- [x] T001 Add `FocusMsg` and `BlurMsg` empty struct types to focus.go with doc comments explaining they are delivered to a model's `Update` when it gains or loses keyboard focus.

**Checkpoint**: New message types available. No existing behavior changed.

---

## Phase 2: User Story 1 - Use tea.Model for Navigation Screens (Priority: P1)

**Goal**: The navigation stack accepts `tea.Model` instead of `Screen` for all operations. The `Screen` interface is removed.

**Independent Test**: Create a standard `tea.Model` (not implementing the old `Screen`) and pass it to `NewStack`. Push, pop, and replace operations work correctly.

### Implementation for User Story 1

- [x] T002 [US1] Rewrite nav.go: remove the `Screen` interface definition. Change `PushMsg.Screen` and `ReplaceMsg.Screen` field types from `Screen` to `tea.Model`. Change `Push()` and `Replace()` parameter types from `Screen` to `tea.Model`. Change `Stack.screens` from `[]Screen` to `[]tea.Model`. Change `NewStack()` parameter from `Screen` to `tea.Model`. Update `handleNav` to store updated `tea.Model` values returned by `Update` (no type assertion needed). Update `Update()` to store the `tea.Model` returned by the active screen's `Update`. Keep all other logic unchanged (inLifecycle guard, pendingOps, Init calls).
- [x] T003 [US1] Rewrite test mocks in nav_test.go: change `mockScreen.Update` return type from `(Screen, tea.Cmd)` to `(tea.Model, tea.Cmd)`. Change `recordingScreen.Update` return type similarly. Change `focusScreen.Update` return type similarly (keep internal FocusManager usage as-is since Focusable still exists). Update `newFocusScreen` to keep working with current Focusable interface.
- [x] T004 [P] [US1] Rewrite test mocks in lifecycle_test.go: change `lifecycleScreen.Update` return type from `(Screen, tea.Cmd)` to `(tea.Model, tea.Cmd)`. Change `appearStateScreen.Update` and `disappearStateScreen.Update` return types similarly.
- [x] T005 [P] [US1] Update examples/basic/main.go: change `homeScreen.Update` and `detailScreen.Update` return types from `(nav.Screen, tea.Cmd)` to `(tea.Model, tea.Cmd)`. Add `tea` import alias if needed. No other changes required.
- [x] T006 [US1] Run `go test ./...` to verify all navigation and lifecycle tests pass with tea.Model.

**Checkpoint**: Screen interface removed. Stack accepts tea.Model. All navigation tests pass. Focus tests still pass (Focusable unchanged).

---

## Phase 3: User Story 2 - Use tea.Model for Focusable Components (Priority: P2)

**Goal**: The focus manager accepts `tea.Model` instead of `Focusable`. Focus/blur is delivered via `FocusMsg`/`BlurMsg` messages through `Update`. The `Focusable` interface is removed.

**Independent Test**: Create a standard `tea.Model` that handles `FocusMsg`/`BlurMsg` in `Update`. Pass it to `NewFocusManager`. Verify Tab/Shift+Tab delivers correct messages.

### Implementation for User Story 2

- [x] T007 [US2] Rewrite focus.go: remove the `Focusable` interface definition. Change `FocusManager.items` from `[]Focusable` to `[]tea.Model`. Change `NewFocusManager` signature from `(items ...Focusable) FocusManager` to `(items ...tea.Model) (FocusManager, tea.Cmd)` — deliver `FocusMsg` to first item via `Update`, store returned model, return command. Change `SetItems` signature from `(items ...Focusable) FocusManager` to `(items ...tea.Model) (FocusManager, tea.Cmd)` — deliver `BlurMsg` to old focused item and `FocusMsg` to new first item via `Update`, collect commands. Rewrite `setFocusWithPrev` to deliver `BlurMsg` via `Update` (store returned model, collect cmd) and `FocusMsg` via `Update` (store returned model, collect cmd) instead of calling `Blur()` and `Focus()`. Maintain blur-before-focus ordering per FR-010. Update `routeMessage` to remove `.(Focusable)` type assertion — store `tea.Model` directly. Update `FocusIndex(-1)` blur-all path to deliver `BlurMsg` via `Update`.
- [x] T008 [US2] Rewrite focus_test.go: remove `Focus()` and `Blur()` methods from `mockFocusable`. Add `FocusMsg`/`BlurMsg` handling in `mockFocusable.Update` to set `focused` field. Remove `Focus()` and `Blur()` from `mockBoundedFocusable` (inherits from mockFocusable). Remove `Focus()` method from `mockFocusableWithCmd`. Update `NewFocusManager` calls throughout to handle the new `(FocusManager, tea.Cmd)` return. Update `SetItems` calls to handle the new `(FocusManager, tea.Cmd)` return. Verify all existing test assertions still pass — focus state checked via `focused` field set by `Update`.
- [x] T009 [P] [US2] Update examples/focus-routing/main.go: remove `Focus()` and `Blur()` methods from `textField`. Add `FocusMsg`/`BlurMsg` handling in `textField.Update` to set `focused` field. Update `newFormModel` to handle `(FocusManager, tea.Cmd)` return from `NewFocusManager`.
- [x] T010 [US2] Run `go test ./...` to verify all focus tests pass with message-based delivery.

**Checkpoint**: Focusable interface removed. FocusManager uses tea.Model + messages. All focus tests pass.

---

## Phase 4: User Story 3 - Integration & Cleanup (Priority: P3)

**Goal**: Update integration points that depend on both US1 (tea.Model navigation) and US2 (tea.Model focus). Update documentation.

**Independent Test**: Verify `Screen` and `Focusable` no longer exist as exported types. All examples compile. All tests pass.

### Implementation for User Story 3

- [x] T011 [US3] Update `focusScreen` helper and `newFocusScreen` in nav_test.go: change `NewFocusManager` call to pass `tea.Model` items (not Focusable) and handle `(FocusManager, tea.Cmd)` return. The mockFocusable items already implement tea.Model, so pass them directly. Update the focusables slice type from `[]Focusable` to `[]tea.Model`.
- [x] T012 [P] [US3] Update examples/focus-nav/main.go: change `listScreen.Update` and `editScreen.Update` return types from `(nav.Screen, tea.Cmd)` to `(tea.Model, tea.Cmd)`. Remove `Focus()` and `Blur()` methods from `textField`. Add `FocusMsg`/`BlurMsg` handling in `textField.Update`. Update `newEditScreen` to handle `(FocusManager, tea.Cmd)` return from `NewFocusManager`. Update `nav.Push()` call (already accepts tea.Model after US1).
- [x] T013 [P] [US3] Update doc.go: change the package documentation example from `nav.Screen` return type to `tea.Model` return type in the `Update` method examples.
- [x] T014 [US3] Run `go test ./...` to verify all tests pass across all packages.

**Checkpoint**: All integration points updated. Screen and Focusable interfaces fully removed from the codebase.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final verification across all files and examples.

- [x] T015 Run `go vet ./...` to verify no vet warnings.
- [x] T016 Run `go build ./examples/...` to verify all examples compile.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: No dependencies — start immediately
- **US1 (Phase 2)**: Depends on Phase 1 only (FocusMsg/BlurMsg types exist but are unused until US2)
- **US2 (Phase 3)**: Depends on Phase 1 (needs FocusMsg/BlurMsg types). Independent of US1 in theory, but sequencing after US1 avoids merge conflicts in shared test files.
- **US3 (Phase 4)**: Depends on both US1 and US2 (integration points use both tea.Model navigation and tea.Model focus)
- **Polish (Phase 5)**: Depends on all previous phases

### User Story Dependencies

- **US1 (P1)**: Depends on Foundational only. Modifies nav.go, nav_test.go, lifecycle_test.go, examples/basic.
- **US2 (P2)**: Depends on Foundational only (technically). Modifies focus.go, focus_test.go, examples/focus-routing. Sequenced after US1 to avoid conflicts in nav_test.go.
- **US3 (P3)**: Depends on US1 + US2. Modifies nav_test.go (focusScreen), examples/focus-nav, doc.go.

### Within Each User Story

- Core implementation task (T002, T007) must complete before test mock updates
- Test mock updates before test verification
- Example updates can run in parallel with test mock updates (different files)

### Parallel Opportunities

- T004 and T005 can run in parallel (different files, both depend on T002)
- T009 can run in parallel with T008 (different files, both depend on T007)
- T012 and T013 can run in parallel (different files, both depend on T011)
- T015 and T016 can run in parallel (read-only verification)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001: Add FocusMsg/BlurMsg types
2. Complete T002-T006: Replace Screen with tea.Model in navigation
3. **STOP and VALIDATE**: `go test ./...` — all navigation tests pass

### Incremental Delivery

1. T001 → Foundational ready
2. T002-T006 → US1 complete (navigation uses tea.Model)
3. T007-T010 → US2 complete (focus uses tea.Model + messages)
4. T011-T014 → US3 complete (integration updated, interfaces removed)
5. T015-T016 → Full verification

---

## Notes

- T002 and T007 are the pivotal tasks — they rewrite the core Stack and FocusManager implementations respectively.
- The `Bounded` interface is unchanged throughout — it is orthogonal and checked via type assertion on `tea.Model`.
- `FocusChangedMsg` behavior is unchanged.
- `NewFocusManager` and `SetItems` signatures change to return `(FocusManager, tea.Cmd)` — this is a breaking change tracked in research.md R3.
- For `FocusIndex(-1)` (blur all), the old focused item must receive `BlurMsg` via `Update` instead of having `Blur()` called directly.
- The nav_test.go `focusScreen` helper is split across US1 (Update return type) and US3 (FocusManager API change) because it depends on both interfaces.
