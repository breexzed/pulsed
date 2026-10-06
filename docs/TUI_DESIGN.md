# sysmon TUI Design Guideline

Version 0.1. Place at `docs/TUI_DESIGN.md`. All UI work follows this file.

## 1. Purpose

This file defines how the sysmon TUI looks, behaves, and is built.
The look comes from an engineering drawing sheet: a framed page with lettered panels and a title block.
The TUI shares its tokens and component names with the SVG diagrams in `tools/diagram.py`.

## 2. Principles

- **One sheet.** The screen is one framed sheet. Every panel has a letter and a place.
- **Addressable.** You can name any spot, for example "panel at 3B".
- **Pure rendering.** Layout and drawing are pure functions. They have no side effects.
- **Mouse and keys are equal.** Every mouse action has a key.
- **Color never works alone.** Every status has a glyph and a color.
- **Resize is normal.** The sheet fits any terminal size, at any time.

## 3. Sheet anatomy

### 3.1 Frame
- The frame is a single-line box in `overlay0`.
- Column numbers sit in the top and bottom border. Row letters sit in the left and right border.
- The grid has 8 columns and 4 row bands (A, B, C, D).
- Each column number is centered in its segment. Joints use `┬` and `┴`.
- Hide the coordinates below 100 columns. They cost too much space.

### 3.2 Row bands
| Band | Content | Height |
|---|---|---|
| A | tree, detail, states | 3 parts |
| B | process list | 4 parts |
| C | limits, events | 2 parts |
| D | title block | 4 rows, fixed |

- Divide the height left after the frame and band D in the ratio 3:4:2.
- A panel is never shorter than 5 rows or narrower than 20 columns.

### 3.3 Layout modes
| Mode | Size | Layout |
|---|---|---|
| Wide | 160 x 30 or more | A cols 1-3, B cols 4-6, C cols 7-8. D list full width. E cols 1-4, F cols 5-8. |
| Standard | 100 x 30 to 159 | A cols 1-3, B cols 4-8. D list full width. E cols 1-4, F cols 5-8. Panel C opens as an overlay with `c`. |
| Compact | 80 x 24 to 99 | One panel at a time. A tab strip lists the panels. Title block stays. |
| Too small | under 80 x 24 | Show one line: "Terminal too small. Need 80 x 24." |

- Choose the mode on every resize message.
- Never crash on any size. Test sizes from 1 x 1 up to 400 x 100.

## 4. Panel anatomy

```
┌─ A ─ Process tree ──────────── 214 processes ┐
│ content                                    ▲ │
│                                            █ │
│                                            ▼ │
└──────────────────────────────────────────────┘
```

- **Tab:** the letter, with padding. Unfocused: `surface2` background. Focused: `lavender` background, `crust` text.
- **Title:** bold, `text`.
- **Subtitle:** right-aligned, `overlay1`. It shows a count or a mode.
- **Border:** single line, `surface2`. Focused: `lavender`.
- **Scrollbar:** on the right border. Show it only when the content overflows.
- Shared borders use joints: `├ ┤ ┬ ┴ ┼`.
- Truncate long text with `…`. Measure display width, not bytes.

## 5. Panel catalog

| ID | Panel | Image source | Shows |
|---|---|---|---|
| A | Process tree | Document structure | Parent and child processes |
| B | Detail | Anatomy of sentences | The selected process, field by field |
| C | States | Verb forms | State codes and their meaning |
| D | Process list | Dictionary entries | All processes, sorted |
| E | Limits | Writing rule limits | Resource bars with limit markers |
| F | Events | History | Timeline of events |
| T | Title block | Title block | Host, kernel, sample, interval |

### A. Process tree
- Draw links with `├─`, `└─`, and `│`.
- `▸` means collapsed. `▾` means expanded.
- Key for each row: PID plus start time.
- Each row shows name, PID, and CPU%.

### B. Detail
- Show label and value pairs: PID, parent, state, threads, start time, CPU%, memory, command line.
- Under a key value, draw a bracket line in `blue`. Add a short meaning, like the annotations in the image.
- Show "Select a process" when nothing is selected.

### C. States
| Code | Meaning | Mark | Color |
|---|---|---|---|
| R | running | ✓ | green |
| S | sleeping | ✓ | overlay1 |
| D | waiting for disk | ! | peach |
| T | stopped | ! | yellow |
| Z | zombie | ✗ | red |
| I | idle | ✓ | overlay1 |

### D. Process list
- Columns: PID, NAME, STATE, CPU%, MEM, THR.
- Zebra bands use `surface0`. The selected row uses `surface1` and bold.
- Default sort: CPU% descending. The header shows `▼` or `▲` on the sorted column.
- Draw only the visible rows. Do not render 5000 rows to show 30.

