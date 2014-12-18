package biocomp

import (
	"math/rand"
	"sort"
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
	maxSize++
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
	for i := ls - 1; i >= 0; i-- {
		mismatch = 0
		for j := 0; j+i < ls && j < lt; j++ {
			if s[i+j] != t[j] { // Suffix-prefix matching
				mismatch++
			}
		}
		if float64(mismatch)/float64(ls-i)-err < 1E-4 { // floating point error considered
			if ls-i >= lt { // s contains t
				return -1
			}
			score = ls - i
		}
	}
	return score
}

func (frags *Sequences) BuildOverlapEdges(err float64, reversed bool) Edges {
	ss := *frags
	es := make(Edges, 0, len(ss))
	ignore := make([]int, 0)
	for i := range ss {
		for j := range ss {
			if i == j {
				continue
			}
			if reversed {
				if score := ss[i].OverlapScore(ss[j], err); score >= 0 {
					es = append(es, &Edge{i, j, score})
				} else {
					ignore = append(ignore, j)
				}
				if score := ss[i].OverlapScore(ss[j].ReverseComplement(), err); score >= 0 {
					es = append(es, &Edge{i, -j, score})
				} else {
					ignore = append(ignore, j)
				}
				if score := ss[i].ReverseComplement().OverlapScore(ss[j], err); score >= 0 {
					es = append(es, &Edge{-i, j, score})
				} else {
					ignore = append(ignore, j)
				}
				if score := ss[i].ReverseComplement().OverlapScore(ss[j].ReverseComplement(), err); score >= 0 {
					es = append(es, &Edge{-i, -j, score})
				} else {
					ignore = append(ignore, j)
				}
			} else {
				if score := ss[i].OverlapScore(ss[j], err); score >= 0 {
					es = append(es, &Edge{i, j, ss[i].OverlapScore(ss[j], err)})
				} else {
					ignore = append(ignore, j)
				}
			}
		}
	}
	if len(ignore) > 0 {
		newFrags := make(Sequences, 0)
		for i := 0; i < len(ss); i++ {
			ign := false
			for _, ignored := range ignore {
				if i == ignored {
					ign = true
					break
				}
			}
			if !(ign) {
				newFrags = append(newFrags, ss[i])
			}
		}
		*frags = newFrags
		return frags.BuildOverlapEdges(err, reversed)
	}
	sort.Sort(es)
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

func EditDistance(s, t string) int {
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}
