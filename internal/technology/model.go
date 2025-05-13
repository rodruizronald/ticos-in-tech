package technology

import (
	"time"

	"github.com/rodruizronald/ticos-in-tech/internal/job_technology"
	"github.com/rodruizronald/ticos-in-tech/internal/technology_alias"
)

// Technology represents a technology skill (programming language, framework, tool, etc.)
// used in job postings.
type Technology struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Category  string    `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Relationships (not stored in database)
	Aliases []technology_alias.TechnologyAlias `json:"aliases,omitempty" db:"-"`
	Jobs    []job_technology.JobTechnology     `json:"jobs,omitempty" db:"-"`
}
