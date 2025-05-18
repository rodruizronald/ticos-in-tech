package company

import "fmt"

// ErrNotFound represents a company not found error
type ErrNotFound struct {
	ID   int
	Name string
}

func (e ErrNotFound) Error() string {
	if e.ID > 0 {
		return fmt.Sprintf("company with ID %d not found", e.ID)
	}
	return fmt.Sprintf("company with name %s not found", e.Name)
}

// ErrDuplicate represents a duplicate company error
type ErrDuplicate struct {
	Name string
}

func (e ErrDuplicate) Error() string {
	return fmt.Sprintf("company with name %s already exists", e.Name)
}

// IsNotFound checks if an error is a company not found error
func IsNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

// IsDuplicate checks if an error is a duplicate company error
func IsDuplicate(err error) bool {
	_, ok := err.(*ErrDuplicate)
	return ok
}
