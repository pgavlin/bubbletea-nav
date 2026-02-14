# Feature Specification: Focus Message Routing

**Feature Branch**: `002-focus-message-routing`
**Created**: 2026-02-14
**Status**: Draft
**Input**: User description: "The FocusManager should handle message routing. Messages should be routed to the focused Model."

## Clarifications

### Session 2026-02-14

- Q: When a mouse click targets a bounded component and moves focus to it, should that click also be forwarded to the component as a routed message? → A: Yes — forward the click after focus change. Click moves focus, then the click is routed to the newly focused component (consistent with 001-screen-nav FR-015).
- Q: Should a new interface be introduced for message-capable components? → A: No — no new interface. Messages should be forwarded to Focusable items that also implement tea.Model, using the standard Update(tea.Msg) (tea.Model, tea.Cmd) flow.
- Q: Is backward compatibility required for focus-only items that do not implement tea.Model? → A: No — backward compatibility is not required. All Focusable items must also implement tea.Model.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Automatic Message Routing (Priority: P1)

A developer building a screen with multiple interactive components
(e.g., text inputs, select lists) currently must manually check which
component is focused and dispatch messages to it. This is repetitive
boilerplate that every screen with focus management must duplicate.
The FocusManager should automatically forward non-focus messages to
the currently focused component, eliminating this boilerplate.

Today the developer writes:

```
1. Delegate to FocusManager for Tab/Shift+Tab handling
2. Check FocusedIndex()
3. Switch on the index to route the message to the correct component
4. Collect the command from the component
```

With this feature, the developer writes:

```
1. Delegate to FocusManager — it handles focus cycling AND message routing
```

**Why this priority**: This is the core value of the feature. Without
automatic routing, the FocusManager only manages focus state but
cannot simplify the most common developer task: dispatching input to
the right component.

**Independent Test**: Can be fully tested by creating a FocusManager
with three mock components (each implementing both Focusable and
tea.Model), sending a non-focus message, and verifying that only the
focused component received the message and its updated state is
retained by the FocusManager.

**Acceptance Scenarios**:

1. **Given** a FocusManager with three components and focus on the
   second, **When** a keyboard message (not Tab/Shift+Tab) is sent,
   **Then** only the second component receives and processes the
   message.
2. **Given** a FocusManager with focus on a component, **When** the
   focused component processes a message and produces a command,
   **Then** the FocusManager returns that command to the caller.
3. **Given** a FocusManager with focus on a component, **When** the
   focused component processes a message and changes its internal
   state, **Then** the FocusManager retains the updated component
   state (not the stale version from before the message).
4. **Given** a FocusManager with three components and focus on the
   first, **When** a non-focus message is sent, **Then** the second
   and third components do not receive the message.

---

### User Story 2 - Keyboard Focus Cycling Consumed (Priority: P2)

When the FocusManager receives a keyboard focus-cycling message (Tab
or Shift+Tab), it should handle focus management and not forward that
message to the focused component. This prevents components from
reacting to Tab presses that are meant for navigation, not for input.
Mouse clicks that target a bounded component are handled differently:
focus moves first, then the click is forwarded to the component.

**Why this priority**: Without this ordering guarantee, components
could misinterpret focus-cycling keys as input, leading to unexpected
behavior (e.g., a text input inserting a tab character when the user
intends to move focus).

**Independent Test**: Can be fully tested by creating a FocusManager
with two components, sending a Tab message, and verifying that focus
moved but the (previously or newly) focused component did not receive
the Tab as a routed message.

**Acceptance Scenarios**:

1. **Given** a FocusManager with focus on the first component,
   **When** Tab is pressed, **Then** focus moves to the second
   component and neither component receives the Tab as a routed
   message.
2. **Given** a FocusManager with focus on the first component,
   **When** Shift+Tab is pressed, **Then** focus wraps to the last
   component and no component receives Shift+Tab as a routed message.
3. **Given** a FocusManager with bounded components, **When** a mouse
   click targets the second component, **Then** focus moves to the
   second component and the click is also forwarded to it as a routed
   message.

---

### Edge Cases

- What happens when the focus index is -1 (no item focused) and a
  non-focus message arrives? The message is not routed to any
  component.
- What happens when the items list is empty and a message arrives?
  The message passes through with no effect.
- What happens when a routed message causes the focused component to
  produce a navigation command (e.g., Push)? The command is returned
  normally; the FocusManager does not interpret it.
- What happens when a mouse click lands within a bounded component's
  area? Focus moves to that component and the click is then forwarded
  to the newly focused component as a routed message (consistent with
  001-screen-nav FR-015: click both moves focus and delivers the event).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The FocusManager MUST forward non-focus messages to the
  currently focused component via its tea.Model Update method.
- **FR-002**: The FocusManager MUST return any commands produced by
  the focused component's message processing.
- **FR-003**: The FocusManager MUST retain the updated tea.Model
  returned by the focused component's Update method (reflecting any
  internal state changes from processing the message).
- **FR-004**: The FocusManager MUST handle keyboard focus-cycling
  messages (Tab, Shift+Tab) without forwarding them to any component
  as routed messages.
- **FR-008**: When a mouse click targets a bounded component and
  moves focus, the FocusManager MUST first move focus and then
  forward the click to the newly focused component as a routed
  message.
- **FR-005**: The FocusManager MUST NOT route messages when no
  component has focus (focus index is -1).
- **FR-006**: The FocusManager MUST combine its own commands (from
  focus changes) with commands from routed message processing, so the
  caller receives all commands.

### Key Entities

- **Focusable Component**: An interactive component that can receive
  and lose focus and process messages. All Focusable items must also
  implement tea.Model. The FocusManager routes messages via the
  standard tea.Model Update method.

### Assumptions

- Only the focused component receives routed messages. Broadcasting
  messages to all components is out of scope.
- The FocusManager does not interpret or filter routed messages beyond
  the focus-cycling keys. All other messages are forwarded as-is.
- Components are responsible for their own message handling logic.
  The FocusManager is only a dispatcher.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can build a screen with multiple interactive
  components without writing any manual message-routing logic (no
  index checking, no switch statements for dispatch).
- **SC-002**: Message routing adds zero perceptible latency to input
  handling (operations complete in under 1ms).
- **SC-003**: All public message-routing capabilities are covered by
  automated tests with no external dependencies required to run them.
