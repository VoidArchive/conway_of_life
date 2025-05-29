package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	alive = '█'
	dead  = ' '
)

func main() {
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		cols, rows = 80, 24
	}
	w, h := cols, rows

	// double buffers
	cur := make([][]bool, h)
	nxt := make([][]bool, h)
	for y := range cur {
		cur[y] = make([]bool, w)
		nxt[y] = make([]bool, w)
		for x := range cur[y] {
			cur[y][x] = rand.Intn(4) == 0 // 25 % alive
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprint(out, "\x1b[2J") // clear screen once

	for {
		draw(out, cur)
		for y := range cur {
			for x := range cur[y] {
				n := neigh(cur, x, y, w, h)
				if cur[y][x] {
					nxt[y][x] = n == 2 || n == 3
				} else {
					nxt[y][x] = n == 3
				}
			}
		}
		cur, nxt = nxt, cur
		time.Sleep(80 * time.Millisecond)
	}
}

func neigh(g [][]bool, x, y, w, h int) int {
	n := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			xx := (x + dx + w) % w
			yy := (y + dy + h) % h
			if g[yy][xx] {
				n++
			}
		}
	}
	return n
}

func draw(out *bufio.Writer, g [][]bool) {
	fmt.Fprint(out, "\x1b[H") // cursor home
	for _, row := range g {
		for _, cell := range row {
			if cell {
				out.WriteRune(alive)
			} else {
				out.WriteRune(dead)
			}
		}
		out.WriteByte('\n')
	}
	out.Flush()
}
