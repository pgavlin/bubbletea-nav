# Feature Specification: Screen Navigation Library

**Feature Branch**: `001-screen-nav`
**Created**: 2026-02-14
**Status**: Draft
**Input**: User description: "A navigation library for Bubble Tea that supports push/pop screen transitions and navigation between interactable user interface components within screens."

## Clarifications

### Session 2026-02-14

- Q: Can the focus order change at runtime, or is it fixed at screen creation? → A: Static per render cycle — developer can update the focus list during the update phase; changes apply on the next render.
- Q: How should the library handle mouse events? → A: Focus-aware routing — mouse clicks on a focusable component move focus to it, then deliver the click event to that component.
- Q: What lifecycle events does a replace operation fire? → A: Both — the old screen receives "disappeared" and the new screen receives "appeared".

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Screen Stack Navigation (Priority: P1)

A developer building a multi-screen terminal application needs to
navigate users between distinct screens (e.g., a list screen, a
detail screen, a settings screen). The developer pushes a new screen
onto a navigation stack when the user selects an item, and pops it
to return to the previous screen. The stack preserves each screen's
state so that returning to a prior screen resumes exactly where the
user left off.

**Why this priority**: Without a screen stack, there is no
multi-screen navigation at all. This is the foundational capability
that every other feature builds on.

**Independent Test**: Can be fully tested by creating two screens,
pushing the second onto the stack, verifying it becomes active, then
popping it and verifying the first screen is restored with its
original state.

**Acceptance Scenarios**:

1. **Given** an application with one active screen, **When** a new
   screen is pushed onto the stack, **Then** the new screen becomes
   the active screen and receives all user input.
2. **Given** a stack with two or more screens, **When** the top
   screen is popped, **Then** the previous screen becomes active and
   its state is exactly as it was before the push.
3. **Given** a stack with only one screen (the root), **When** a pop
   is requested, **Then** the pop is ignored and the root screen
   remains active.
4. **Given** a stack with multiple screens, **When** the active
   screen is replaced, **Then** the new screen takes the position of
   the old top screen without changing the rest of the stack.

---

### User Story 2 - Component Focus Navigation (Priority: P2)

A developer building a screen with multiple interactable components
(e.g., text inputs, buttons, toggles, lists) needs users to be able
to move focus between those components. The user presses Tab to move
focus to the next component and Shift+Tab to move to the previous
one. Only the focused component receives keyboard input; unfocused
components are visible but inert. The developer declares which
components are focusable and in what order.

**Why this priority**: Screens with more than one interactive
element are unusable without focus management. This is the second
most critical capability after screen navigation itself.

**Independent Test**: Can be fully tested by creating a single
screen with three focusable components, verifying that Tab cycles
focus forward, Shift+Tab cycles backward, and only the focused
component responds to input.

**Acceptance Scenarios**:

1. **Given** a screen with three focusable components and focus on
   the first, **When** the user presses Tab, **Then** focus moves to
   the second component.
2. **Given** a screen with focus on the last component, **When** the
   user presses Tab, **Then** focus wraps to the first component.
3. **Given** a screen with focus on the first component, **When**
   the user presses Shift+Tab, **Then** focus wraps to the last
   component.
4. **Given** a screen with focus on any component, **When** the user
   types text, **Then** only the focused component receives and
   processes that input.
5. **Given** a screen with a mix of focusable and non-focusable
   components, **When** the user cycles through focus, **Then**
   non-focusable components are skipped.
6. **Given** a screen with focus on the first component, **When**
   the user clicks on the third focusable component, **Then** focus
   moves to the third component and it receives the click event.
7. **Given** a screen with focusable and non-focusable components,
   **When** the user clicks on a non-focusable area, **Then** focus
   does not change and the click is delivered to the active screen.

---

### User Story 3 - Screen Lifecycle Events (Priority: P3)

A developer needs to perform setup when a screen appears and
teardown when it disappears. For example, a screen might start a
periodic data refresh when it becomes visible and stop it when
another screen is pushed on top. The library notifies each screen
when it gains or loses visibility due to stack changes.

**Why this priority**: Lifecycle events enable resource management
and dynamic behavior but are not required for basic navigation to
function.

**Independent Test**: Can be fully tested by creating two screens
that record lifecycle events, pushing and popping them, and
verifying the correct sequence of appear/disappear notifications.

**Acceptance Scenarios**:

1. **Given** a screen is pushed onto the stack, **When** it becomes
   the active screen, **Then** it receives an "appeared" event.
2. **Given** an active screen, **When** a new screen is pushed on
   top, **Then** the previously active screen receives a
   "disappeared" event.
3. **Given** a screen that was hidden by a push, **When** the screen
   above it is popped, **Then** the revealed screen receives an
   "appeared" event.
4. **Given** an active screen, **When** it is popped from the stack,
   **Then** it receives a "disappeared" event before removal.
5. **Given** an active screen, **When** it is replaced by a new
   screen, **Then** the old screen receives a "disappeared" event
   and the new screen receives an "appeared" event.

---

### User Story 4 - Combined Stack and Focus (Priority: P4)

A developer builds an application where each screen has its own
set of focusable components. When a screen is pushed, its focus
state is independent. When a screen is popped and the previous
screen is restored, the focus state of the restored screen resumes
at whichever component was last focused.

