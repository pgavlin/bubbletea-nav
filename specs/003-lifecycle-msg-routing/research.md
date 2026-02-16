# Research: Message-Based Lifecycle Notifications

**Feature**: 003-lifecycle-msg-routing
**Date**: 2026-02-15

## R1: Synchronous vs Asynchronous Lifecycle Message Delivery

**Decision**: Synchronous delivery — call `screen.Update(msg)` directly inside `handleNav`, matching the current `Appeared()`/`Disappeared()` call pattern.

**Rationale**: The current implementation calls lifecycle methods synchronously during navigation. Switching to asynchronous `tea.Cmd`-based delivery would change the ordering guarantees (messages would arrive in a later event loop tick, after `handleNav` returns). Synchronous delivery preserves the existing contract: by the time `handleNav` returns, all lifecycle notifications have been delivered and screens have updated their state.

**Alternatives considered**:
- *Return `tea.Cmd` that produces the message*: Would require the Bubble Tea runtime to deliver the message back to `Stack.Update`, which would then need to route it to the correct screen (not just the active one). Significantly more complex and changes timing semantics.
- *Batch all lifecycle messages as commands*: Same routing problem, plus ordering between disappear and appear becomes non-deterministic.

## R2: Storing Updated Screen State After Lifecycle Messages

**Decision**: Always call `screen.Update(msg)` and use the returned `Screen` value. For screens remaining in the stack (pushed-over, pop-revealed), store the updated screen back. For screens being removed (popped, replaced-away), discard the updated screen but still collect its returned command.

**Rationale**: The `Screen.Update` method returns `(Screen, tea.Cmd)`. Since screens are value types in the Bubble Tea model, state changes only persist if the returned value is stored. For screens staying in the stack, this is essential (FR-004). For screens being removed, the state is discarded anyway, but the command may trigger meaningful side effects (e.g., cleanup cancellation).

**Alternatives considered**:
- *Skip `Update` for removed screens*: Would break FR-002/FR-003 — disappeared screens must still receive the message and have their commands processed.
- *Store removed screens temporarily*: Unnecessary complexity; commands are the only meaningful output.

## R3: Handling the `inLifecycle` / Pending-Ops Mechanism

**Decision**: Maintain the existing `inLifecycle` guard and `pendingOps` queue. When `Update` is called on a screen during lifecycle message delivery and it returns a navigation command (via `tea.Cmd`), that command will be executed by the Bubble Tea runtime in a subsequent tick and hit `Stack.Update` normally. The `inLifecycle` flag prevents direct re-entrant calls during the synchronous dispatch phase.

**Rationale**: The current mechanism already handles re-entrant navigation correctly. Since lifecycle messages are delivered synchronously via direct `screen.Update()` calls (not through `Stack.Update`), any navigation messages returned as commands will flow through the normal Bubble Tea event loop and arrive at `Stack.Update` later. The `inLifecycle` guard protects against the edge case where `screen.Update()` somehow triggers a synchronous navigation message (which shouldn't happen with commands, but is defensive).

**Alternatives considered**:
- *Remove `inLifecycle` guard*: Risky — even though direct re-entrance is unlikely with the message-based approach, keeping it is low-cost defensive programming.

## R4: Impact on Existing Test Mock Pattern

**Decision**: Replace the `lifecycleScreen` test mock. The new mock will record lifecycle events in its `Update` method by matching on `ScreenAppearedMsg` and `ScreenDisappearedMsg`, rather than implementing separate `Appeared()`/`Disappeared()` methods.

**Rationale**: The test mock must match the new contract. Since `Update` now handles lifecycle messages, the mock's `Update` method needs a type switch to detect and record them. This is simpler than the old approach (one method instead of two extra interface methods).

**Alternatives considered**:
- *Generic recording screen that records all messages*: Would work but makes assertions less precise — lifecycle tests should explicitly match the message types.
