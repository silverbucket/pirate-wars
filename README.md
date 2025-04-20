# pirate-wars
A pirate-themed game, written in Go, using the [Fyne UI Toolkit](https://github.com/fyne-io/fyne).

**_NOTE_**: This is a hobby project, a work in progress.


![Pirate Wars! Title Screen](https://storage.5apps.com/silverbucket/public/shares/250420-0050-pwicon.png)
![Game Play Screenshot](https://storage.5apps.com/silverbucket/public/shares/250420-0046-pwss5.png)


## Overview

You are a pirate, sailing the seas. You can sail around, explore the map, and examine other ships you encounter.

Currently there are NPC ships (`⏏`) which have basic pathfinding capabilities. They travel from one town to another (a "trade route"). 

Towns are also generated throughout the map, with the red (`⩎`) characters. 

You currently cannot interaction with the towns or ships (other than examining).

## Keybindings

### Navigation
```
 q w e        y k u
 a   d  -or-  l   h
 z x c        b j n
```
*(or arrow keys)*

### Commands
* `ctrl-q`: Quit
* `m`: Mini-map
* `x`: Examine something on the map
* ~~`i`: View your info~~
* ~~`?`: Help screen~~

## Features
* Move around in your boat
* Explore the map
* Visit towns (currently you cannot enter them)
* View mini-map of entire world, with towns listed (`m`)
* NPC boats with basic pathfinding AI
* View NPC ship details

### Towns
* Towns don't spawn towns in small land-locked areas, however larger inaccessible can form with the terrain generation.

## Todo

#### Visuals
* Use Tilemaps
* Animate tranistions
* Nice borders for panels
* Examine data popup over ship (rather than in side-panel)

#### Towns
* Enter towns
* Make towns look better
* Buy/sell goods
* Found your own town? (Pirate hideaway?)

#### Travel
* Engage with NPCs
* Improved NPC AI
* Hire/Dig channels pathways?
* Land defenses/fortifications
* Don't allow overlap of ships (collision detection)
* Wind direction determines ease of travel (consume more food when going against wind)

### Ships 
* Fire from boat
* Upgrade
* Repair
* Buy/capture 
* Name your ship(s)
* Maintain a fleet
* Appoint Captains?

### Misc
* Lipgloss adaptive colors, for highlighting entities
* Bubbles loading spinner
* Bubbles help hints on bottom of screen