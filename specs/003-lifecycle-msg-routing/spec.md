# Feature Specification: Message-Based Lifecycle Notifications

**Feature Branch**: `003-lifecycle-msg-routing`
**Created**: 2026-02-15
**Status**: Draft
**Input**: User description: "Replace the functionality of LifecycleScreen with the currently unused ScreenAppearedMsg and ScreenDisappearedMsg. The current approach makes it difficult for a Model to change its state."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Update Screen State on Appear (Priority: P1)

A developer building a TUI app needs their screen to refresh data when it becomes visible again after being revealed by a pop. With the current approach, the `Appeared()` callback is a direct method call that cannot update the screen's own state through the normal message-processing flow. With message-based notifications, the screen receives a `ScreenAppearedMsg` through its standard `Update` method, allowing it to modify its state and return commands like any other message.

**Why this priority**: This is the core motivation for the change. The inability to update state in response to lifecycle events is the primary pain point.

**Independent Test**: Can be fully tested by pushing screen A, then pushing screen B, then popping screen B, and verifying that screen A receives a `ScreenAppearedMsg` through its `Update` method and can update its state accordingly.

**Acceptance Scenarios**:

1. **Given** a stack with screens A and B (B on top), **When** screen B is popped, **Then** screen A receives a `ScreenAppearedMsg` through its `Update` method.
2. **Given** a stack with screen A, **When** screen B is pushed, **Then** screen B receives a `ScreenAppearedMsg` through its `Update` method after initialization.
3. **Given** a stack with screen A, **When** screen A is replaced by screen B, **Then** screen B receives a `ScreenAppearedMsg` through its `Update` method after initialization.
4. **Given** a screen that modifies its own state in response to `ScreenAppearedMsg`, **When** the screen appears, **Then** the state change is reflected in subsequent `View` calls.

---

### User Story 2 - Clean Up Screen State on Disappear (Priority: P2)

A developer needs their screen to perform cleanup (e.g., stop animations, clear transient state) when it loses visibility. The screen receives a `ScreenDisappearedMsg` through its `Update` method, allowing it to update its state and optionally return commands.

**Why this priority**: Disappear notifications are the counterpart to appear notifications and complete the lifecycle story, but are used less frequently in practice.

**Independent Test**: Can be fully tested by pushing a new screen on top of an existing screen, and verifying that the previously-active screen receives a `ScreenDisappearedMsg` through its `Update` method.

**Acceptance Scenarios**:

1. **Given** a stack with screen A on top, **When** screen B is pushed, **Then** screen A receives a `ScreenDisappearedMsg` through its `Update` method.
2. **Given** a stack with screens A and B (B on top), **When** screen B is popped, **Then** screen B receives a `ScreenDisappearedMsg` through its `Update` method.
3. **Given** a stack with screen A on top, **When** screen A is replaced by screen B, **Then** screen A receives a `ScreenDisappearedMsg` through its `Update` method.
4. **Given** a screen that returns a command in response to `ScreenDisappearedMsg`, **When** the screen disappears, **Then** the returned command is executed.

---

### User Story 3 - Remove LifecycleScreen Interface (Priority: P3)

A developer currently implementing the `LifecycleScreen` interface migrates to the message-based approach. The `LifecycleScreen` interface, along with its `Appeared()` and `Blur()` methods, is removed. All screens now receive lifecycle notifications uniformly through their `Update` method without needing to implement an extra interface.

**Why this priority**: Removing the old interface simplifies the public API surface and eliminates the confusing dual-path lifecycle model, but depends on the message-based approach being in place first.

**Independent Test**: Can be tested by verifying that the `LifecycleScreen` interface no longer exists in the public API and that all lifecycle behavior works exclusively through messages.

**Acceptance Scenarios**:

1. **Given** a screen that previously implemented `LifecycleScreen`, **When** the developer removes the `Appeared`/`Disappeared` methods and instead handles `ScreenAppearedMsg`/`ScreenDisappearedMsg` in `Update`, **Then** the screen receives the same lifecycle notifications as before.
2. **Given** a screen that does not handle lifecycle messages in `Update`, **When** the screen appears or disappears, **Then** the messages are delivered and silently ignored (no errors, no panics).

---

### Edge Cases

- What happens when a screen's `Update` returns a navigation command (e.g., push or pop) in response to a lifecycle message? The stack must handle queued navigation operations correctly, as it does today.
- What happens when the root screen is the only screen on the stack? It should receive a `ScreenAppearedMsg` during initial setup but no `ScreenDisappearedMsg` (since it cannot be popped).
- What happens when `Update` for a disappeared screen returns a state change? The updated screen must be stored back into the stack so the state change persists.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The navigation stack MUST deliver a `ScreenAppearedMsg` to a screen's `Update` method when that screen becomes the active (topmost) screen due to a push, pop-reveal, or replace operation.
- **FR-002**: The navigation stack MUST deliver a `ScreenDisappearedMsg` to a screen's `Update` method when that screen loses active status due to being pushed over, popped, or replaced.
- **FR-003**: The navigation stack MUST process any commands returned by a screen's `Update` in response to lifecycle messages, including navigation commands.
- **FR-004**: The navigation stack MUST store the updated screen returned by `Update` after delivering a lifecycle message, so that state changes persist.
- **FR-005**: The `LifecycleScreen` interface MUST be removed from the public API.
- **FR-006**: For push operations, the disappeared screen MUST receive its `ScreenDisappearedMsg` before the new screen receives its `ScreenAppearedMsg`.
- **FR-007**: For pop operations, the popped screen MUST receive its `ScreenDisappearedMsg` before the revealed screen receives its `ScreenAppearedMsg`.
- **FR-008**: For replace operations, the old screen MUST receive its `ScreenDisappearedMsg` before the new screen receives its `ScreenAppearedMsg`.

### Key Entities

- **ScreenAppearedMsg**: A message delivered to a screen when it becomes the active screen. Carries no data fields.
- **ScreenDisappearedMsg**: A message delivered to a screen when it loses active status. Carries no data fields.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of existing lifecycle test scenarios pass using the message-based approach instead of the interface-based approach.
- **SC-002**: Screens can modify their own state in response to lifecycle events, verified by a test that checks state changes persist after an appear/disappear cycle.
- **SC-003**: The public API surface is reduced by one interface (`LifecycleScreen`) and two methods (`Appeared`, `Disappeared`), verified by examining the exported symbols.
- **SC-004**: All existing examples and tests compile and pass without referencing `LifecycleScreen`.

## Assumptions

- The ordering of lifecycle messages (disappear before appear) matches the current `LifecycleScreen` dispatch order and is preserved.
- The `ScreenAppearedMsg` and `ScreenDisappearedMsg` types already exist in the codebase and their definitions do not need to change.
- Navigation commands returned from lifecycle message handling are queued and processed, consistent with the existing pending-ops mechanism in the stack.
