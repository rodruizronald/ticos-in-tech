package job_technology

import (
	"time"
)

// JobTechnology represents the association between a job and a technology,
// including additional metadata about the relationship.
type JobTechnology struct {
	ID           int       `json:"id" db:"id"`
	JobID        int       `json:"job_id" db:"job_id"`
	TechnologyID int       `json:"technology_id" db:"technology_id"`
	IsPrimary    bool      `json:"is_primary" db:"is_primary"`
	IsRequired   bool      `json:"is_required" db:"is_required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
