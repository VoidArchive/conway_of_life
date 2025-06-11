package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

// ANSI color codes for colorful output
const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorBlue   = "\x1b[34m"
	colorPurple = "\x1b[35m"
	colorCyan   = "\x1b[36m"
	colorWhite  = "\x1b[37m"

	// Bright colors
	colorBrightRed    = "\x1b[91m"
	colorBrightGreen  = "\x1b[92m"
	colorBrightYellow = "\x1b[93m"
	colorBrightBlue   = "\x1b[94m"
	colorBrightPurple = "\x1b[95m"
	colorBrightCyan   = "\x1b[96m"
	colorBrightWhite  = "\x1b[97m"
)

// Cell represents a cell in the game with its age for color variation
type Cell struct {
	alive bool
	age   int
}

// Game represents the Conway's Game of Life state
type Game struct {
	width  int
	height int
	cur    [][]Cell
	nxt    [][]Cell
	out    *bufio.Writer
}

// NewGame creates a new Game instance with the given dimensions
func NewGame(width, height int) *Game {
	g := &Game{
		width:  width,
		height: height,
		out:    bufio.NewWriter(os.Stdout),
	}

	// Initialize double buffers
	g.cur = make([][]Cell, height)
	g.nxt = make([][]Cell, height)
	for y := range g.cur {
		g.cur[y] = make([]Cell, width)
		g.nxt[y] = make([]Cell, width)
		for x := range g.cur[y] {
			// 25% chance of being alive initially
			g.cur[y][x] = Cell{
				alive: rand.Intn(4) == 0,
				age:   0,
			}
		}
	}

	return g
}

// getColorForAge returns an ANSI color code based on cell age
func getColorForAge(age int) string {
	colors := []string{
		colorBrightWhite,  // newborn - bright white
		colorBrightCyan,   // young - bright cyan
		colorBrightGreen,  // mature - bright green
		colorBrightYellow, // aging - bright yellow
		colorYellow,       // old - yellow
		colorRed,          // very old - red
		colorPurple,       // ancient - purple
	}

	if age >= len(colors) {
		return colorPurple // Ancient cells stay purple
	}
	return colors[age]
}

// countNeighbors counts living neighbors around a cell
func (g *Game) countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			// Wrap around edges (toroidal topology)
			nx := (x + dx + g.width) % g.width
			ny := (y + dy + g.height) % g.height
			if g.cur[ny][nx].alive {
				count++
			}
		}
	}
	return count
}

// update applies Conway's Game of Life rules and updates cell ages
func (g *Game) update() {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			neighbors := g.countNeighbors(x, y)
			currentCell := g.cur[y][x]

			// Apply Conway's rules
			if currentCell.alive {
				// Living cell survives with 2 or 3 neighbors
				g.nxt[y][x].alive = neighbors == 2 || neighbors == 3
				if g.nxt[y][x].alive {
					// Cell survives, increment age
					g.nxt[y][x].age = currentCell.age + 1
				} else {
					// Cell dies, reset age
					g.nxt[y][x].age = 0
				}
			} else {
				// Dead cell becomes alive with exactly 3 neighbors
				g.nxt[y][x].alive = neighbors == 3
				if g.nxt[y][x].alive {
					// New cell is born, age is 0
					g.nxt[y][x].age = 0
				} else {
					// Cell remains dead
					g.nxt[y][x].age = 0
				}
			}
		}
	}

	// Swap buffers
	g.cur, g.nxt = g.nxt, g.cur
}

// draw renders the current game state with colors
func (g *Game) draw() {
	fmt.Fprint(g.out, "\x1b[H") // Move cursor to home position

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			cell := g.cur[y][x]
			if cell.alive {
				color := getColorForAge(cell.age)
				fmt.Fprintf(g.out, "%s█%s", color, colorReset)
			} else {
				g.out.WriteRune(' ')
			}
		}
		g.out.WriteByte('\n')
	}
	g.out.Flush()
}

// cleanup restores terminal state
func (g *Game) cleanup() {
	fmt.Fprint(g.out, colorReset)
	fmt.Fprint(g.out, "\x1b[?25h") // Show cursor
	g.out.Flush()
}

// run starts the game loop
func (g *Game) run() {
	// Setup signal handler for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Hide cursor and clear screen
	fmt.Fprint(g.out, "\x1b[?25l") // Hide cursor
	fmt.Fprint(g.out, "\x1b[2J")   // Clear screen

	defer g.cleanup()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutting down gracefully...")
			return
		case <-ticker.C:
			g.draw()
			g.update()
		}
	}
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Get terminal size
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		cols, rows = 80, 24
	}

	// Adjust for proper display (leave space for newlines)
	height := rows - 1

	// Create and run the game
	game := NewGame(cols, height)
	game.run()
}
