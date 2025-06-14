package jobs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database interface to support pgxpool and mocks
type Database interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

// Repository handles database operations for the Job model.
type Repository struct {
	db Database
}

// NewRepository creates a new Repository instance.
func NewRepository(db Database) *Repository {
	return &Repository{db: db}
}

// SearchJobsWithCount performs a full-text search and returns both results and total count
func (r *Repository) SearchJobsWithCount(ctx context.Context, params *SearchParams) ([]*JobWithCompany, int, error) {
	// Trim whitespace from query
	params.Query = strings.TrimSpace(params.Query)

	// Build additional WHERE conditions
	whereConditions := []string{}
	args := []any{params.Query}
	argCount := 2 // Starting at 2 because $1 is the search query

	// Add optional filters
	if params.ExperienceLevel != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("j.experience_level = $%d", argCount))
		args = append(args, *params.ExperienceLevel)
		argCount++
	}

	if params.EmploymentType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("j.employment_type = $%d", argCount))
		args = append(args, *params.EmploymentType)
		argCount++
	}

	if params.Location != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("j.location = $%d", argCount))
		args = append(args, *params.Location)
		argCount++
	}

	if params.WorkMode != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("j.work_mode = $%d", argCount))
		args = append(args, *params.WorkMode)
		argCount++
	}

	if params.Company != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("LOWER(c.name) LIKE LOWER($%d)", argCount))
		args = append(args, "%"+*params.Company+"%")
		argCount++
	}

	if params.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("j.created_at >= $%d", argCount))
		args = append(args, *params.DateFrom)
		argCount++
	}

	if params.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("j.created_at <= $%d", argCount))
		args = append(args, *params.DateTo)
		argCount++
	}

	// Build additional WHERE clause
	additionalWhere := ""
	if len(whereConditions) > 0 {
		additionalWhere = " AND " + strings.Join(whereConditions, " AND ")
	}

	// Build final search query with ordering and pagination
	searchQuery := searchJobsWithCountBaseQuery + additionalWhere +
		fmt.Sprintf(" ORDER BY j.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)

	// Add pagination parameters
	args = append(args, params.Limit, params.Offset)

	// Execute search query
	rows, err := r.db.Query(ctx, searchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*JobWithCompany
	var total int

	for rows.Next() {
		job := &JobWithCompany{}
		err = rows.Scan(
			&job.ID,
			&job.CompanyID,
			&job.Title,
			&job.Description,
			&job.ExperienceLevel,
			&job.EmploymentType,
			&job.Location,
			&job.WorkMode,
			&job.ApplicationURL,
			&job.IsActive,
			&job.Signature,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.CompanyName,
			&job.CompanyLogoURL,
			&total, // Window function gives us the same total for each row
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan job row: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating job rows: %w", err)
	}

	// If no results, total should be 0
	if len(jobs) == 0 {
		total = 0
	}

	return jobs, total, nil
}

// Create inserts a new job into the database.
func (r *Repository) Create(ctx context.Context, job *Job) error {
	err := r.db.QueryRow(
		ctx,
		createJobQuery,
		job.CompanyID,
		job.Title,
		job.Description,
		job.ExperienceLevel,
		job.EmploymentType,
		job.Location,
		job.WorkMode,
		job.ApplicationURL,
		job.IsActive,
		job.Signature,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation (duplicate job signature)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &DuplicateError{Signature: job.Signature}
		}
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

// GetByID retrieves a job by its ID.
func (r *Repository) GetByID(ctx context.Context, id int) (*Job, error) {
	job := &Job{}
	err := r.db.QueryRow(ctx, getJobByIDQuery, id).Scan(
		&job.ID,
		&job.CompanyID,
		&job.Title,
		&job.Description,
		&job.ExperienceLevel,
		&job.EmploymentType,
		&job.Location,
		&job.WorkMode,
		&job.ApplicationURL,
		&job.IsActive,
		&job.Signature,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{ID: id}
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

// Update updates an existing job in the database.
func (r *Repository) Update(ctx context.Context, job *Job) error {
	err := r.db.QueryRow(
		ctx,
		updateJobQuery,
		job.CompanyID,
		job.Title,
		job.Description,
		job.ExperienceLevel,
		job.EmploymentType,
		job.Location,
		job.WorkMode,
		job.ApplicationURL,
		job.IsActive,
		job.Signature,
		job.ID,
	).Scan(&job.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &NotFoundError{ID: job.ID}
		}

		// Check for unique constraint violation (duplicate job signature)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &DuplicateError{Signature: job.Signature}
		}

		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

// Delete removes a job from the database.
func (r *Repository) Delete(ctx context.Context, id int) error {
	commandTag, err := r.db.Exec(ctx, deleteJobQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &NotFoundError{ID: id}
	}

	return nil
}

// GetBySignature retrieves a job by its signature.
func (r *Repository) GetBySignature(ctx context.Context, signature string) (*Job, error) {
	job := &Job{}
	err := r.db.QueryRow(ctx, getJobBySignatureQuery, signature).Scan(
		&job.ID,
		&job.CompanyID,
		&job.Title,
		&job.Description,
		&job.ExperienceLevel,
		&job.EmploymentType,
		&job.Location,
		&job.WorkMode,
		&job.ApplicationURL,
		&job.IsActive,
		&job.Signature,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Signature: signature}
		}
		return nil, fmt.Errorf("failed to get job by signature: %w", err)
	}

	return job, nil
}
