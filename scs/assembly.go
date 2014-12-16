package biocomp

import (
	"fmt"
	"math/rand"
	"sort"
)

type Parent map[int]*Edge
type Parents [2]Parent

func (path Path) GenomeAssembly(fragments Sequences) (genome string) {
	first := path.Edges[0].A
	if first < 0 {
		genome = string(fragments[Abs(first)].ReverseComplement())
	} else {
		genome = string(fragments[first])
	}
	for _, edge := range path.Edges {
		if edge.B < 0 {
			genome = fmt.Sprintf("%s%s", genome, fragments[Abs(edge.B)][edge.S:].ReverseComplement())
		} else {
			genome = fmt.Sprintf("%s%s", genome, fragments[edge.B][edge.S:])
		}
	}
	return genome
}

func GeneticAssembly(edges Edges, edgesMap map[int]Edges, vertices int, parents Parents, cutoff int) (path Path) {
	choice, err := edges.Pick(cutoff, Edges{})
	if err != nil {
		panic("Couldn't even start")
	}
	path.AddEdge(choice)
	for len(path.Edges) < vertices {
		lastVertice := path.Edges[len(path.Edges)-1].B
		genes := make(Edges, 0, 2)
		// Find possible edges from parents
		for _, parent := range parents {
			if parent != nil {
				if nextEdge, ok := parent[lastVertice]; ok && !(path.Contains(nextEdge.B)) {
					genes = append(genes, nextEdge)
				}
			}
		}
		if len(genes) > 0 { // Pick random edge from parents
			path.AddEdge(genes[rand.Intn(len(genes))])
		} else { // Add random possible edge based on cutoff rate
			choice, err = edgesMap[lastVertice].Pick(cutoff, path.Edges)
			if err != nil {
				break
			}
			path.AddEdge(choice)
		}
	}
	return
}

func MakeParents(possible Paths) (p Parents) {
	for i := 0; i < 2; i++ {
		p[i] = make(Parent)
		for _, e := range possible[i].Edges {
			p[i][e.A] = e
		}
	}
	return
}

func Greedy(edges Edges, edgesMap map[int]Edges, vertices int) Path {
	// Consider edges already sorted by reverse score
	p := Path{Edges{edges[0]}, edges[0].S} // Start with highest weight
	// Run until all veritices are added
	for len(p.Edges) < vertices {
		last := p.Edges[len(p.Edges)-1]
		// Find best possible edge and add to graph
		for _, next := range edgesMap[last.B] {
			if !(p.Contains(next.B)) {
				p.AddEdge(next)
				break
			}
		}
	}
	return p
}

func Genetic(edges Edges, edgesMap map[int]Edges, vertices int, cutoff, generations, children int) (paths Paths) {
	paths = make(Paths, 0, generations*children)
	var parents Parents
	// Run for given generations
	for g := 0; g < generations; g++ {
		children := make(Paths, children)
		// Generate a number of children per generation
		for c := 0; c < len(children); c++ {
			children[c] = GeneticAssembly(edges, edgesMap, vertices, parents, cutoff)
		}
		// Get best children as parents for next generation
		sort.Sort(children)
		parents = MakeParents(children)
		paths = append(paths, children...)
	}
	// Sort results by reverse score
	sort.Sort(paths)
	return
}
