package biocomp

import (
	"log"
	"math/rand"
)

type Sequence string
type Sequences []Sequence

func GenerateSequence(length int, alphabet string) Sequence {
	runes := []rune(alphabet)
	alphaLength := len(runes)
	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = runes[rand.Intn(alphaLength)]
	}
	return Sequence(string(result))
}

func GenerateDNA(length int) Sequence {
	return GenerateSequence(length, "ATCG")
}

func (s Sequence) ReverseComplement() Sequence {
	pairs := map[rune]rune{
		'A': 'T',
		'T': 'A',
		'C': 'G',
		'G': 'C',
	}
	r := []rune(string(s))
	for i := 0; i < len(s)/2; i++ {
		r[i], r[len(s)-i-1] = pairs[r[len(s)-i-1]], pairs[r[i]]
	}
	return Sequence(r)
}

func (s Sequence) Fragmentize(minSize, maxSize int, coverage float64, err float64, reverse float64) *Sequences {
	bases, target := 0, int(coverage*float64(len(s)))
	fragments := make(Sequences, 0, 2*target/(minSize+maxSize))
	for bases < target {
		// Pick random fragment
		size := rand.Intn(maxSize-minSize) + minSize
		start := rand.Intn(len(s) - size)
		fragment := s[start : start+size]
		// Reverse fragment
		if reverse > 0 && rand.Float64() < reverse {
			fragment = fragment.ReverseComplement()
		}
		// Add errors
		if err > 0 {
			rs := []rune(string(fragment))
			bases := []rune("ATCG")
			for i := range rs {
				if rand.Float64() < err {
					rs[i] = bases[rand.Intn(4)]
				}
			}
			fragment = Sequence(rs)
		}
		fragments = append(fragments, fragment)
		bases += size
	}
	return &fragments
}

func (s Sequence) OverlapScore(t Sequence, err float64) int {
	ls, lt := len(s), len(t)
	score, mismatch := 0, 0
	for i := 1; i < ls && i < lt; i++ {
		mismatch = 0
		for j := 0; j < i; j++ {
			if s[ls-j-1] != t[j] {
				mismatch++
			}
		}
		if float64(mismatch)/float64(i) <= err {
			score = i
		}
	}
	return score
}

func (ss Sequences) BuildOverlapEdges(err float64, reversed bool) Edges {
	es := make(Edges, 0, len(ss))
	if reversed {
		for i := range ss {
			for j := range ss {
				if i == j {
					continue
				}
				es = append(es, &Edge{i, j, ss[i].OverlapScore(ss[j], err)})
				es = append(es, &Edge{i, -j, ss[i].OverlapScore(ss[j].ReverseComplement(), err)})
				es = append(es, &Edge{-i, j, ss[i].ReverseComplement().OverlapScore(ss[j], err)})
				es = append(es, &Edge{-i, -j, ss[i].ReverseComplement().OverlapScore(ss[j].ReverseComplement(), err)})
			}
		}
	} else {
		for i := range ss {
			for j := range ss {
				if i == j {
					continue
				}
				edge := Edge{i, j, ss[i].OverlapScore(ss[j], err)}
				es = append(es, &edge)
				if edge.S > 5 {
					log.Println(edge)
				}
			}
		}
	}
	return es
}

func (edges Edges) BuildEdgesMap() (edgesMap map[int]Edges) {
	edgesMap = make(map[int]Edges)
	for _, edge := range edges {
		if edgesMap[edge.A] == nil {
			edgesMap[edge.A] = make(Edges, 0)
		}
		edgesMap[edge.A] = append(edgesMap[edge.A], edge)
	}
	return
}
