# Feature Specification: Unify on tea.Model

**Feature Branch**: `004-unify-tea-model`
**Created**: 2026-02-15
**Status**: Draft
**Input**: User description: "Remove the Screen and Focusable interfaces in favor of tea.Model. Replace the Focus() / Blur() functionality with Focused / Blurred messages."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Use tea.Model for Navigation Screens (Priority: P1)

A developer using bubbletea-nav currently must implement the custom `Screen` interface, whose `Update` method returns `(Screen, tea.Cmd)` instead of the standard `(tea.Model, tea.Cmd)`. This means existing `tea.Model` implementations cannot be used directly with the navigation stack — they require wrapper adapters or rewriting. By replacing `Screen` with `tea.Model`, developers can pass any standard Bubble Tea model to the stack without modification.

**Why this priority**: The `Screen` interface is used by every consumer of the navigation stack. Eliminating it removes the single biggest barrier to drop-in compatibility with the Bubble Tea ecosystem.

**Independent Test**: Create a standard `tea.Model` (not implementing `Screen`) and pass it to `NewStack`. Push, pop, and replace operations work correctly, and the model's `Update` returns `(tea.Model, tea.Cmd)` as normal.

**Acceptance Scenarios**:

1. **Given** a standard `tea.Model` implementation, **When** it is passed to `NewStack`, **Then** the stack accepts it without requiring any additional interface.
2. **Given** a stack with model A, **When** model B is pushed via `Push`, **Then** model B becomes the active screen and receives messages through its standard `Update` method.
3. **Given** a stack with models A and B (B on top), **When** B is popped, **Then** A becomes active and its state is preserved.
4. **Given** a stack with model A, **When** A is replaced by model B, **Then** B becomes active and receives lifecycle messages through its standard `Update` method.
5. **Given** navigation message types (`PushMsg`, `PopMsg`, `ReplaceMsg`), **When** checked for their screen field type, **Then** they accept `tea.Model` instead of `Screen`.

---

### User Story 2 - Use tea.Model for Focusable Components (Priority: P2)

A developer building a form with focusable fields currently must implement the `Focusable` interface, which extends `tea.Model` with `Focus()` and `Blur()` methods. This creates a custom interface that existing Bubble Tea components (like `textinput.Model`) don't satisfy directly. By replacing `Focusable` with `tea.Model` and delivering focus/blur as messages, any `tea.Model` can participate in focus management by simply handling `FocusMsg` and `BlurMsg` in its `Update` method.

**Why this priority**: The `Focusable` interface is the second custom interface in the library. Removing it completes the unification on `tea.Model` and makes focus management compatible with the broader Bubble Tea component ecosystem.

**Independent Test**: Create a standard `tea.Model` that handles `FocusMsg` and `BlurMsg` in its `Update` method. Pass it to `NewFocusManager`. Verify Tab/Shift+Tab cycling delivers the correct messages, and the model can update its state in response.

**Acceptance Scenarios**:

1. **Given** a standard `tea.Model`, **When** it is passed to `NewFocusManager`, **Then** the focus manager accepts it.
2. **Given** a focus manager with models A, B, C (A focused), **When** Tab is pressed, **Then** A receives a `BlurMsg` and B receives a `FocusMsg`, both through their `Update` method.
3. **Given** a model that updates its state in response to `FocusMsg`, **When** it receives focus, **Then** the state change persists (the updated model returned by `Update` is stored).
4. **Given** a model that returns a command in response to `FocusMsg`, **When** it receives focus, **Then** the command is included in the batched output.
5. **Given** the first model in a `NewFocusManager` call, **When** the focus manager is created, **Then** the first model receives a `FocusMsg` through its `Update` method.

---

### User Story 3 - Remove Screen and Focusable Interfaces (Priority: P3)

A developer currently maintaining code that implements `Screen` or `Focusable` migrates to the unified `tea.Model` approach. The `Screen` interface and `Focusable` interface are removed from the public API. All components use `tea.Model` exclusively.

**Why this priority**: Removing the old interfaces simplifies the public API and eliminates the risk of consumers accidentally implementing the wrong interface. Depends on US1 and US2 being in place first.

**Independent Test**: Verify that `Screen` and `Focusable` no longer exist as exported types. Verify all functionality works exclusively through `tea.Model`.

