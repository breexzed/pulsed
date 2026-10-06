# BLUEPRINT: Pulsed

## How to read this file
You can read this file with no other context. Every term is defined here.

Two tags show who said what:
- [OWNER]: the owner said it.
- [PROPOSED]: Claude suggested it. The owner has not confirmed it. Do not treat it as settled.

This file is a plan. It is not a claim that Pulsed exists. Almost nothing is built. A part counts as built only when `docs/epistemic/ARCHITECTURE.md` lists it.

`docs/epistemic/CANON.md` holds the project's ground rules and the full word list. Read it first if you can.

## What Pulsed is
Pulsed is a Linux program, written in Go. It shows which processes run on the computer, and how much CPU and memory each one uses.

It works in four steps:
1. It reads the `/proc` folder once per second. Linux exposes process data there as files. The file `/proc/1234/stat` holds the state and CPU use of process 1234.
2. It compares each read to the one before. It calls each difference an event, for example "process 1234 started".
3. It writes the events to a log file.
4. It shows the processes as a table in the terminal. The user can select a process and kill it, if they own it.

## Why it exists
[OWNER] The owner is a beginner. They want to become a kernel-native engineer, which means writing programs that run inside the Linux kernel using eBPF. Pulsed is practice for phase 1 of that plan: learn `/proc`, process state, and Go.

The owner writes the core code. The AI helps under rules. The finished program matters less than what the owner learns.

## Words used in this file
- **Sample**: one full read of `/proc` for every process at one moment.
- **Event**: a difference between two samples. Example: PID 500 is in sample 1 but not in sample 2. That is an "exited" event.
- **Identity**: a process's PID plus its start time. The kernel reuses PIDs, so the PID alone is not enough.
- **Source**: a Go interface the sampler reads from. The real one reads `/proc`. A fake one returns fixed data, for tests.
- **Learn list**: code the owner writes. The AI may only explain, hint, and review it.
- **Ship list**: code the AI builds.
- **Gate**: four checks a feature must pass before any document says it exists. It is written as code. Something calls it when the program runs. A test passes. Someone ran the program and saw it work.
- **ADR**: a short note that records one decision and why it was made.

## Requirements
| ID | The program must | Tag |
|---|---|---|
| FR-1 | read process data from `/proc` with no root | [OWNER] |
| FR-2 | sample all processes on an interval, default 1 second | [OWNER], default [PROPOSED] |
| FR-3 | identify a process by PID plus start time | [PROPOSED] |
| FR-4 | find events by comparing two samples | [OWNER] |
| FR-5 | report these events: started, exited, state changed, CPU spike, thread count changed | [PROPOSED] |
| FR-6 | write events as JSON lines to a file, using Go's `log/slog` | [OWNER], format [PROPOSED] |
| FR-7 | show processes sorted by CPU use, in the layout in `docs/TUI_DESIGN.md` | [OWNER] |
| FR-8 | show details for a process, and kill the user's own processes after a confirmation | [OWNER] |
| FR-9 | support mouse click, wheel, and drag, and resize with the terminal. Every mouse action has a key | [OWNER], keys [PROPOSED] |
| FR-10 | build and test on Linux only | [PROPOSED] |

## The parts of the program
Each part is a folder under `internal/`. None exists yet.

| Part | What it does | Written by |
|---|---|---|
| `proc` | reads and parses files from `/proc` | owner |
| `model` | defines the shared data types: Sample, Event, Identity | owner |
| `sample` | runs the read on a timer and sends samples on a channel | owner |
| `event` | compares two samples and returns events | owner |
| `obs` | writes events to the log file | AI |
| `tui` | draws the terminal screen | AI |
| `cmd/pulsed` | starts the program and connects the parts | AI |

Import rules [PROPOSED]. A rule says which part may use which. The reason: parts that depend on few things are easier to test.
- `model` uses only the Go standard library.
- `proc`, `sample`, and `event` use `model` and the standard library.
- `tui` uses `model` and Bubble Tea. It never uses `proc`.
- Only `cmd/pulsed` uses every part.

