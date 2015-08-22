package jig

import (
	"container/heap"
	"strings"

	log "github.com/Sirupsen/logrus"
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

func reverse(s []string) chan string {
	c := make(chan string, 1)
	go func() {
		defer close(c)
		n := len(s)
		for i, _ := range s {
			m := s[n-1-i]
			if len(m) > 0 {
				c <- m
			}
		}
	}()
	return c
}

// BestScore returns the best score for a candidate. Repo names are weighted
// higher than owner names. They are then added together for the total.
func BestScore(term, candidate string) float64 {
	logger := log.WithFields(log.Fields{
		"term":      term,
		"candidate": candidate,
	})
	logger.Debug("Scoring term")
	split := reverse(strings.Split(term, "/"))
	candidates := reverse(strings.Split(candidate, "/"))

	var (
		a, b  string
		score float64
	)

	for a = range split {
		b = <-candidates
		score += smetrics.JaroWinkler(a, b, 0.7, 4)
	}
	for b = range candidates {
		score += smetrics.JaroWinkler(a, b, 0.7, 4)
	}
	score += 2 * smetrics.JaroWinkler(term, candidate, 0.7, 4)
	return score
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
	results := matches.ToStringArray()
	return results
}