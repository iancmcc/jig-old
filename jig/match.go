package jig

import (
	"container/heap"
	"math"
	"sort"
	"strings"

	"github.com/xrash/smetrics"
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
	n := strings.Count(term, "/")
	if n > 0 {
		split := strings.SplitN(candidate, "/", 3-n)
		candidate = split[len(split)-1]
		// Just return the score comparing the two
		return smetrics.JaroWinkler(term, candidate, 0.7, 4)
	}
	i := 0.0
	segments := strings.Split(candidate, "/")
	sort.Reverse(sort.StringSlice(segments))
	for _, s := range segments {
		score := smetrics.JaroWinkler(term, s, 0.7, 4)
		i = math.Max(score, i)
		if i == 1 {
			break
		}
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
		})
	}
	return matches.ToStringArray()
}