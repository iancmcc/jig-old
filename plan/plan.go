package plan

import (
	"encoding/json"
	"io"
)

// Plan is the representation of the set of all repositories.
type Plan struct {
	Repos map[string]*Repo `json:"repos"`
}

func NewPlan() *Plan {
	return &Plan{make(map[string]*Repo)}
}

// NewPlanFromJSON creates a Plan from JSON.
func NewPlanFromJSON(reader io.Reader) (*Plan, error) {
	var plan Plan
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&plan)
	if err != nil {
		return nil, err
	}
	return cleanUpPlan(&plan)
}

// cleanUpPlan will create a brand new Plan with proper names and URLs
func cleanUpPlan(p *Plan) (*Plan, error) {
	newplan := NewPlan()
	for uri, repo := range p.Repos {
		newrepo, err := NewRepo(uri)
		if err != nil {
			return nil, err
		}
		if repo.FullName != "" {
			newrepo.FullName = repo.FullName
		}
		newrepo.RefSpec = repo.RefSpec
		newrepo.UserType = repo.UserType
		newplan.Repos[uri] = &newrepo
	}
	return newplan, nil
}

func (p *Plan) ToJSON(writer io.Writer) error {
	if indented, err := json.MarshalIndent(p, "", "    "); err != nil {
		return err
	} else {
		writer.Write(indented)
	}
	return nil
}