**Acceptance Scenarios**:

1. **Given** a developer who previously implemented `Screen`, **When** they switch their `Update` to return `(tea.Model, tea.Cmd)` and remove any Screen-specific code, **Then** their model works with the navigation stack.
2. **Given** a developer who previously implemented `Focusable`, **When** they remove the `Focus()`/`Blur()` methods and instead handle `FocusMsg`/`BlurMsg` in `Update`, **Then** their model works with the focus manager.
3. **Given** a model that does not handle `FocusMsg` or `BlurMsg` in `Update`, **When** it receives focus or blur, **Then** the messages are delivered and silently ignored (no errors, no panics).

---

### Edge Cases

- What happens when a model's `Update` returns a navigation command in response to a focus/blur message? The focus manager should return the command normally; the stack processes it in a subsequent tick.
- What happens when `FocusManager.Update` calls a model's `Update` and the returned `tea.Model` is a different concrete type? The focus manager stores the returned model, matching standard Bubble Tea value-type semantics.
- What happens when `NewFocusManager` is called with zero items? No `FocusMsg` is delivered. Same behavior as today.
- What happens when `SetItems` is called? The previously focused item receives `BlurMsg` and the first new item receives `FocusMsg`, both via `Update`.
- What happens when the `Bounded` interface is checked? The `Bounded` interface remains unchanged — it is an optional, orthogonal interface that does not conflict with `tea.Model`.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The navigation stack MUST accept `tea.Model` (instead of `Screen`) for all operations: `NewStack`, `Push`, `Pop`, `Replace`.
- **FR-002**: The `PushMsg` and `ReplaceMsg` types MUST hold a `tea.Model` field instead of a `Screen` field.
- **FR-003**: The `FocusManager` MUST accept `tea.Model` (instead of `Focusable`) for all operations: `NewFocusManager`, `SetItems`.
- **FR-004**: The `FocusManager` MUST deliver a `FocusMsg` to a model's `Update` method when that model gains focus (replacing the `Focus()` method call).
- **FR-005**: The `FocusManager` MUST deliver a `BlurMsg` to a model's `Update` method when that model loses focus (replacing the `Blur()` method call).
- **FR-006**: The `FocusManager` MUST store the updated model returned by `Update` after delivering `FocusMsg` or `BlurMsg`.
- **FR-007**: The `FocusManager` MUST collect and return any commands from `Update` calls in response to `FocusMsg` and `BlurMsg`.
- **FR-008**: The `Screen` interface MUST be removed from the public API.
- **FR-009**: The `Focusable` interface MUST be removed from the public API.
- **FR-010**: For focus changes, the losing model MUST receive `BlurMsg` before the gaining model receives `FocusMsg`.

### Key Entities

- **FocusMsg**: A message delivered to a model's `Update` when it gains keyboard focus. Carries no data fields.
- **BlurMsg**: A message delivered to a model's `Update` when it loses keyboard focus. Carries no data fields.
- **Bounded**: An optional interface for models that support mouse-click targeting (unchanged).
- **FocusChangedMsg**: Emitted by `FocusManager` when focus changes between items (unchanged).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Any standard `tea.Model` can be used with the navigation stack and focus manager without implementing additional interfaces, verified by tests using plain `tea.Model` implementations.
- **SC-002**: 100% of existing navigation and focus test scenarios pass using the `tea.Model`-based approach.
- **SC-003**: The public API surface is reduced by two interfaces (`Screen`, `Focusable`) and two methods (`Focus`, `Blur`), verified by examining exported symbols.
- **SC-004**: All existing examples compile and function correctly after migrating from `Screen`/`Focusable` to `tea.Model`.
- **SC-005**: Models can update their own state in response to focus/blur events, verified by a test that checks state changes persist after a focus transition.

## Assumptions

- The `Bounded` interface is retained as-is. It is orthogonal to `tea.Model` and provides optional mouse-click targeting.
- The `FocusChangedMsg` type and behavior are unchanged.
- The ordering of blur-before-focus matches the current `Blur()`-before-`Focus()` dispatch order.
- Navigation message types (`PushMsg`, `ReplaceMsg`) change their field type from `Screen` to `tea.Model` — this is a breaking change, acceptable pre-v1.0.
- The `FocusMsg` and `BlurMsg` types are new additions to the public API.
