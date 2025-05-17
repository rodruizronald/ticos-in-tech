package company

import (
	"time"

	"github.com/rodruizronald/ticos-in-tech/internal/job"
)

// Company represents a company that posts jobs on the platform.
type Company struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	LogoURL   string    `json:"logo_url" db:"logo_url"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Relationships (not stored in database)
	Jobs []job.Job `json:"jobs,omitempty" db:"-"`
}
