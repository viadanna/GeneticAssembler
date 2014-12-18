package biocomp

import (
	"testing"
)

func TestAssembly(t *testing.T) {
	frags := &Sequences{"AAAAAGT", "AAA", "GTAAC", "AGTA", "AACC", "GTAACC"}
	edges := frags.BuildOverlapEdges(0, false)
	edgesMap := edges.BuildEdgesMap()
	path := Greedy(edges, edgesMap, len(*frags))
	assembly := path.GenomeAssembly(*frags)
	if assembly != "AAAAAGTAACC" {
		t.Error("Failed to assemble genome from fragments", path.Edges, *frags, assembly)
	}
}