### E. Limits
- One bar per resource: CPU, memory, load, process count, thread count.
- Fill uses `█` plus partial blocks `▏▎▍▌▋▊▉` for sub-cell precision.
- The track uses `surface0`.
- The limit marker is a `│` in `peach` at the limit position.
- Fill color: green below 60%, yellow below 85%, red from 85%.
- Each bar shows its value and its max on the right, like "max 20 words" in the image.

### F. Events
- One line per event: time, glyph, text.
| Glyph | Event | Color |
|---|---|---|
| ● | process started | green |
| ○ | process ended | overlay1 |
| ◆ | state changed | blue |
| ▲ | CPU spike | peach |
| ✗ | error | red |
- The panel follows the newest event. Scrolling up pauses it. The subtitle then shows "PAUSED. End to follow."

### T. Title block
- Two lines of fields, each with a small label above its value.
- Fields: Title, Host, Kernel, Sample, Interval, Uptime, Sheet.
- Background `mantle`. Labels `subtext0`. Values `text`.
- The block is never hidden.

## 6. Typography

A terminal sets its own font. The TUI cannot choose it.
The TUI controls only color and text style.

### 6.1 Text roles
| Role | Style | Color |
|---|---|---|
| Title | bold | `text` |
| Subtitle, meta | normal | `overlay1` |
| Label | normal | `subtext0` |
| Data | normal | `text` |
| Emphasis | bold | `text` |
| Note | italic | `blue` |
| Clickable | underline | `sapphire` |
| Error | bold | `red` |
| Disabled | normal | `overlay0` |
| Selected | bold, `surface1` background | `text` |

- Use color tokens for dim text. Do not use the faint attribute. Many terminals ignore it.
- Italic may not render. A note must still read well without it.
- Never use blink.

### 6.2 Font for the user
- Recommended terminal fonts: SF Mono where you have it, else JetBrains Mono or Fira Code.
- Apple licenses SF Mono for Apple platform work, so sysmon never bundles it.

### 6.3 Font for HTML and SVG outputs
- Mono stack: `"SF Mono", SFMono-Regular, "JetBrains Mono", "Fira Code", ui-monospace, monospace`
- Title stack: `system-ui, -apple-system, "Segoe UI", sans-serif`
- Weights: 400 regular, 500 medium, 600 semibold, 700 bold. Italic is available.
- Sizes in px: 11 label, 12 data, 14 table title, 18 panel title.
- Data and code use mono. Panel titles use the title stack, as in the image.

## 7. Theme

The default theme is Catppuccin Mocha. The image is light, so the roles invert for a dark base.

### 7.1 Tokens
| Token | Hex | Token | Hex |
|---|---|---|---|
| rosewater | `#f5e0dc` | overlay2 | `#9399b2` |
| flamingo | `#f2cdcd` | overlay1 | `#7f849c` |
| pink | `#f5c2e7` | overlay0 | `#6c7086` |
| mauve | `#cba6f7` | surface2 | `#585b70` |
| red | `#f38ba8` | surface1 | `#45475a` |
| maroon | `#eba0ac` | surface0 | `#313244` |
| peach | `#fab387` | base | `#1e1e2e` |
| yellow | `#f9e2af` | mantle | `#181825` |
| green | `#a6e3a1` | crust | `#11111b` |
| teal | `#94e2d5` | text | `#cdd6f4` |
| sky | `#89dceb` | subtext1 | `#bac2de` |
| sapphire | `#74c7ec` | subtext0 | `#a6adc8` |
| blue | `#89b4fa` | lavender | `#b4befe` |

### 7.2 Roles
| Role | Token |
|---|---|
| Screen background | base |
| Title block background | mantle |
| Outer frame | overlay0 |
| Panel border | surface2 |
| Focused border and tab | lavender |
| Body text | text |
| Label | subtext0 |
| Meta, subtitle | overlay1 |
| Zebra band | surface0 |
| Selected row | surface1 |
| Hover row | surface0 |
| Annotation, note | blue |
| Clickable | sapphire |
| OK, running | green |
| Warning | yellow |
| Attention | peach |
| Error, zombie | red |
| Gauge track | surface0 |
| Dialog background | mantle |
| Kill dialog border | red |

### 7.3 Rules
- Define tokens once, in `internal/tui/theme`. No hex value appears anywhere else.
- Add no color package. Store the hex values in the theme package.
- Detect color depth. Use true color, then 256 colors, then 16 colors.
- If `NO_COLOR` is set, drop all color and keep the text styles and glyphs.
- If the locale is not UTF-8, use ASCII glyphs: `+` for ✓, `x` for ✗, `!` for !, `#` for █.
- Add a second theme later. Catppuccin Latte fits print and light diagrams.