## The log line
The log is a contract with the outside world. One JSON object per line:

```json
{"time":"2026-10-04T09:12:03Z","level":"INFO","msg":"event","seq":42,"kind":"state_changed","pid":1234,"start_ticks":98765,"name":"bash","from":"S","to":"R"}
```

- `kind` is one of: `started`, `exited`, `state_changed`, `cpu_spike`, `threads_changed`.
- The log never holds a command line. Command lines can hold passwords.
- The program never prints logs to the terminal while the screen runs. It writes to the file.
- Default path: `$XDG_STATE_HOME/pulsed/pulsed.log`. If that variable is unset, `~/.local/state/pulsed/pulsed.log`.

All of this is [PROPOSED] except the JSON-lines choice. The owner designs the internal Go types in steps G1 to G3.

## The repo grows in stages
The old plan listed about 40 files for 4 parts that do not exist. That was too much. The owner would spend their attention reviewing files, not learning Go.

New rule: **a file is added only when its trigger happens.** If an AI proposes a file, ask which trigger happened.

### Stage 1: now
```
pulsed/
├── internal/
│   └── proc/                 owner writes this. ListPIDs goes here.
├── learning/
│   ├── check.py
│   ├── verify.py
│   ├── meta.md
│   ├── ledger.md
│   ├── experiments.md
│   └── misconceptions.md
├── docs/
│   ├── BLUEPRINT.md
│   ├── TUI_DESIGN.md         kept. Not used until step UI1.
│   └── epistemic/
│       ├── CANON.md
│       ├── QUESTIONS.md
│       ├── ARCHITECTURE.md
│       └── (5 more files, kept, not read yet:
│            METHOD, CURRENT_MODEL, CLAIMS,
│            AGENT_PROTOCOL, EPISTEMIC_CONTROL)
├── .opencode                 configurations for opencode
├── AGENTS.md
├── README.md
├── LICENSE                    MIT
├── .gitignore
└── go.mod
```

Add one line at the top of each Go file: `// SPDX-License-Identifier: MIT`. A script can check it later.

### Stage 2: add when it hurts
| File or folder | Add it when |
|---|---|
| `cmd/pulsed/main.go` | step G1: you first need to run something |
| `learning/maps/` | step G1: you write the first map file |
| `.github/workflows/ci.yml` (only `go test`) | `ListPIDs` has its first test |
| `internal/model/` | two parts need the same type |
| `internal/sample/` | step G4 or G5 |
| `internal/event/` | step G6 |
| `internal/obs/` | step G6 |
| `docs/reference/log-schema.md` | step G6 |
| `docs/decisions/` | the first decision worth an ADR. Likely the channel design in G5 |
| `internal/tui/` | after G4, for UI1 |
| `Makefile` | you type the same three commands twice in a day |
| `.golangci.yml` and a lint job in CI | after G3 |
| `.githooks/` | after your first bad commit, such as a forgotten `gofmt` |
| `learning/overrides.md` | the first "ship it" on a learn-list file |
| `learning/disagreements.md` | the first time you overrule the AI |
| `tools/ste_lint.py`, `ste_words.tsv` | the README and BLUEPRINT stop changing daily |
| `tools/diagram.py`, `page.py` | you have 3 map files |
| `opencode.json`, `.editorconfig` | an editor or tool needs them |

