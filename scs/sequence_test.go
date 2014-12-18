package biocomp

import (
	"testing"
)

func TestOverlapScore1(t *testing.T) {
	s1 := Sequence("ATCGT")
	t1 := Sequence("GTCA")
	if s1.OverlapScore(t1, 0) != 2 {
		t.Error("OverlapScore expected 2, got instead", s1.OverlapScore(t1, 0))
	}
}
func TestOverlapScore2(t *testing.T) {
	s1 := Sequence("AAAAAGT")
	t1 := Sequence("AAAA")
	if s1.OverlapScore(t1, 0) != -1 {
		t.Error("OverlapScore expected -1, got instead", s1.OverlapScore(t1, 0))
	}
}
func TestOverlapScore3(t *testing.T) {
	s1 := Sequence("ATCG")
	t1 := Sequence("GCTA")
	if s1.OverlapScore(t1, 0.25) != 1 {
		t.Error("OverlapScore expected 1, got instead", s1.OverlapScore(t1, 0.25))
	}
}
func TestOverlapGraph(t *testing.T) {
	frags := &Sequences{"AAAAAGT", "AAA", "GTAACC"}
	if edges := frags.BuildOverlapEdges(0, false); len(edges) != 2 {
		t.Error("Failed to build overlap Graph")
	}
}
func TestOverlapEqualSequences(t *testing.T) {
	s1 := Sequence("ATCG")
	t1 := Sequence("ATCG")
	if s1.OverlapScore(t1, 0) != -1 {
		t.Error("OverlapScore expected -1, got instead", s1.OverlapScore(t1, 0))
	}
}