## 8. Interaction

### 8.1 Focus
- Exactly one panel has focus. Keys go to it.
- The mouse wheel scrolls the panel under the pointer. It does not move focus.
- A click on a panel moves focus to it.

### 8.2 Mouse
| Action | Result |
|---|---|
| Click a panel | Focus it |
| Click a row | Select it |
| Double-click a row, or click `▸` | Expand or collapse |
| Click a column header | Sort. Click again to reverse. |
| Wheel | Scroll 3 rows per notch |
| Drag a scrollbar thumb | Scroll |
| Drag a divider between panels | Resize both panels |
| Click a tab in compact mode | Show that panel |
| Click a dialog button | Press it |
| Click outside a dialog | Cancel |

- Enable cell motion reporting. Dragging needs it.
- A press that moves 1 cell or more becomes a drag. A press that does not move is a click.
- Reordering panels by drag is not in version 0.1.

### 8.3 Keyboard
| Key | Action |
|---|---|
| Tab, Shift+Tab | Next, previous panel |
| Alt+A to Alt+F | Jump to a panel |
| Up, Down, `j`, `k` | Move selection |
| PgUp, PgDn | Move one page |
| Home, End | First, last row |
| Enter, Space | Expand or collapse |
| `s` | Cycle the sort column |
| `S` | Reverse the sort |
| `[`, `]` | Shrink, grow the focused panel by 1 step |
| `r` | Reset the layout |
| `c` | Toggle panel C |
| `x` | Kill the selected process |
| `?` | Help |
| Esc | Close a dialog |
| `q` | Quit |

### 8.4 Resize
- Handle every window size message. Recompute the layout in one pass.
- Store split sizes as ratios, not cells. Ratios survive a resize.
- Clamp to the minimum panel size.
- Keep the selected row in view after a resize.
- The `r` key restores the default ratios.

### 8.5 Drag state
`Idle` → `Pressed(target)` → `Dragging(target, origin)` → `Idle`
- A release in `Pressed` is a click.
- A release in `Dragging` ends the drag. Save the new ratio.
- Esc during a drag cancels it and restores the old ratio.

### 8.6 Kill dialog
- Open it with `x` or a button.
- Show the name, PID, and owner. Show "Send SIGTERM?".
- Enter confirms. Esc cancels.
- If the process still runs after 3 seconds, offer "Force kill (SIGKILL)".
- If the process is not yours, disable the action. Show the reason.

## 9. Rendering and safety

- `Layout(width, height, ratios) -> map[PanelID]Rect` is a pure function.
- `HitTest(x, y) -> Hit{Panel, Region, Index}` is a pure function. Regions: tab, body, scrollbar, divider, header, row.
- `Render(model, width, height) string` is a pure function.
- The UI reads data only through the `Source` interface.
- **Sanitize every process name and command line.** A process can put escape codes in its name. Strip all control characters before you draw.
- Truncate by display width. Names can hold wide characters.
- Render only visible rows.

## 10. Shared with diagrams

`tools/diagram.py` draws SVG with the same tokens and component names.
- Frame, column numbers, row letters, and lettered panels match this file.
- The title block holds: Title, Specification, Owner, Source, Sheet.
- Themes: `--theme mocha` (default) and `--theme latte`.
- A map file has 5 keys: `starts`, `reads`, `state`, `uses`, `fails`.
- The script draws one panel per key: A, B, C, D, E.
- It draws only what the map file states.

## 11. Tests

- **Golden files:** snapshot `Render` at 80 x 24, 120 x 36, and 200 x 50.
- **Hit tests:** a table of points and expected hits for each mode.
- **Resize:** random sizes from 1 x 1 to 400 x 100. No panic.
- **Theme:** a test fails if a hex value appears outside `internal/tui/theme`.
- **Sanitize:** a name with escape codes draws as plain text.
- **No color:** `NO_COLOR` output has no color codes.

## 12. Build slices

| Slice | Content | Starts after |
|---|---|---|
| 1 | Static sheet, Mocha theme, resize, fake data | Go step 4 |
| 2 | Scroll and keyboard selection | Slice 1 |
| 3 | Mouse click and wheel | Go step 8 |
| 4 | Drag dividers and scrollbars | Slice 3 |
| 5 | Expand a process and kill it | Slice 4 |

Each slice ends with its tests passing and the docs updated.

## 13. Not in version 0.1

- Drag to reorder panels
- Right-click menus
- Themes beyond Mocha and Latte
- A config file for layout
- Mouse support inside multiplexers