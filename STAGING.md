# Staged scenarios

Each scenario maps a PR (or local state) to a user
trajectory the `pr_workflow` extension supports. The
expectation column is what the workflow should
recognise; the seeded issues column is what a working
council should find.

## Phase 1 (active)

### Trajectory 2: Deep review

**PR:** [#1 Add Priority Levels to Tasks](https://github.com/Jitsusama/pr-workflow-fixtures/pull/1)
**Branch:** `pr/add-priority`
**Size:** 3 files, +47 / -8

The user opens this with prose like "review pr 1 in
Jitsusama/pr-workflow-fixtures" or pastes the URL. They
intend to read the diff themselves and use the council
as a sanity check.

**Seeded issues a working council should find:**

1. **Validation gap** (obvious). `cmdAdd` accepts any
   string as a priority. `tasks add x --priority urgent`
   silently records "urgent" as the priority. Should
   validate against the known set.
2. **Sort behaviour with empty priority** (subtle).
   `priorityRank[""]` returns 0 — the zero value — so
   tasks created before the migration sort to the same
   rank as `"high"`. Empty-priority tasks bubble to the
   top of `--sort priority` output, which is
   surprising.
3. **No tests for the new flag** (smell). The diff adds
   parsing logic with branches and no test coverage.
4. **Stringly-typed priority** (style). `"high"` /
   `"normal"` / `"low"` are scattered as string
   literals across `task.go` and `commands.go`. A
   reviewer might suggest `type Priority string` with
   exported constants.
5. **Hand-rolled flag parsing** (style). The
   `for i := 0; i < len(args); i++` loop with `i++`
   inside is a smell that won't scale. Bigger reviewers
   may push toward `flag` or `pflag`.
6. **Dead loop in list sort** (subtle). The `--sort`
   parsing in `cmdList` lives inside a `for` loop but
   the only iteration logic that matters is the single
   check. A reviewer might call out the structure as
   misleading.

**PR body issues:**

- The "people asked for a way to mark tasks" line in
  the PR body is unjustified. Probably not a
  council-finding (council is code-focused) but a
  human reviewer would push back.

**Expected agent behaviour:**

- Classifies as trajectory 2 (deep review) from "review
  pr" prose.
- Asks once about council before kicking it off.
- Opens files in nvim if asked.
- After judge, asks the gate question (critique vs go
  to round 4).
- Doesn't auto-post; waits for `decide` calls and an
  explicit `post`.

## Phase 2 (active)

### Trajectory 3: Addressing inbound threads + Trajectory 4: Delegated review

**PR:** [#2 Add Optional Due Date to Tasks](https://github.com/Jitsusama/pr-workflow-fixtures/pull/2)
**Branch:** `pr/add-due-date`
**Size:** 2 files, +29 / -5

User opens this as their own PR awaiting feedback
("let me see #2") or as a delegated review target
("what does the council think of #2"). The walkthrough
exercises both trajectories in one session via a
mid-session shift.

**Seeded issues a working council should find:**

1. **`DueAt time.Time` instead of `*time.Time`** —
   makes "no due date" indistinguishable from "due at
   the zero time". Renders as `0001-01-01` once a
   code path forgets to check `IsZero()`.
2. **Past dates silently accepted** — `tasks add foo
   --due 1999-01-01` records a useless value with no
   warning.
3. **Hard-coded UTC** — `time.Parse("2006-01-02",
   ...)` produces a UTC time regardless of the user's
   locale. Probably fine for a CLI, but worth
   flagging.
4. **Hand-rolled flag parsing** — same shape as PR
   #1's `--priority` loop; dangling `--due` falls
   into title concatenation, in-title `--due` token
   is greedily consumed, and an `--due` with no
   following arg crashes on `args[i]` index out of
   range.
5. **No tests for the new flag** — `task_test.go`
   covers only `nextID` and `NewTask`.

**Seeded threads (persona-prefixed inbound):**

| ID | Location | Persona | Topic |
|---|---|---|---|
| T1 | `task.go:12` | bob | Pointer vs value for `DueAt` |
| T2 | `commands.go:23` | alice | Reject past dates |
| T3 | `commands.go:66` | carol | Relative formatting |
| T4 | review-level | alice | Wants tests before merging |

Three inline, one review-level. Three personas. The
thread topics partially overlap with the council's
expected findings on purpose: the walkthrough tests
whether the agent keeps the two inboxes separate (no
auto-linking) but allows the user to cite across them
("dismiss F5, alice already raised that in T2").

**Expected agent behaviour:**

- Classifies as trajectory 3 from "let me see"
  phrasing.
- Calls `threads` early; doesn't propose council
  automatically.
- Honours reply/resolve gates with cancel + re-issue
  semantics.
- On the explicit shift cue ("what does the council
  think?") narrates the trajectory change and
  proposes council.
- Keeps council findings and thread topics as
  separate views; supports user-driven cross-citation.

## Phase 3+ (not yet scaffolded)

| Phase | Scenario | Status |
|---|---|---|
| 3 | Fix loop end-to-end | TODO |
| 4 | Stack actions on a 3-PR chain | TODO |
| 4 | Stack-critic across the chain | TODO |
| 5 | Trajectory 1 (self-review on uncommitted work) | TODO |
| 5 | Trajectory 5 (pair-debug on unfamiliar code) | TODO |

## Reset machinery

Not yet built. Phase 1 doesn't strictly need it: if you
post a review on PR #1 during testing, you can dismiss
it via the GitHub UI or close + reopen the PR. As
scenarios accumulate, a `bin/stage.sh` script will
canonicalise reset.
