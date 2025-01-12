# pirate-wars
A pirate-themed roguelike, written in Go.

![pirate-wars](https://storage.5apps.com/silverbucket/public/shares/250110-1732-Screenshot%202025-01-10%20at%2018.31.52.jpg)

**_NOTE_**: This is a hobby project, a work in progress.

## Overview

You are a pirate, sailing the seas. 

Currently there are NPC ships (`⏏`) which have basic pathfinding capabilities. They travel from one town to another (a "trade route"). 

Towns are also generated throughout the map, with the red (`⩎`) characters. 

You currently cannot interaction with the towns or ships.

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

### Developer Commands
* `p`: View heat-map for town 0 (debug purposes)

## Features
* Move around in your boat
* Explore the map
* Visit towns (currently you cannot enter them)
* View mini-map of entire world, with towns listed (`m`)
* NPC boats with basic pathfinding AI

### Towns
* Towns don't spawn towns in small land-locked areas, however larger inaccessible can form with the terrain generation.

## Todo

#### Towns
* Enter towns
* Make towns look better
* Buy/sell goods
* Found your own town? (Pirate hideaway?)

#### World Map
* Engage with NPCs
* Improved NPC AI
* Hire/Dig channels pathways?
* Land defenses/fortifications
* Dont allow overlap of ships (collision detection)

### Ships 
* View ship details
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