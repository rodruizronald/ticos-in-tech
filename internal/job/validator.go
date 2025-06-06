package job

import (
	"errors"
	"slices"
	"strings"
	"time"
)

// validateEnum checks if a value is in the allowed enum values
func validateEnum(value string, validValues []string, fieldName string) error {
	if value == "" {
		return nil // Optional fields
	}

	if slices.Contains(validValues, value) {
		return nil
	}

	return &ValidationError{Field: fieldName}
}

// validateSearchRequest validates the search request parameters
func validateSearchRequest(req *SearchRequest) []string {
	var errs []string

	// Validate enum fields
	if err := validateEnum(req.ExperienceLevel, validExperienceLevels, "experience_level"); err != nil {
		errs = append(errs, err.Error())
	}

	if err := validateEnum(req.EmploymentType, validEmploymentTypes, "employment_type"); err != nil {
		errs = append(errs, err.Error())
	}

	if err := validateEnum(req.Location, validLocations, "location"); err != nil {
		errs = append(errs, err.Error())
	}

	if err := validateEnum(req.WorkMode, validWorkModes, "work_mode"); err != nil {
		errs = append(errs, err.Error())
	}

	// Validate date range - both must be provided if one is provided
	hasDateFrom := req.DateFrom != ""
	hasDateTo := req.DateTo != ""

	if hasDateFrom != hasDateTo {
		errs = append(errs, "both date_from and date_to must be provided together")
	}

	// Validate date format if provided
	if hasDateFrom && hasDateTo {
		dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			errs = append(errs, "date_from must be in YYYY-MM-DD format")
		}

		dateTo, err := time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			errs = append(errs, "date_to must be in YYYY-MM-DD format")
		}

		// Check date range if both dates are valid
		if err == nil && dateFrom.After(dateTo) {
			errs = append(errs, "date_from cannot be after date_to")
		}
	}

	return errs
}

// validateSearchParams ensures search parameters are within acceptable bounds
func validateSearchParams(sp *SearchParams) error {
	if strings.TrimSpace(sp.Query) == "" {
		return errors.New("search query cannot be empty")
	}

	if sp.Limit <= 0 {
		sp.Limit = 20 // Default limit
	}

	if sp.Limit > 100 {
		sp.Limit = 100 // Max limit to prevent abuse
	}

	if sp.Offset < 0 {
		sp.Offset = 0
	}

	// Validate date range if provided
	if sp.DateFrom != nil && sp.DateTo != nil {
		if sp.DateFrom.After(*sp.DateTo) {
			return errors.New("date from cannot be after date to")
		}
	}

	return nil
}
