package job

import "fmt"

// ErrNotFound represents a job not found error
type ErrNotFound struct {
    ID int
}

func (e ErrNotFound) Error() string {
    return fmt.Sprintf("job with ID %d not found", e.ID)
}

// ErrDuplicate represents a duplicate job error
type ErrDuplicate struct {
    Signature string
}

func (e ErrDuplicate) Error() string {
    return fmt.Sprintf("job with signature %s already exists", e.Signature)
}

// IsNotFound checks if an error is a job not found error
func IsNotFound(err error) bool {
    _, ok := err.(*ErrNotFound)
    return ok
}

// IsDuplicate checks if an error is a duplicate job error
func IsDuplicate(err error) bool {
    _, ok := err.(*ErrDuplicate)
    return ok
}