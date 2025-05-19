package job

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_List(t *testing.T) {
	t.Parallel()
	now := time.Now()
	dbError := errors.New("database error")

	// Helper function to create a pointer to a value
	intPtr := func(i int) *int { return &i }
	boolPtr := func(b bool) *bool { return &b }
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name         string
		filter       Filter
		mockSetup    func(mock pgxmock.PgxPoolIface, filter Filter)
		checkResults func(t *testing.T, jobs []*Job, err error)
	}{
		{
			name:   "no filters",
			filter: Filter{},
			mockSetup: func(mock pgxmock.PgxPoolIface, _ Filter) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(listJobsBaseQuery + " ORDER BY created_at DESC")).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}).AddRow(
						1, 1, "Software Engineer", "Job description", "Mid-Level", "Full-Time",
						"San Francisco", "Remote", "https://example.com/apply", true, "job-signature-1", now, now,
					).AddRow(
						2, 2, "Product Manager", "Another description", "Senior", "Full-Time",
						"New York", "Hybrid", "https://example.com/apply2", true, "job-signature-2", now, now,
					))
			},
			checkResults: func(t *testing.T, jobs []*Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Len(t, jobs, 2)
				assert.Equal(t, 1, jobs[0].ID)
				assert.Equal(t, "Software Engineer", jobs[0].Title)
				assert.Equal(t, 2, jobs[1].ID)
				assert.Equal(t, "Product Manager", jobs[1].Title)
			},
		},
		{
			name: "filter by company ID",
			filter: Filter{
				CompanyID: intPtr(1),
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, filter Filter) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(listJobsBaseQuery + " AND company_id = $1 ORDER BY created_at DESC")).
					WithArgs(*filter.CompanyID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}).AddRow(
						1, 1, "Software Engineer", "Job description", "Mid-Level", "Full-Time",
						"San Francisco", "Remote", "https://example.com/apply", true, "job-signature-1", now, now,
					))
			},
			checkResults: func(t *testing.T, jobs []*Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Len(t, jobs, 1)
				assert.Equal(t, 1, jobs[0].ID)
				assert.Equal(t, 1, jobs[0].CompanyID)
				assert.Equal(t, "Software Engineer", jobs[0].Title)
			},
		},
		{
			name: "filter by multiple criteria",
			filter: Filter{
				IsActive:        boolPtr(true),
				WorkMode:        strPtr("Remote"),
				ExperienceLevel: strPtr("Mid-Level"),
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, filter Filter) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(
					listJobsBaseQuery+
						" AND is_active = $1"+
						" AND work_mode = $2"+
						" AND experience_level = $3"+
						" ORDER BY created_at DESC")).
					WithArgs(*filter.IsActive, *filter.WorkMode, *filter.ExperienceLevel).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}).AddRow(
						1, 1, "Software Engineer", "Job description", "Mid-Level", "Full-Time",
						"San Francisco", "Remote", "https://example.com/apply", true, "job-signature-1", now, now,
					))
			},
			checkResults: func(t *testing.T, jobs []*Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Len(t, jobs, 1)
				assert.Equal(t, "Mid-Level", jobs[0].ExperienceLevel)
				assert.Equal(t, "Remote", jobs[0].WorkMode)
				assert.True(t, jobs[0].IsActive)
			},
		},
		{
			name: "no results",
			filter: Filter{
				Location: strPtr("Antarctica"),
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, filter Filter) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(listJobsBaseQuery + " AND location = $1 ORDER BY created_at DESC")).
					WithArgs(*filter.Location).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}))
			},
			checkResults: func(t *testing.T, jobs []*Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Empty(t, jobs)
			},
		},
		{
			name: "database error",
			filter: Filter{
				EmploymentType: strPtr("Full-Time"),
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, filter Filter) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(listJobsBaseQuery + " AND employment_type = $1 ORDER BY created_at DESC")).
					WithArgs(*filter.EmploymentType).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, jobs []*Job, err error) {
				t.Helper()
				require.Error(t, err)
				assert.Nil(t, jobs)
				require.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.filter)

			jobs, err := repo.List(context.Background(), tt.filter)
			tt.checkResults(t, jobs, err)

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		job          *Job
		mockSetup    func(mock pgxmock.PgxPoolIface, job *Job)
		checkResults func(t *testing.T, result *Job, err error)
	}{
		{
			name: "successful creation",
			job: &Job{
				CompanyID:       1,
				Title:           "Software Engineer",
				Description:     "Job description",
				ExperienceLevel: "Mid-Level",
				EmploymentType:  "Full-Time",
				Location:        "San Francisco",
				WorkMode:        "Remote",
				ApplicationURL:  "https://example.com/apply",
				IsActive:        true,
				Signature:       "job-signature-1",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(createJobQuery)).
					WithArgs(
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
					).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "created_at", "updated_at",
					}).AddRow(1, now, now))
			},
			checkResults: func(t *testing.T, result *Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, now, result.CreatedAt)
				assert.Equal(t, now, result.UpdatedAt)
			},
		},
		{
			name: "duplicate job signature",
			job: &Job{
				CompanyID:       2,
				Title:           "Product Manager",
				Description:     "Job description",
				ExperienceLevel: "Senior",
				EmploymentType:  "Full-Time",
				Location:        "New York",
				WorkMode:        "Hybrid",
				ApplicationURL:  "https://example.com/apply2",
				IsActive:        true,
				Signature:       "duplicate-signature",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				pgErr := &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "idx_jobs_signature",
				}
				mock.ExpectQuery(regexp.QuoteMeta(createJobQuery)).
					WithArgs(
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
					).
					WillReturnError(pgErr)
			},
			checkResults: func(t *testing.T, _ *Job, err error) {
				t.Helper()
				require.Error(t, err)

				var duplicateErr *DuplicateError
				require.ErrorAs(t, err, &duplicateErr)
				assert.Equal(t, "duplicate-signature", duplicateErr.Signature)
			},
		},
		{
			name: "database error",
			job: &Job{
				CompanyID:       3,
				Title:           "Data Scientist",
				Description:     "Job description",
				ExperienceLevel: "Entry-Level",
				EmploymentType:  "Contract",
				Location:        "Chicago",
				WorkMode:        "On-Site",
				ApplicationURL:  "https://example.com/apply3",
				IsActive:        true,
				Signature:       "job-signature-3",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(createJobQuery)).
					WithArgs(
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
					).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, _ *Job, err error) {
				t.Helper()
				require.Error(t, err)
				require.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.job)

			err = repo.Create(context.Background(), tt.job)
			tt.checkResults(t, tt.job, err)

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	t.Parallel()
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		jobID        int
		mockSetup    func(mock pgxmock.PgxPoolIface, jobID int)
		checkResults func(t *testing.T, result *Job, err error)
	}{
		{
			name:  "job found",
			jobID: 1,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(getJobByIDQuery)).
					WithArgs(jobID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}).AddRow(
						1, 1, "Software Engineer", "Job description", "Mid-Level", "Full-Time",
						"San Francisco", "Remote", "https://example.com/apply", true, "job-signature-1", now, now,
					))
			},
			checkResults: func(t *testing.T, result *Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, 1, result.CompanyID)
				assert.Equal(t, "Software Engineer", result.Title)
				assert.Equal(t, "Job description", result.Description)
				assert.Equal(t, "Mid-Level", result.ExperienceLevel)
				assert.Equal(t, "Full-Time", result.EmploymentType)
				assert.Equal(t, "San Francisco", result.Location)
				assert.Equal(t, "Remote", result.WorkMode)
				assert.Equal(t, "https://example.com/apply", result.ApplicationURL)
				assert.True(t, result.IsActive)
				assert.Equal(t, "job-signature-1", result.Signature)
				assert.Equal(t, now, result.CreatedAt)
				assert.Equal(t, now, result.UpdatedAt)
			},
		},
		{
			name:  "job not found",
			jobID: 999,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(getJobByIDQuery)).
					WithArgs(jobID).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, result *Job, err error) {
				t.Helper()
				require.Error(t, err)
				assert.Nil(t, result)

				var notFoundErr *NotFoundError
				require.ErrorIs(t, err, dbError)
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name:  "database error",
			jobID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(getJobByIDQuery)).
					WithArgs(jobID).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *Job, err error) {
				t.Helper()
				require.Error(t, err)
				assert.Nil(t, result)
				require.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.jobID)

			result, err := repo.GetByID(context.Background(), tt.jobID)
			tt.checkResults(t, result, err)

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		job          *Job
		mockSetup    func(mock pgxmock.PgxPoolIface, job *Job)
		checkResults func(t *testing.T, result *Job, err error)
	}{
		{
			name: "successful update",
			job: &Job{
				ID:              1,
				CompanyID:       1,
				Title:           "Updated Software Engineer",
				Description:     "Updated job description",
				ExperienceLevel: "Senior",
				EmploymentType:  "Full-Time",
				Location:        "San Francisco",
				WorkMode:        "Hybrid",
				ApplicationURL:  "https://example.com/apply-updated",
				IsActive:        true,
				Signature:       "job-signature-1-updated",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(updateJobQuery)).
					WithArgs(
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
					).
					WillReturnRows(pgxmock.NewRows([]string{"updated_at"}).AddRow(now))
			},
			checkResults: func(t *testing.T, result *Job, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, now, result.UpdatedAt)
			},
		},
		{
			name: "job not found",
			job: &Job{
				ID:              999,
				CompanyID:       1,
				Title:           "Nonexistent Job",
				Description:     "Job description",
				ExperienceLevel: "Mid-Level",
				EmploymentType:  "Full-Time",
				Location:        "Remote",
				WorkMode:        "Remote",
				ApplicationURL:  "https://example.com/apply",
				IsActive:        true,
				Signature:       "nonexistent-signature",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(updateJobQuery)).
					WithArgs(
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
					).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, _ *Job, err error) {
				t.Helper()
				require.Error(t, err)

				var notFoundErr *NotFoundError
				require.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name: "duplicate job signature",
			job: &Job{
				ID:              2,
				CompanyID:       1,
				Title:           "Product Manager",
				Description:     "Job description",
				ExperienceLevel: "Senior",
				EmploymentType:  "Full-Time",
				Location:        "New York",
				WorkMode:        "Hybrid",
				ApplicationURL:  "https://example.com/apply2",
				IsActive:        true,
				Signature:       "duplicate-signature",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				pgErr := &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "idx_jobs_signature",
				}
				mock.ExpectQuery(regexp.QuoteMeta(updateJobQuery)).
					WithArgs(
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
					).
					WillReturnError(pgErr)
			},
			checkResults: func(t *testing.T, _ *Job, err error) {
				t.Helper()
				require.Error(t, err)

				var duplicateErr *DuplicateError
				require.ErrorAs(t, err, &duplicateErr)
				assert.Equal(t, "duplicate-signature", duplicateErr.Signature)
			},
		},
		{
			name: "database error",
			job: &Job{
				ID:              3,
				CompanyID:       1,
				Title:           "Error Job",
				Description:     "Job description",
				ExperienceLevel: "Mid-Level",
				EmploymentType:  "Full-Time",
				Location:        "Chicago",
				WorkMode:        "On-Site",
				ApplicationURL:  "https://example.com/apply3",
				IsActive:        true,
				Signature:       "error-signature",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, job *Job) {
				t.Helper()
				mock.ExpectQuery(regexp.QuoteMeta(updateJobQuery)).
					WithArgs(
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
					).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, _ *Job, err error) {
				t.Helper()
				require.Error(t, err)
				require.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.job)

			err = repo.Update(context.Background(), tt.job)
			tt.checkResults(t, tt.job, err)

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		id           int
		mockSetup    func(mock pgxmock.PgxPoolIface, id int)
		checkResults func(t *testing.T, err error)
	}{
		{
			name: "successful deletion",
			id:   1,
			mockSetup: func(mock pgxmock.PgxPoolIface, id int) {
				t.Helper()
				mock.ExpectExec(regexp.QuoteMeta(deleteJobQuery)).
					WithArgs(id).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			checkResults: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		{
			name: "job not found",
			id:   999,
			mockSetup: func(mock pgxmock.PgxPoolIface, id int) {
				t.Helper()
				mock.ExpectExec(regexp.QuoteMeta(deleteJobQuery)).
					WithArgs(id).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			checkResults: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)

				var notFoundErr *NotFoundError
				require.ErrorIs(t, err, dbError)
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name: "database error",
			id:   1,
			mockSetup: func(mock pgxmock.PgxPoolIface, id int) {
				t.Helper()
				mock.ExpectExec(regexp.QuoteMeta(deleteJobQuery)).
					WithArgs(id).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				require.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.id)

			err = repo.Delete(context.Background(), tt.id)
			tt.checkResults(t, err)

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
