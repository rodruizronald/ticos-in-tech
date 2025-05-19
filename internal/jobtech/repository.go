package jobtech

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SQL query constants
const (
	createJobTechnologyQuery = `
        INSERT INTO job_technologies (job_id, technology_id, is_primary, is_required)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `

	getJobTechnologyByJobAndTechQuery = `
        SELECT id, job_id, technology_id, is_primary, is_required, created_at
        FROM job_technologies
        WHERE job_id = $1 AND technology_id = $2
    `

	updateJobTechnologyQuery = `
        UPDATE job_technologies
        SET is_primary = $1, is_required = $2
        WHERE id = $3
    `

	deleteJobTechnologyQuery = `DELETE FROM job_technologies WHERE id = $1`

	listJobTechnologiesByJobQuery = `
        SELECT id, job_id, technology_id, is_primary, is_required, created_at
        FROM job_technologies
        WHERE job_id = $1
        ORDER BY is_primary DESC, id
    `

	listJobTechnologiesByTechnologyQuery = `
        SELECT id, job_id, technology_id, is_primary, is_required, created_at
        FROM job_technologies
        WHERE technology_id = $1
        ORDER BY created_at DESC
    `
)

// Database interface to support pgxpool and mocks
type Database interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

// Repository handles database operations for the JobTechnology model.
type Repository struct {
	db Database
}

// NewRepository creates a new Repository instance.
func NewRepository(db Database) *Repository {
	return &Repository{db: db}
}

// Create inserts a new job-technology association into the database.
func (r *Repository) Create(ctx context.Context, jobTech *JobTechnology) error {
	err := r.db.QueryRow(
		ctx,
		createJobTechnologyQuery,
		jobTech.JobID,
		jobTech.TechnologyID,
		jobTech.IsPrimary,
		jobTech.IsRequired,
	).Scan(&jobTech.ID, &jobTech.CreatedAt)

	if err != nil {
		// Check for unique constraint violation (duplicate job-technology association)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &DuplicateError{
				JobID:        jobTech.JobID,
				TechnologyID: jobTech.TechnologyID,
			}
		}
		return fmt.Errorf("failed to create job technology association: %w", err)
	}

	return nil
}

// GetByJobAndTechnology retrieves a job-technology association by job ID and technology ID.
func (r *Repository) GetByJobAndTechnology(ctx context.Context, jobID, technologyID int) (*JobTechnology, error) {
	jobTech := &JobTechnology{}
	err := r.db.QueryRow(ctx, getJobTechnologyByJobAndTechQuery, jobID, technologyID).Scan(
		&jobTech.ID,
		&jobTech.JobID,
		&jobTech.TechnologyID,
		&jobTech.IsPrimary,
		&jobTech.IsRequired,
		&jobTech.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{
				JobID:        jobID,
				TechnologyID: technologyID,
			}
		}
		return nil, fmt.Errorf("failed to get job technology association: %w", err)
	}

	return jobTech, nil
}

// Update updates an existing job-technology association in the database.
func (r *Repository) Update(ctx context.Context, jobTech *JobTechnology) error {
	commandTag, err := r.db.Exec(
		ctx,
		updateJobTechnologyQuery,
		jobTech.IsPrimary,
		jobTech.IsRequired,
		jobTech.ID,
	)

	if err != nil {
		// Check for unique constraint violation (duplicate job-technology association)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &DuplicateError{
				JobID:        jobTech.JobID,
				TechnologyID: jobTech.TechnologyID,
			}
		}
		return fmt.Errorf("failed to update job technology association: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &NotFoundError{ID: jobTech.ID}
	}

	return nil
}

// Delete removes a job-technology association from the database.
func (r *Repository) Delete(ctx context.Context, id int) error {
	commandTag, err := r.db.Exec(ctx, deleteJobTechnologyQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete job technology association: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return &NotFoundError{ID: id}
	}

	return nil
}

// ListByJob retrieves all technology associations for a specific job.
func (r *Repository) ListByJob(ctx context.Context, jobID int) ([]*JobTechnology, error) {
	rows, err := r.db.Query(ctx, listJobTechnologiesByJobQuery, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to list job technologies: %w", err)
	}
	defer rows.Close()

	var jobTechnologies []*JobTechnology
	for rows.Next() {
		jobTech := &JobTechnology{}
		err = rows.Scan(
			&jobTech.ID,
			&jobTech.JobID,
			&jobTech.TechnologyID,
			&jobTech.IsPrimary,
			&jobTech.IsRequired,
			&jobTech.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job technology row: %w", err)
		}
		jobTechnologies = append(jobTechnologies, jobTech)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating job technology rows: %w", err)
	}

	return jobTechnologies, nil
}

// ListByTechnology retrieves all job associations for a specific technology.
func (r *Repository) ListByTechnology(ctx context.Context, technologyID int) ([]*JobTechnology, error) {
	rows, err := r.db.Query(ctx, listJobTechnologiesByTechnologyQuery, technologyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list technology jobs: %w", err)
	}
	defer rows.Close()

	var jobTechnologies []*JobTechnology
	for rows.Next() {
		jobTech := &JobTechnology{}
		err = rows.Scan(
			&jobTech.ID,
			&jobTech.JobID,
			&jobTech.TechnologyID,
			&jobTech.IsPrimary,
			&jobTech.IsRequired,
			&jobTech.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan technology job row: %w", err)
		}
		jobTechnologies = append(jobTechnologies, jobTech)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating technology job rows: %w", err)
	}

	return jobTechnologies, nil
}
