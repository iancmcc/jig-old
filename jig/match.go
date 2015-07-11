package jig

import (
	"container/heap"
	"math"
	"strings"

	"github.com/antzucaro/matchr"
)

type Match struct {
	score float64
	match string
}

type MatchHeap []*Match

func (h MatchHeap) Len() int { return len(h) }
func (h MatchHeap) Less(i, j int) bool {
	return h[i].score > h[j].score
}

func (h MatchHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MatchHeap) Push(x interface{}) {
	match := x.(*Match)
	*h = append(*h, match)
}

func (h *MatchHeap) Pop() interface{} {
	old := *h
	n := len(old)
	match := old[n-1]
	*h = old[0 : n-1]
	return match
}

func (h MatchHeap) ToStringArray() []string {
	result := []string{}
	for h.Len() > 0 {
		match := heap.Pop(&h).(*Match)
		result = append(result, match.match)
	}
	return result
}

func BestScore(term, candidate string) float64 {
	if strings.Contains(term, "/") {
		// Just return the score comparing the two
		return matchr.Jaro(term, candidate)
	}
	i := 0.0
	segments := strings.Split(candidate, "/")
	for x, s := range segments {
		i = math.Max(matchr.Jaro(term, s)+float64(x), i)
	}
	return i
}

func SortedMatches(term string, candidates []string) []string {
	matches := &MatchHeap{}
	heap.Init(matches)
	for _, s := range candidates {
		heap.Push(matches, &Match{
			match: s,
			score: BestScore(term, s),
			//score: matchr.DamerauLevenshtein(term, rejoined),
		})
	}
	return matches.ToStringArray()
}