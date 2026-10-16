<!-- readme-pattern: off — owner wants a short learning-project README (2026-10-17). -->
# Pulsed

**This is a learning project.** I am learning Linux internals and Go by building a process monitor from `/proc` up. It is pre-alpha. Nothing runs yet.

## What works today

- `ListPIDs` lists the process IDs in `/proc`.
- A parser for `/proc/PID/stat`. It survives process names with spaces and `)` in them, and returns an error on short lines instead of panicking.
- 125 lines of Go, 136 lines of tests. Standard library only. Tests run in CI.

## What I am testing next

The plan was simple: read `/proc` once per second, compare reads, report when a process starts or exits.

The claim that plan has to survive: can a once-per-second read see a process that lives 5 ms? If it can't, most "started" and "exited" events would never be reported. I am measuring that next. If polling fails, the fix I want to try is eBPF: let the kernel report each exec as it happens.

## How I work here

I write `internal/proc`, `internal/model`, `internal/sample` and `internal/event` by hand. AI tools may explain and review that code. They do not write it. The rules are in `AGENTS.md`.

`learning/` is my log: what I predicted, what happened, and what I got wrong.

Pull requests to those folders won't be merged. Issues and ideas are welcome.

## Run the tests

Needs Linux and Go (version in `go.mod`).

```
go test ./...
```

## License

MIT. See `LICENSE`.
