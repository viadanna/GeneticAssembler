package main

import (
	. "biocomp/scs"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"time"
)

var (
	LEN      = flag.Int("genome", 1000, "Size of genome to generate")
	MIN_SIZE = flag.Int("min", 25, "Minimum size of fragments")
	MAX_SIZE = flag.Int("max", 100, "Maximum size of fragments")
	COVERAGE = flag.Float64("coverage", 11, "Genome coverage of fragments")
	SEQ_ERR  = flag.Float64("seqerror", 0.1, "Sequencing error rate [0,1]")
	SEQ_REV  = flag.Float64("seqrev", 0.5, "Sequencing reverses rate [0,1]")
	ERR      = flag.Float64("error", 0.1, "Error considered in graph building [0,1]")
	GEN      = flag.Int("gen", 10, "Number of generations to try")
	CHILD    = flag.Int("child", 10, "Number of children per generation")
	CUTOFF   = flag.Int("cutoff", 90, "Percentage of max score to accept edges in random selector")
	CPUPROF  = flag.String("cpuprofile", "", "Write CPU profile to file")
	MP       = flag.Bool("mp", false, "Run Parallel")
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

	genome := GenerateDNA(*LEN)
	log.Println("Generated Genome")
	log.Println(genome)
	fragments := genome.Fragmentize(*MIN_SIZE, *MAX_SIZE, *COVERAGE, *SEQ_ERR, *SEQ_REV)
	log.Println("Generated Fragments")
	edges := fragments.BuildOverlapEdges(*ERR, *SEQ_REV > 0)
	log.Println("Generated Edges from fragments")
	sort.Sort(edges)
	log.Println("Sorted edges")
	edgesMap := edges.BuildEdgesMap()
	log.Println("Built edges map")

	// Run naive greedy algorithm
	naiveEdges := fragments.BuildOverlapEdges(0, false)
	naiveMap := naiveEdges.BuildEdgesMap()
	naivePath := Greedy(naiveEdges, naiveMap, len(*fragments))
	fmt.Println("Naive Greedy Score:", naivePath.Score)
	fmt.Println("Naive Assembly:", naivePath.GenomeAssembly(*fragments))
	// Run greedy algorithm with error and reversals
	greedyPath := Greedy(edges, edgesMap, len(*fragments))
	fmt.Println("Greedy Score:", greedyPath.Score)
	fmt.Println("Greedy Assembly:", greedyPath.GenomeAssembly(*fragments))
	// Run genetic algorithm
	geneticPath := Genetic(edges, edgesMap, len(*fragments), *CUTOFF, *GEN, *CHILD)[0]
	fmt.Println("Genetic Score:", geneticPath.Score)
	fmt.Println("Genetic Assembly:", geneticPath.GenomeAssembly(*fragments))
	// fmt.Println(paths)
}
