package job

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SQL query constants
const (
	createJobQuery = `
        INSERT INTO jobs (
            company_id, title, description, experience_level, employment_type,
            location, work_mode, application_url, is_active, signature
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `

	getJobByIDQuery = `
        SELECT id, company_id, title, description, experience_level, employment_type,
               location, work_mode, application_url, is_active, signature, created_at, updated_at
        FROM jobs
        WHERE id = $1
    `

	updateJobQuery = `
        UPDATE jobs
        SET company_id = $1, title = $2, description = $3, experience_level = $4,
            employment_type = $5, location = $6, work_mode = $7, application_url = $8,
            is_active = $9, signature = $10, updated_at = NOW()
        WHERE id = $11
        RETURNING updated_at
    `

	deleteJobQuery = `DELETE FROM jobs WHERE id = $1`

	listJobsBaseQuery = `
        SELECT id, company_id, title, description, experience_level, employment_type,
               location, work_mode, application_url, is_active, signature, created_at, updated_at
        FROM jobs
        WHERE 1=1
    `
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

// Filter defines the available filters for job queries
type Filter struct {
	CompanyID       *int
	IsActive        *bool
	Location        *string
	WorkMode        *string
	ExperienceLevel *string
	EmploymentType  *string
}

// List retrieves all jobs from the database with optional filtering.
func (r *Repository) List(ctx context.Context, filter Filter) ([]*Job, error) {
	// Start with base query
	query := listJobsBaseQuery

	// Build query based on filters
	args := []any{}
	argCount := 1

	if filter.CompanyID != nil {
		query += fmt.Sprintf(" AND company_id = $%d", argCount)
		args = append(args, *filter.CompanyID)
		argCount++
	}

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filter.IsActive)
		argCount++
	}

	if filter.Location != nil {
		query += fmt.Sprintf(" AND location = $%d", argCount)
		args = append(args, *filter.Location)
		argCount++
	}

	if filter.WorkMode != nil {
		query += fmt.Sprintf(" AND work_mode = $%d", argCount)
		args = append(args, *filter.WorkMode)
		argCount++
	}

	if filter.ExperienceLevel != nil {
		query += fmt.Sprintf(" AND experience_level = $%d", argCount)
		args = append(args, *filter.ExperienceLevel)
		argCount++
	}

	if filter.EmploymentType != nil {
		query += fmt.Sprintf(" AND employment_type = $%d", argCount)
		args = append(args, *filter.EmploymentType)
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		job := &Job{}
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan job row: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating job rows: %w", err)
	}

	return jobs, nil
}
