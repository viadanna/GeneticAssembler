package biocomp

import (
	"errors"
	"fmt"
	"math/rand"
)

type Path struct {
	Edges
	Score int
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p *Path) Contains(vertice int) bool {
	last := len(p.Edges) - 1
	for i, edge := range p.Edges {
		if Abs(edge.B) == Abs(vertice) || (i != last && Abs(edge.A) == Abs(vertice)) {
			return true
		}
	}
	return false
}

func (p *Path) AddEdge(e *Edge) {
	if len(p.Edges) > 0 && p.Edges[len(p.Edges)-1].B != e.A {
		panic("Trying to add invalid path")
	}
	p.Edges = append(p.Edges, e)
	p.Score += e.S
}

func (p Path) String() string {
	s := fmt.Sprintf("%d", p.Edges[0].A)
	for _, edge := range p.Edges {
		s = fmt.Sprintf("%s %d", s, edge.B)
	}
	s = fmt.Sprintf("%s (%d)", s, p.Score)
	return s
}

type Paths []Path

func (p Paths) Len() int           { return len(p) }
func (p Paths) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Paths) Less(i, j int) bool { return p[i].Score > p[j].Score } // Reverse sorting

type Edge struct {
	A int
	B int
	S int
}

type Edges []*Edge

func (e Edges) Len() int           { return len(e) }
func (e Edges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e Edges) Less(i, j int) bool { return e[i].S > e[j].S } // Reverse sorting

func (e Edges) String() string {
	s := ""
	for _, v := range e {
		s = fmt.Sprintf("%s(%v)", s, *v)
	}
	return s
}

func (edges Edges) Pick(cutoff int, ignore Edges) (*Edge, error) {
	possible := make([]*Edge, 0)
	maxScore := -1
	for _, edge := range edges {
		already := false
		for _, ig := range ignore {
			if edge.B == ig.A {
				already = true
				break
			}
		}
		if !(already) {
			possible = append(possible, edge)
			if maxScore < edge.S {
				maxScore = edge.S
			}
		}
	}
	minScore := maxScore * cutoff / 100
	enoughScore := make([]*Edge, 0)
	for _, edge := range possible {
		if edge.S >= minScore {
			enoughScore = append(enoughScore, edge)
		}
	}
	if len(enoughScore) == 0 {
		return nil, errors.New("No possible edges found")
	}
	return enoughScore[rand.Intn(len(enoughScore))], nil
}
