# pr-workflow-fixtures

Staged test repo for the `pr_workflow` extension in
[agentic-harness.pi](https://github.com/Jitsusama/agentic-harness.pi).

This repo exists to give the workflow real PRs to chew
on — deliberately seeded with reviewable code shapes.

**Don't use this for anything else.** Force-pushes are
expected. PRs get closed and re-opened as scaffolding
evolves.

See [`STAGING.md`](./STAGING.md) for which PR exercises
which trajectory and how to reset state between runs.

## The code

A tiny Go CLI task tracker. Add, list, done, remove.
Stdlib only. Five files, around 150 lines. Big enough to
have reviewable shape, small enough to read in five
minutes.

```
$ tasks add buy groceries
added task 1: buy groceries
$ tasks list
[ ] 1. buy groceries
$ tasks done 1
$ tasks list
[x] 1. buy groceries
$ tasks rm 1
```

## Trajectories staged

| Trajectory | PR | Status |
|---|---|---|
| 2 deep review | `pr/add-priority` | Phase 1 |
| 3 addressing feedback | (next) | Phase 2 |
| 4 delegated review | (next) | Phase 2 |
| Fix loop | (next) | Phase 3 |
| Stack actions | (next) | Phase 4 |
| 1 self-review | (next) | Phase 5 |
| 5 pair-debug | (next) | Phase 5 |

## Persona convention for inbound threads

Trajectory 3 needs review comments from "someone else".
Because we're using one GitHub account, comments meant to
be inbound are prefixed with a persona marker:

```
[alice]: This loop reads stale state. Did you mean to
re-load between iterations?
```

Markers in use:
- `[alice]` — pedantic reviewer who hunts edge cases
- `[bob]` — pragmatist who flags style and consistency
- `[carol]` — security-minded, distrusts user input

The persona is convention only. Threads are still
authored by `@Jitsusama`; the test exercises the
workflow's handling of inbound feedback, not GitHub's
identity model.
