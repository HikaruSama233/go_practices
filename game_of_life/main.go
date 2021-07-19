package main
import (
	"fmt"
	"math/rand"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
	"os"
	"os/exec"
	"log"
	"runtime"
)

var (
	m = kingpin.Flag("row", "number of rows").Default("32").Int()
	n = kingpin.Flag("col", "number of cols").Default("64").Int()
	gens = kingpin.Flag("generations", "number of generations").Default("-1").Int()
	fps = kingpin.Flag("fps", "Frames per second").Default("10").Int()
	pcg = kingpin.Flag("percentage", "Percentage of living cells at the start").Short('p').Default("33").Int()
)

func main() {
	kingpin.Parse()
	fmt.Printf("m=%d, n=%d\n", *m, *n)
	new_world := NewWorld(*m, *n)
	new_world.RandomInit(*pcg)
	sleep_time := time.Duration(1000 / *fps) * time.Millisecond

	for g := 0; ; g++ {
		if g == *gens {
			break
		}
		cmd := exec.Command("clear")
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "cls")
		}
		cmd.Stdout = os.Stdout
		cmd.Run()
		new_world.ShowWorld()
		fmt.Printf("Generation: %d\n", g)
		new_world.Evolute()
		time.Sleep(sleep_time)
	}
}

type World struct {
	generation int
	m int
	n int
	board []bool
}

func NewWorld(m int, n int) *World {
	board := make([]bool, m*n)
	return &World{generation:0, m:m, n:n, board:board}
}

func (world *World) RandomInit(percentage int) {
	total_live_cells := len(world.board) * percentage / 100
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(world.board))
	for i, v := range perm {
		// fmt.Printf("%d, %d\n", i, v)
		if i < total_live_cells {
			world.board[v] = true
		} else {
			world.board[v] = false
		}
		 
	}
}

func (world *World) Evolute() {
	old_world := NewWorld(world.m, world.n)
	copy(old_world.board, world.board)
	var num_neighbors = 0
	for y := 0; y < world.m; y++ {
		for x := 0; x < world.n; x++ {
			if old_world.Get(x, y) {
				num_neighbors = old_world.CalculateNeighbors(x, y, true)
				if num_neighbors < 2 {
					world.Set(x, y, false)
				} else if num_neighbors > 3 {
					world.Set(x, y, false)
				}
			} else {
				num_neighbors = old_world.CalculateNeighbors(x, y, false)
				if num_neighbors == 3 {
					world.Set(x, y, true)
				}
			}
		}
	}
}

func (world *World) Get(x int, y int) bool {
	if x < 0 || y < 0 || x >= world.n || y >= world.m {
		return false
		// log.Fatal("Invalid Coordinate")
	}
	return world.board[y*world.n+x]
}

func (world *World) Set(x int, y int, status bool) {
	if x >= 0 && y >= 0 && x < world.n && y < world.m {
		world.board[y*world.n+x] = status
	} else {
		log.Fatal("Invalid Coordinate")
	}
}

func (world *World) CalculateNeighbors(x int, y int, live bool) int {
	var num = 0
	for i := y-1; i < y+2; i++ {
		for j := x-1; j < x+2; j++ {
			if world.Get(j, i) {
				num += 1
			}
		}
	}
	if live {
		num -= 1
	}
	return num
}

func (world *World) ShowWorld() {
	fmt.Printf("┌─")
	for i := 0; i < world.n; i++ {
		fmt.Printf("──")
	}
	fmt.Printf("─┐\n")

	for y := 0; y < world.m; y++ {
		fmt.Printf("│ ")
		for x := 0; x < world.n; x++ {
			if world.Get(x, y) {
				fmt.Printf("██")
			} else {
				fmt.Printf("[]")
			}
		}
		fmt.Printf(" │\n")
	}

	fmt.Printf("└─")
	for i := 0; i < world.n; i++ {
		fmt.Printf("──")
	}
	fmt.Printf("─┘\n")
}
