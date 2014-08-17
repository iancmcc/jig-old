package plan

import (
	"encoding/json"
	"io"
)

// Plan is the representation of the set of all repositories.
type Plan struct {
	Repos map[string]*Repo `json:"repos"`
}

// NewPlanFromJSON creates a Plan from JSON.
func NewPlanFromJSON(reader io.Reader) (*Plan, error) {
	var plan Plan
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&plan)
	if err != nil {
		return nil, err
	}
	plan.setRepoNames()
	return &plan, nil
}

// setRepoNames updates all the repos in the Plan with their full names.
func (p *Plan) setRepoNames() {
	for name, repo := range p.Repos {
		repo.FullName = name
	}
}