### Stage 3: before the first public release
`NOTICE`, `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, `CHANGELOG.md`, `CODEOWNERS`, issue and PR templates, `dependabot.yml`, `release.yml`, `tools/spdx_check.py`, and `docs/assets/`.


## How to build it
Until a `Makefile` exists, the quality check is three commands, run from the repo root:
```
gofmt -s -l .
go vet ./...
go test ./...
```
`gofmt` should print nothing. The other two should pass.

### The loop for every learn-list step
1. Write a map of the piece in 5 lines. Answer: what starts it? What does it read? What state does it hold? Who uses its output? What happens when it fails?
2. Write your prediction.
3. Run a small experiment. Record the result in `learning/experiments.md`.
4. Write the code and one test. Commit your first try with the prefix `attempt:`. First tries are meant to be bad.
5. If a bug cost more than 15 minutes, add an entry to `learning/misconceptions.md`. The number 15 is [PROPOSED].
6. Run `python3 learning/check.py`. Then propose a ledger entry. It is `passed` only after a cold artifact.

### The steps
The owner first target is G0 and G1. Do not measure ownere against the whole table.

| Step | You do | Done when | Closes |
|---|---|---|---|
| G0 | Make `ListPIDs` take a folder path and test it with a temp folder | The test passes with fake numbered folders and one non-numbered folder | none |
| G1 | Read one process from `/proc/self/stat` | Owner predicted the fields first, and the program prints them | none |
| G2 | Parse the process name field | A test passes for a name like `a b)` | none |
| G3 | Compute CPU percent | OWner numbers are close to `top` for a busy loop | Q3 |
| G4 | Write the `Source` interface and a fake | A test runs against the fake | Q8 |
| UI1 | AI builds a static screen with fake data | The screen draws and resizes | Q5, Q9 |
| G5 | Goroutine and channel for the sampler | A slow reader does not stop the sampler | Q4 |
| G6 | Find events and write the log | The log shows start and exit events for a test process | Q1, Q2 |
| G7 | Sort processes, then read threads | Sorted output is correct in a test | none |
| G8 | Kill action | A test process gets killed after confirmation | none |
| UI2 to UI5 | AI builds scroll, mouse, drag, and the kill dialog | Each passes the tests in `docs/TUI_DESIGN.md` | Q7 |

The Q numbers are experiments listed in `docs/epistemic/QUESTIONS.md`.

### Reference material
- `man 5 proc` is the main reference.
- The Go library `github.com/c9s/goprocinfo` is the second. Owner reads it after own attempt, never before. Do not import it. The core uses only the Go standard library. If you copy a line, keep its MIT license notice.

## Open-source plan
Applies at stage 3.
- License: MIT [OWNER].
- Commit messages: Conventional Commits, such as `feat(proc): read stat fields`.
- The `main` branch is protected. Pull requests need passing checks. History stays linear.
- The version stays `v0.x` until phase 1 is done. Tag `v0.1.0` then.
- Changelog: Keep a Changelog format, with SemVer.

### CI later
`ci.yml` runs on every push, on `ubuntu-latest`. The jobs: build and test with `-race`; `gofmt`, `go vet`, and `golangci-lint`; `govulncheck`; doc style with `ste_lint.py`; tests for `tools/`. The `learning/check.py` job never blocks a merge.

`release.yml` runs on a tag that starts with `v`. It builds `linux/amd64` and `linux/arm64`, writes a `sha256sums` file, and creates a GitHub Release from the changelog. A release plus `go install` is the only deployment.

## Phase 1 is done when
- The screen shows live data, and the log records events, end to end.
- Every part in this file is listed in `docs/epistemic/ARCHITECTURE.md` with a cleared gate.
- The three quality commands pass on a fresh clone.
- `python3 learning/check.py` reports clean.

## Risks
- Free models change or stop without warning. `AGENTS.md` names a fallback model. Experiment Q6 measures this.
- Mouse support differs by terminal. Keys cover every action. Experiment Q7 checks it.
- The screen design is larger than the core. The UI steps stay behind the core steps.
- Overuse of "ship it" on learn-list files. `learning/overrides.md` counts each use, once it exists.
- `learning/` is public if the repo is public. It shows your ledger and your experiments. Decide if you want that.

## Open decisions for the owner
- Where does `ListPIDs` live today? Move it into `internal/proc`.
- Is the stage 1 tree small enough, or still too big?
- Keep or change the 15 minute rule in the loop.