package main

import (
	. "biocomp/scs"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	RUNS     = flag.Int("runs", 10, "Number of experiment runs")
	LEN      = flag.Int("genome", 1000, "Size of genome to generate")
	MIN_SIZE = flag.Int("min", 25, "Minimum size of fragments")
	MAX_SIZE = flag.Int("max", 100, "Maximum size of fragments")
	COVERAGE = flag.Float64("coverage", 1, "Genome coverage of fragments")
	SEQ_ERR  = flag.Float64("seqerror", 0.01, "Sequencing error rate [0,1]")
	SEQ_REV  = flag.Float64("seqrev", 0.5, "Sequencing reverses rate [0,1]")
	ERR      = flag.Float64("error", 0.1, "Error considered in graph building [0,1]")
	GEN      = flag.Int("gen", 10, "Number of generations to try")
	CHILD    = flag.Int("child", 10, "Number of children per generation")
	CUTOFF   = flag.Int("cutoff", 90, "Percentage of max score to accept edges in random selector")
	CPUPROF  = flag.String("cpuprofile", "", "Write CPU profile to file")
	MP       = flag.Bool("mp", false, "Run Parallel")
)

type Experiment struct {
	RefGenome          Sequence
	RefLength          int
	GreedyAssembly     string
	GreedyLength       int
	GreedyScore        int
	GreedyDistance     int
	GreedyPlusAssembly string
	GreedyPlusLength   int
	GreedyPlusScore    int
	GreedyPlusDistance int
	GeneticAssembly    string
	GeneticLength      int
	GeneticScore       int
	GeneticDistance    int
	RandomAssembly     string
	RandomLength       int
	RandomScore        int
	RandomDistance     int
}

type Results struct {
	Experiments            []Experiment
	GreedyMeanLength       int
	GreedyMeanScore        int
	GreedyMeanDistance     int
	GreedyPlusMeanLength   int
	GreedyPlusMeanScore    int
	GreedyPlusMeanDistance int
	GeneticMeanLength      int
	GeneticMeanScore       int
	GeneticMeanDistance    int
	RandomMeanLength       int
	RandomMeanScore        int
	RandomMeanDistance     int
}

func RunExperiment(length, min, max, gen, child, cutoff int, coverage, seqErr, seqRev, errorRate float64) Experiment {
	e := Experiment{}
	e.RefGenome = GenerateDNA(length)
	e.RefLength = len(e.RefGenome)
	fragments := e.RefGenome.Fragmentize(min, max, coverage, seqErr, seqRev)
	edges := fragments.BuildOverlapEdges(*ERR, *SEQ_REV > 0)
	edgesMap := edges.BuildEdgesMap()

	// Run naive greedy algorithm
	if seqErr > 0 || seqRev > 0 {
		naiveEdges := fragments.BuildOverlapEdges(0, false)
		naivePath := Greedy(naiveEdges, naiveEdges.BuildEdgesMap(), len(*fragments))
		e.GreedyScore = naivePath.Score
		e.GreedyAssembly = naivePath.GenomeAssembly(*fragments)
		e.GreedyLength = len(e.GreedyAssembly)
		e.GreedyDistance = EditDistance(string(e.RefGenome), e.GreedyAssembly)
	}

	// Run greedy algorithm with error and reversals
	greedyPath := Greedy(edges, edgesMap, len(*fragments))
	e.GreedyPlusScore = greedyPath.Score
	e.GreedyPlusAssembly = greedyPath.GenomeAssembly(*fragments)
	e.GreedyPlusLength = len(e.GreedyPlusAssembly)
	e.GreedyPlusDistance = EditDistance(string(e.RefGenome), e.GreedyPlusAssembly)

	// Run genetic algorithm
	geneticPath := Genetic(edges, edgesMap, len(*fragments), cutoff, gen, child)[0]
	e.GeneticScore = geneticPath.Score
	e.GeneticAssembly = geneticPath.GenomeAssembly(*fragments)
	e.GeneticLength = len(e.GeneticAssembly)
	e.GeneticDistance = EditDistance(string(e.RefGenome), e.GeneticAssembly)

	// Run genetic algorithm without parents
	randomPath := Genetic(edges, edgesMap, len(*fragments), cutoff, 1, gen*child)[0]
	e.RandomScore = randomPath.Score
	e.RandomAssembly = randomPath.GenomeAssembly(*fragments)
	e.RandomLength = len(e.RandomAssembly)
	e.RandomDistance = EditDistance(string(e.RefGenome), e.RandomAssembly)

	return e
}

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

	results := Results{}
	results.Experiments = make([]Experiment, *RUNS)
	// Setup Multiprocessing, not implemented yet
	if *MP {
		log.Println("Using", runtime.NumCPU(), "Processes")
		runtime.GOMAXPROCS(runtime.NumCPU())
		ch := make(chan Experiment)
		for i := 0; i < *RUNS; i++ {
			go func(res chan Experiment) {
				res <- RunExperiment(*LEN, *MIN_SIZE, *MAX_SIZE, *GEN, *CHILD, *CUTOFF, *COVERAGE, *SEQ_ERR, *SEQ_REV, *ERR)
			}(ch)
		}
		for i := 0; i < *RUNS; i++ {
			results.Experiments[i] = <-ch
		}
	} else {
		for i := 0; i < *RUNS; i++ {
			results.Experiments[i] = RunExperiment(*LEN, *MIN_SIZE, *MAX_SIZE, *GEN, *CHILD, *CUTOFF, *COVERAGE, *SEQ_ERR, *SEQ_REV, *ERR)
		}
	}
	for i := 0; i < *RUNS; i++ {
		results.GeneticMeanLength += results.Experiments[i].GeneticLength
		results.GreedyMeanLength += results.Experiments[i].GreedyLength
		results.GreedyPlusMeanLength += results.Experiments[i].GreedyPlusLength
		results.RandomMeanLength += results.Experiments[i].RandomLength
		results.GeneticMeanScore += results.Experiments[i].GeneticScore
		results.GreedyMeanScore += results.Experiments[i].GreedyScore
		results.GreedyPlusMeanScore += results.Experiments[i].GreedyPlusScore
		results.RandomMeanScore += results.Experiments[i].RandomScore
		results.GeneticMeanDistance += results.Experiments[i].GeneticDistance
		results.GreedyMeanDistance += results.Experiments[i].GreedyDistance
		results.GreedyPlusMeanDistance += results.Experiments[i].GreedyPlusDistance
		results.RandomMeanDistance += results.Experiments[i].RandomDistance
	}
	results.GeneticMeanLength /= *RUNS
	results.GreedyMeanLength /= *RUNS
	results.GreedyPlusMeanLength /= *RUNS
	results.RandomMeanLength /= *RUNS
	results.GeneticMeanScore /= *RUNS
	results.GreedyMeanScore /= *RUNS
	results.GreedyPlusMeanScore /= *RUNS
	results.RandomMeanScore /= *RUNS
	results.GeneticMeanDistance /= *RUNS
	results.GreedyMeanDistance /= *RUNS
	results.GreedyPlusMeanDistance /= *RUNS
	results.RandomMeanDistance /= *RUNS
	out, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
	}
	os.Stdout.Write(out)
}
