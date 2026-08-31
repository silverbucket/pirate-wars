# pirate-wars
A pirate-themed game, written in Go, using [Ebitengine](https://ebitengine.org/).

**_NOTE_**: This is a hobby project, a work in progress.


![Pirate Wars! Title Screen](https://storage.5apps.com/silverbucket/public/shares/250421-2239-pirate-wars.png)
![Game Play Screenshot](https://storage.5apps.com/silverbucket/public/shares/250421-2238-Screenshot%202025-04-22%20at%2000.38.36.jpg)


## Overview

You are a pirate, sailing the seas. You sail by wind, trade between towns, and
examine or hail the ships you meet.

NPC traders sail their own routes between towns. Towns are generated across the
map and can be docked at to trade goods, drink in the tavern, and refit.

## Sailing

You steer by **tacking**, not by pointing at a compass direction. `A` and `D`
turn one octant to port or starboard **relative to your current heading**, and
the helm answers once per sailing tick — a ship cannot spin on the spot.

Your speed is `hull x sail x point of sail x wind`, a spread of roughly 45x
between the best and worst combinations. The side panel reports every term, and
the compass shows your heading against the wind: needles together is a run
(fastest), opposed is in irons (a standstill).

**Furled sail is a dead stop, not a slow crawl.** `S` at furled leaves you
becalmed; the panel says so, and `W` sets canvas again.

### Keybindings

| Key | Command |
|-----|---------|
| `A` / `D` | Tack one octant to port / starboard |
| `W` / `S` | Set more / less sail |
| `1` / `2` / `3` | Jump to full / half / furled sail |
| `Enter` / `O` | Dock at an adjacent town |
| `H` | Hail a ship you have come alongside |
| `X` | Examine ships and towns in sight |
| `M` | World minimap |
| `?` / `F1` | Help screen |
| `Esc` | Back out of any screen |
| `F3` | Debug overlay |
| `Ctrl+Q` | Quit (asks first) |

Every binding is also named on the action bar at the bottom of the window, and
every bar button is clickable.

## Features
* Sail by wind: heading, sail trim, point of sail and wind strength all feed speed
* Compass showing heading against wind, and a full instrument readout
* Explore an 800x800 generated world
* Dock at towns to trade goods, hear tavern rumours, and refit at the shipwright
* World minimap with towns and the on-screen area outlined
* NPC boats with basic pathfinding AI, running trade routes
* Examine and hail NPC ships

### Towns
* Towns don't spawn towns in small land-locked areas, however larger inaccessible areas can form with the terrain generation.

## Todo

#### Visuals
* ~~Use Tilemaps~~
* Rounded edges
* Animate tranistions
* Nice borders for panels
* Examine data popup over ship (rather than in side-panel)

#### Towns
* ~~Enter towns~~
* Make towns look better
* ~~Buy/sell goods~~
* Found your own town? (Pirate hideaway?)

#### Travel
* ~~Use wind and rotating ship to sail, speed etc.~~
* Engage with NPCs
* Improved NPC AI
* Hire/Dig channels pathways?
* Land defenses/fortifications
* ~~Don't allow overlap of ships (collision detection)~~
* ~~Wind direction determines ease of travel~~ (consume more food when going against wind)

### Ships 
* Fire from boat
* Upgrade
* Repair
* Buy/capture 
* Name your ship(s)
* Maintain a fleet
* Appoint Captains?

### Misc
* Combat: the game is called Pirate Wars and currently has neither
* Persistence between runs