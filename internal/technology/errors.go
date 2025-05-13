package technology

import "fmt"

// ErrNotFound represents a technology not found error
type ErrNotFound struct {
    ID   int
    Name string
}

func (e ErrNotFound) Error() string {
    if e.ID > 0 {
        return fmt.Sprintf("technology with ID %d not found", e.ID)
    }
    return fmt.Sprintf("technology with name %s not found", e.Name)
}

// ErrDuplicate represents a duplicate technology error
type ErrDuplicate struct {
    Name string
}

func (e ErrDuplicate) Error() string {
    return fmt.Sprintf("technology with name %s already exists", e.Name)
}

// IsNotFound checks if an error is a technology not found error
func IsNotFound(err error) bool {
    _, ok := err.(*ErrNotFound)
    return ok
}

// IsDuplicate checks if an error is a duplicate technology error
func IsDuplicate(err error) bool {
    _, ok := err.(*ErrDuplicate)
    return ok
}