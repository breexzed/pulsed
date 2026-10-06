# Pulsed system diagram

Pulsed is a Linux program that watches processes. It reads files the kernel
already keeps about every running process, notices what changed, writes those
changes to a log file, and shows a live table in the terminal.

This diagram reads from the bottom up: the kernel is the ground, the face is
what the user sees.

```mermaid
flowchart TD
    subgraph K["Ground — the Linux kernel"]
        P0["/proc — one folder per process,<br/>stat, status, threads list and more.<br/>Read about once per second."]
    end

    subgraph D["Reading and comparing"]
        R["proc — reads and parses /proc files"]
        T["model — the shared types:<br/>Sample, Event, Identity"]
        S["sample — runs the read on a timer,<br/>hands samples over a channel"]
        E["event — compares two samples,<br/>produces started / exited /<br/>state_changed / cpu_spike /<br/>threads_changed"]
    end

    subgraph W["Wiring and logging"]
        O["obs — appends events to a<br/>JSON-lines log file"]
        C["cmd/pulsed — starts everything"]
    end

    subgraph F["The face "]
        U["tui — the terminal screen:<br/>sorted process table, details,<br/>kill with confirmation,<br/>mouse and keyboard"]
    end

    P0 --> R
    R --> S
    S -->|"one Sample per tick"| E
    E -->|"events"| O
    E -->|"events"| U
    S -->|"latest Sample"| U
    C -.starts.-> S
    C -.starts.-> U
    T --- R
    T --- S
    T --- E
    T --- U
```

Three things to notice:

1. Data flows one way. Kernel → proc → sample → event → log and screen.
2. `model` sits under everything. It holds the types all parts agree on,
   so no part needs to know how another part works.
3. Each box is a separate Go package. Only the starter program talks to
   all of them; everything else talks to as few neighbors as possible.

Build order is incremental: ListPIDs first, then reading one process,
parsing, CPU percentage, timer + channel, events, log, table, kill.
The screen comes after the data path works.
