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

## Phase 2+ (not yet scaffolded)

| Phase | Scenario | Status |
|---|---|---|
| 2 | Trajectory 3 (addressing inbound threads) | TODO |
| 2 | Trajectory 4 (delegated review) | TODO |
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
