package job_technology

import "fmt"

// ErrNotFound represents a job technology association not found error
type ErrNotFound struct {
	ID           int
	JobID        int
	TechnologyID int
}

func (e ErrNotFound) Error() string {
	if e.ID > 0 {
		return fmt.Sprintf("job technology with ID %d not found", e.ID)
	}
	return fmt.Sprintf("job technology association for job %d and technology %d not found", e.JobID, e.TechnologyID)
}

// IsNotFound checks if an error is a job technology not found error
func IsNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

// ErrDuplicate represents a duplicate job technology association error
type ErrDuplicate struct {
	JobID        int
	TechnologyID int
}

func (e ErrDuplicate) Error() string {
	return fmt.Sprintf("job technology association for job %d and technology %d already exists", e.JobID, e.TechnologyID)
}

// IsDuplicate checks if an error is a duplicate job technology error
func IsDuplicate(err error) bool {
	_, ok := err.(*ErrDuplicate)
	return ok
}
