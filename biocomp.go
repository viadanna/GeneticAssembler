package main

import (
	. "biocomp/scs"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"time"
)

var (
	INPUT   = flag.String("input", "-", "Edges file to be read or - for stdin(default)")
	GEN     = flag.Int("gen", 10, "Number of generations to try")
	CHILD   = flag.Int("child", 10, "Number of children per generation")
	CUTOFF  = flag.Int("cutoff", 90, "Percentage of max score to accept edges in random selector")
	CPUPROF = flag.String("cpuprofile", "", "Write CPU profile to file")
	MP      = flag.Bool("mp", false, "Run Parallel")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	// Profiling
	if *CPUPROF != "" {
		f, err := os.Create(*CPUPROF)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Setup Multiprocessing, not implemented yet
	if *MP {
		log.Println("Using", runtime.NumCPU(), "Processes")
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Open file
	var input *os.File
	var err error
	if *INPUT == "-" {
		input = os.Stdin
	} else {
		input, err = os.Open(*INPUT)
		if err != nil {
			log.Fatal("Failed to open file", *INPUT)
		}
	}

	// Read edges
	edges := make(Edges, 0)
	var a, b, s int
	for {
		_, err = fmt.Fscanf(input, "%d %d %d\n", &a, &b, &s)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Failed to read file")
		}
		edges = append(edges, &Edge{a, b, s})
	}

	// Build nested dictionary of edges
	edgesMap := make(map[int]Edges)
	for _, edge := range edges {
		if edgesMap[edge.A] == nil {
			edgesMap[edge.A] = make(Edges, 0)
		}
		edgesMap[edge.A] = append(edgesMap[edge.A], edge)
	}

	// Sort edges
	sort.Sort(edges)

	// Run genetic algorithm
	Genetic(edges, edgesMap, *CUTOFF, *GEN, *CHILD)
	// fmt.Println(paths)
}
