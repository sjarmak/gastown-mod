# Gas Town RTS Visualization — Product Requirements Document

## Overview

A browser-based real-time strategy game visualization of the Gas Town multi-agent orchestration system. Users see animated avatar sprites representing agents (polecats, witnesses, refineries, mayor, deacon, dogs) moving within rig territories on a 2D map, with real-time communication visualization and deep inspection on click.

## Core Principles

1. **It's a game view, not a dashboard** — agents are animated avatars that walk around, not cards in boxes
2. **Communication is visible** — when agents nudge, mail, or sling work, you SEE it happen (visual arcs, message particles, avatar movement between rigs)
3. **Click to inspect deeply** — clicking an agent shows its recent messages, nudges received, current task, bead assignments — the actual content, not just status labels
4. **All entities visible** — every rig, every polecat (even idle ones), every witness, refinery, dog, and the mayor/deacon must appear on the map

## User Experience

### The Map

- **2D top-down or slight isometric view** — NOT a card/flexbox layout
- **Rig territories** are distinct regions on the map, like territories in Risk or zones in an RTS
  - Each rig (codescalebench, gastown, oreilly, background_agents) is a labeled region with clear boundaries
  - Town Hall (mayor + deacon) is a central building/zone
- **Scale**: the map should be large enough that entities within a rig have room to move without overlapping
- **Pan** with click-drag, **zoom** with scroll wheel
- **Stable layout** — rig positions NEVER swap between updates. Use deterministic placement.

### Agent Avatars

- **Rendered as distinct sprite characters on a canvas** — not HTML divs
- Each role has a unique visual appearance (shape, color, size):
  - **Polecat** (🐱 orange, small worker sprite) — most numerous
  - **Witness** (🦉 green, medium overseer sprite) — one per rig
  - **Refinery** (🏭 pink/magenta, building-like sprite) — one per rig
  - **Mayor** (🎩 gold, larger central sprite) — one, at Town Hall
  - **Deacon** (🐺 cyan, medium sprite) — one, near Town Hall
  - **Dog** (🐕 gray, small sprite) — utility workers
  - **Crew** (👷 orange, like polecat) — named workers
- **Size**: sprites should be at minimum 32x32px at 1x zoom, with clear pixel art style
- **Name label** below each sprite: readable at default zoom (12px+ font)
- **Status indicator**: colored dot or glow (green=working, gray=idle, red=stuck, pulsing red=escalated)

### Movement & Animation

- **Within-rig wandering**: working agents slowly drift around within their rig territory. Idle agents stay mostly still with a gentle breathing/bob animation.
- **Cross-rig movement**: when an agent sends a nudge/mail to another rig, a visual messenger or particle arc travels from sender to recipient rig
- **Sling animation**: when the mayor slings work to a polecat, show the work item flying from Town Hall to the target rig
- **State transitions**: when an agent goes from idle→working, brief flash/sparkle. When stuck, pulsing red glow.

### Communication Visualization

- **Nudge**: a small glowing orb or arrow flies from sender to recipient across the map
- **Mail**: a letter/envelope icon flies between agents (slower than nudge, more prominent)
- **Sling**: a bright arc from Town Hall to the target rig with the bead ID briefly displayed
- **Speech bubbles**: when an agent is actively working, a small speech bubble floats above showing the truncated current activity (fade after 5 seconds, refresh on new activity)

### Click-to-Inspect (Detail Panel)

When you click an agent avatar, a side panel opens showing:

1. **Identity**: name, role, rig, account (if known)
2. **Status**: working/idle/stuck/dead, session alive/dead
3. **Current Work**:
   - Hooked bead ID + title
   - Current activity description
   - Time since last activity
4. **Recent Messages** (last 5):
   - Nudges received (from whom, content preview)
   - Mail received (subject, from, preview)
   - Nudges/mail sent (to whom, content preview)
5. **Session Info**:
   - tmux session name
   - Last 3-5 lines of terminal output (from tmux capture-pane)

The panel stays open until closed. Clicking another agent switches the panel to that agent.

### HUD (Heads-Up Display)

Top bar showing:
- **GAS TOWN** title
- Agent count breakdown: `12 agents (4 polecat · 4 witness · 4 refinery · 1 mayor · 1 deacon · 2 dog)`
- Active / Idle / Stuck counts
- Merge queue count
- Escalation alerts (pulsing red if any)
- Connection status (LIVE / RECONNECTING)

## Technical Architecture

### Frontend: Phaser 3 + HTML overlay

- **Phaser 3** for the game canvas: sprites, movement, animations, particle effects, camera/zoom
- **HTML overlay** for: HUD stats bar, detail panel (too text-heavy for canvas), legend
- **No card-based layout** — all agents are canvas sprites positioned in world coordinates
- Phaser `pixelArt: true` for crisp rendering
- Sprites generated procedurally (colored rectangles with eyes, or simple pixel art)