**Why this priority**: This is the integration of US1 and US2.
It is important for real-world applications but depends on both
prior stories being complete.

**Independent Test**: Can be fully tested by creating two screens
each with focusable components, pushing the second screen, changing
focus within it, popping it, and verifying the first screen's focus
is restored to its prior position.

**Acceptance Scenarios**:

1. **Given** Screen A with focus on its second component, **When**
   Screen B is pushed and then popped, **Then** Screen A's focus is
   still on its second component.
2. **Given** Screen B is pushed and the user moves focus to its
   third component, **When** Screen B is popped and later pushed
   again, **Then** Screen B starts with default focus (first
   component), not the previously focused component.

---

### Edge Cases

- What happens when a developer pushes the same screen instance
  onto the stack twice? The library treats each push as a separate
  stack entry; the screen instance appears twice in the stack.
- What happens when all focusable components on a screen are
  removed from the focus list during an update phase? On the next
  render, focus enters a "no focus" state and Tab/Shift+Tab have
  no effect until a focusable component is added to the list.
- What happens when a screen is pushed during a lifecycle event
  handler? The push is queued and processed after the current event
  completes, preventing re-entrant stack modifications.
- What happens when the developer provides zero screens at startup?
  The library requires at least one root screen; providing none is
  a usage error surfaced at initialization.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The library MUST maintain an ordered stack of screens
  where only the topmost screen is active and receives user input.
- **FR-002**: The library MUST support pushing a new screen onto the
  top of the stack.
- **FR-003**: The library MUST support popping the top screen from
  the stack, restoring the previous screen to active status.
- **FR-004**: The library MUST support replacing the top screen with
  a different screen without affecting the rest of the stack.
- **FR-005**: The library MUST ignore pop requests when only the
  root screen remains on the stack.
- **FR-006**: The library MUST preserve the full state of inactive
  screens so they can be restored exactly when they become active
  again.
- **FR-007**: The library MUST support declaring focusable
  components within a screen and managing sequential focus order.
  The focus list is static per render cycle; the developer MAY
  update it during the update phase and changes MUST take effect
  on the next render.
- **FR-008**: The library MUST move focus to the next focusable
  component when Tab is pressed and to the previous when Shift+Tab
  is pressed.
- **FR-009**: The library MUST wrap focus cyclically (last to first
  on Tab, first to last on Shift+Tab).
- **FR-010**: The library MUST skip non-focusable components during
  focus traversal.
- **FR-011**: The library MUST provide a mechanism for screens to
  identify the currently focused component so that keyboard input
  can be routed only to it. The library is not responsible for
  intercepting and forwarding all keyboard events; screens use
  FocusedIndex to determine which component should receive input.
- **FR-015**: When a mouse click targets a focusable component, the
  library MUST move focus to that component and then deliver the
  click event to it.
- **FR-016**: Mouse events that do not target a focusable component
  MUST be delivered to the active screen without changing focus.
- **FR-012**: The library MUST notify screens when they become
  visible (appeared) or hidden (disappeared) due to stack changes,
  including push, pop, and replace operations. A replace MUST fire
  "disappeared" on the old screen and "appeared" on the new screen.
- **FR-013**: The library MUST prevent re-entrant stack
  modifications during lifecycle event processing.
- **FR-014**: Each screen on the stack MUST maintain independent
  focus state that is preserved across push/pop cycles.

### Key Entities

- **Screen**: A distinct, self-contained view that occupies the full
  terminal area. Has its own state, view logic, and optionally a set
  of focusable components. Can be active (top of stack) or inactive
  (below another screen).
- **Navigation Stack**: An ordered collection of screens. The
  topmost screen is the active screen. Supports push, pop, and
  replace operations.
- **Focusable Component**: An interactive UI element within a screen
  that can receive keyboard input when focused. Has a position in
  the screen's focus order and a spatial region for mouse click
  targeting.
- **Focus State**: Tracks which component within a screen currently
  has focus. Each screen maintains its own independent focus state.
- **Lifecycle Event**: A notification sent to a screen when its
  visibility changes due to stack operations (appeared or
  disappeared).

### Assumptions

- The library targets terminal UI applications built with the
  Bubble Tea framework.
- Focus navigation uses Tab/Shift+Tab (sequential) rather than
  arrow-key-based spatial navigation. Spatial navigation may be
  added in a future iteration.
- Screens occupy the full terminal area; split-screen or overlay
  layouts are out of scope.
- Screen transitions are instantaneous; animated transitions are
  out of scope for the initial version.
- The library does not manage window resizing directly; it passes
  resize events to the active screen for handling.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can build a 3-screen application with
  push/pop navigation in under 30 minutes using the library.
- **SC-002**: Pushing and popping screens preserves prior screen
  state with 100% fidelity (no lost state on round-trip).
- **SC-003**: Focus traversal visits exactly the declared focusable
  components in order, with zero mis-deliveries of input to
  unfocused components.
- **SC-004**: Lifecycle events fire in the correct order for all
  stack operations (push, pop, replace) with no missed or duplicate
  notifications.
- **SC-005**: All public library capabilities are covered by
  automated tests with no external dependencies required to run
  them.
