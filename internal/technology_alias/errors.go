package technology_alias

import "fmt"

// ErrNotFound represents a technology alias not found error
type ErrNotFound struct {
    ID    int
    Alias string
}

func (e ErrNotFound) Error() string {
    if e.ID != 0 {
        return fmt.Sprintf("technology alias with ID %d not found", e.ID)
    }
    return fmt.Sprintf("technology alias with value %q not found", e.Alias)
}

// IsNotFound checks if an error is a technology alias not found error
func IsNotFound(err error) bool {
    _, ok := err.(*ErrNotFound)
    return ok
}

// ErrDuplicate represents a duplicate technology alias error
type ErrDuplicate struct {
    Alias string
}

func (e ErrDuplicate) Error() string {
    return fmt.Sprintf("technology alias %q already exists", e.Alias)
}

// IsDuplicate checks if an error is a duplicate technology alias error
func IsDuplicate(err error) bool {
    _, ok := err.(*ErrDuplicate)
    return ok
}