### Backend: Go SSE endpoint (existing)

Extend `/api/rts-state` SSE endpoint to include:

1. **All entities** — including idle polecats, refineries (currently missing from the data)
2. **Recent messages per agent** — last 5 nudges/mail received and sent
3. **Terminal output** — last 3 lines from tmux capture-pane for active sessions

New data structures needed:

```json
{
  "entities": [
    {
      "id": "co-slit",
      "type": "polecat",
      "name": "slit",
      "rig": "codescalebench",
      "status": "working",
      "activity": "Running OH gap fill via Daytona...",
      "activityAge": "2m",
      "hookBead": "co-wlr",
      "hookTitle": "Relaunch OH gap fill via Daytona",
      "sessionAlive": true,
      "account": "account1",
      "recentMessages": [
        {
          "direction": "received",
          "type": "nudge",
          "from": "mayor",
          "preview": "CRITICAL: You are working on co-wlr...",
          "timestamp": "2026-03-19T15:24:00Z"
        }
      ],
      "terminalOutput": "$ csb run --agent openhands --parallel 30\n[csb run] Launching..."
    }
  ],
  "rigs": [...],
  "events": [
    {
      "type": "nudge",
      "from": "mayor",
      "to": "codescalebench/polecats/slit",
      "fromRig": "__townhall",
      "toRig": "codescalebench",
      "timestamp": "2026-03-19T15:24:00Z"
    }
  ]
}
```

The `events` array contains recent communication events (last 30 seconds) for the frontend to animate as flying particles/arcs.

### Data Sources

All data sources already exist — no new infrastructure needed:

1. **ConvoyFetcher** (existing): FetchWorkers, FetchSessions, FetchRigs, FetchDogs, FetchMayor, FetchHooks, FetchMergeQueue, FetchEscalations
2. **Tmux capture-pane**: `tmux capture-pane -t <session> -p -S -3` for terminal output
3. **Mail inbox**: `gt mail inbox --json` or direct Dolt query for recent messages
4. **Nudge log**: check if nudge history is available in Dolt or activity log

### Missing from current API

The current `FetchWorkers()` doesn't return:
- Idle polecats (only active ones with sessions)
- Refineries as entities (they're returned as worker type but may be filtered)

The current `FetchSessions()` doesn't return:
- Account information (which Claude account is running)

These need to be fixed in `rts.go` `fetchGameState()`.

## File Structure

```
internal/web/
├── static/rts/
│   └── index.html          # Phaser 3 game + HTML overlays (single file)
├── rts.go                   # SSE handler + data enrichment
└── handler.go               # Route wiring (already has /rts)
```

Keep it as a single HTML file with inline JS for simplicity (no build system). Phaser loaded from CDN.

## What We Learned from v1-v3

### What worked
- HTML detail panel with structured rows (keep this as overlay, not canvas)
- HUD stats bar (keep as HTML overlay)
- SSE for real-time updates (keep, 1-second polling is fine)
- Distinct colors per role (keep: orange polecat, green witness, pink refinery, gold mayor, cyan deacon, gray dog)

### What didn't work
- **Isometric view** (v1): too cramped, depth sorting caused overlap, tiles wasted space
- **Card-based layout** (v2-v3): not game-like, no sense of movement or space
- **Hash-based entity positioning** (v1): caused collisions
- **Circular rig layout** (v2): rigs overlapped
- **Re-rendering full HTML on every SSE update** (v3): caused position swapping, broke animations
- **Phaser speech bubbles overlapping sprites** (v1): canvas text at 7px is unreadable
- **Tiny sprites** (12-24px in v1): unreadable

### Key requirements from user feedback
1. "I want to see them more like moving around" — canvas sprites with pathfinding/wandering
2. "when I click on it I should see the most recent message state" — nudge/mail content
3. "what it was tasked with (like if it received a nudge from you what that looked like)" — show actual nudge text
4. "which account it's associated with" — show account info
5. "the rectangles still overlap" — proper spacing, collision avoidance
6. "super low resolution and hard to read" — larger sprites, readable labels
7. "witnesses say 'working' but you can't see what they're working on" — show actual activity
8. "codescalebench and town hall keep swapping labels" — stable deterministic layout

## Acceptance Criteria

1. All 4 rigs + Town Hall visible on the map with NO overlap, NO position swapping
2. Every entity type rendered as a distinct animated sprite on canvas
3. Agents wander within their rig territory (working agents move more, idle agents mostly still)
4. Communication events (nudge, mail, sling) visualized as flying particles between rigs
5. Click any agent → side panel shows: status, current bead/task, last 5 messages with actual content, terminal output
6. HUD shows accurate counts for all agent types
7. Readable at default zoom — sprite labels ≥ 12px, names visible
8. Pan (drag) and zoom (scroll) work smoothly
9. Layout is deterministic — positions don't change between SSE updates
