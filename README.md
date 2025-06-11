# Conway's Game of Life

A colorful simulation of Conway's Game of Life in the terminal, written in modern Go.

## Features

- **Colorful cells**: Cells change color based on their age
  - 🤍 Newborn cells are bright white
  - 🩵 Young cells are bright cyan  
  - 💚 Mature cells are bright green
  - 💛 Aging cells are bright yellow
  - 🟡 Old cells are yellow
  - 🔴 Very old cells are red
  - 🟣 Ancient cells are purple

- **Modern Go implementation**: 
  - Proper struct-based design
  - Method receivers following Go idioms
  - Graceful shutdown with signal handling
  - Efficient double buffering

- **Enhanced terminal experience**:
  - Auto-detects terminal size
  - Smooth animation with cursor hiding
  - Clean exit with `Ctrl+C`

## Usage

```bash
./life
```

Press `Ctrl+C` to exit gracefully.

## Building

```bash
go build -o life life.go
```

## Requirements

- Go 1.24.3 or later
- Terminal with ANSI color support

```
./life
```
