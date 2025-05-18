package job_technology

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

func TestRepository_Create(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		jobTech      *JobTechnology
		mockSetup    func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology)
		checkResults func(t *testing.T, result *JobTechnology, err error)
	}{
		{
			name: "successful creation",
			jobTech: &JobTechnology{
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				mock.ExpectQuery(regexp.QuoteMeta(createJobTechnologyQuery)).
					WithArgs(
						jobTech.JobID,
						jobTech.TechnologyID,
						jobTech.IsPrimary,
						jobTech.IsRequired,
					).
					WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))
			},
			checkResults: func(t *testing.T, result *JobTechnology, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, now, result.CreatedAt)
			},
		},
		{
			name: "duplicate job-technology association",
			jobTech: &JobTechnology{
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				pgErr := &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "job_technologies_job_id_technology_id_key",
				}
				mock.ExpectQuery(regexp.QuoteMeta(createJobTechnologyQuery)).
					WithArgs(
						jobTech.JobID,
						jobTech.TechnologyID,
						jobTech.IsPrimary,
						jobTech.IsRequired,
					).
					WillReturnError(pgErr)
			},
			checkResults: func(t *testing.T, result *JobTechnology, err error) {
				assert.Error(t, err)

				var duplicateErr *ErrDuplicate
				assert.ErrorAs(t, err, &duplicateErr)
				assert.Equal(t, 1, duplicateErr.JobID)
				assert.Equal(t, 2, duplicateErr.TechnologyID)
			},
		},
		{
			name: "database error",
			jobTech: &JobTechnology{
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				mock.ExpectQuery(regexp.QuoteMeta(createJobTechnologyQuery)).
					WithArgs(
						jobTech.JobID,
						jobTech.TechnologyID,
						jobTech.IsPrimary,
						jobTech.IsRequired,
					).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *JobTechnology, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.jobTech)

			err = repo.Create(context.Background(), tt.jobTech)
			tt.checkResults(t, tt.jobTech, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetByJobAndTechnology(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		jobID        int
		techID       int
		mockSetup    func(mock pgxmock.PgxPoolIface, jobID, techID int)
		checkResults func(t *testing.T, result *JobTechnology, err error)
	}{
		{
			name:   "association found",
			jobID:  1,
			techID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID, techID int) {
				mock.ExpectQuery(regexp.QuoteMeta(getJobTechnologyByJobAndTechQuery)).
					WithArgs(jobID, techID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", "technology_id", "is_primary", "is_required", "created_at",
					}).AddRow(
						1, jobID, techID, true, true, now,
					))
			},
			checkResults: func(t *testing.T, result *JobTechnology, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, 1, result.JobID)
				assert.Equal(t, 2, result.TechnologyID)
				assert.True(t, result.IsPrimary)
				assert.True(t, result.IsRequired)
				assert.Equal(t, now, result.CreatedAt)
			},
		},
		{
			name:   "association not found",
			jobID:  999,
			techID: 888,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID, techID int) {
				mock.ExpectQuery(regexp.QuoteMeta(getJobTechnologyByJobAndTechQuery)).
					WithArgs(jobID, techID).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, result *JobTechnology, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)

				var notFoundErr *ErrNotFound
				assert.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, 999, notFoundErr.JobID)
				assert.Equal(t, 888, notFoundErr.TechnologyID)
			},
		},
		{
			name:   "database error",
			jobID:  1,
			techID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID, techID int) {
				mock.ExpectQuery(regexp.QuoteMeta(getJobTechnologyByJobAndTechQuery)).
					WithArgs(jobID, techID).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *JobTechnology, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.jobID, tt.techID)

			result, err := repo.GetByJobAndTechnology(context.Background(), tt.jobID, tt.techID)
			tt.checkResults(t, result, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Update(t *testing.T) {
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		jobTech      *JobTechnology
		mockSetup    func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology)
		checkResults func(t *testing.T, err error)
	}{
		{
			name: "successful update",
			jobTech: &JobTechnology{
				ID:           1,
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   false,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				mock.ExpectExec(regexp.QuoteMeta(updateJobTechnologyQuery)).
					WithArgs(
						jobTech.IsPrimary,
						jobTech.IsRequired,
						jobTech.ID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			checkResults: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "job technology not found",
			jobTech: &JobTechnology{
				ID:           999,
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   false,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				mock.ExpectExec(regexp.QuoteMeta(updateJobTechnologyQuery)).
					WithArgs(
						jobTech.IsPrimary,
						jobTech.IsRequired,
						jobTech.ID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			checkResults: func(t *testing.T, err error) {
				assert.Error(t, err)

				var notFoundErr *ErrNotFound
				assert.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name: "duplicate constraint error",
			jobTech: &JobTechnology{
				ID:           1,
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   false,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				pgErr := &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "job_technologies_job_id_technology_id_key",
				}
				mock.ExpectExec(regexp.QuoteMeta(updateJobTechnologyQuery)).
					WithArgs(
						jobTech.IsPrimary,
						jobTech.IsRequired,
						jobTech.ID,
					).
					WillReturnError(pgErr)
			},
			checkResults: func(t *testing.T, err error) {
				assert.Error(t, err)

				var duplicateErr *ErrDuplicate
				assert.ErrorAs(t, err, &duplicateErr)
				assert.Equal(t, 1, duplicateErr.JobID)
				assert.Equal(t, 2, duplicateErr.TechnologyID)
			},
		},
		{
			name: "database error",
			jobTech: &JobTechnology{
				ID:           1,
				JobID:        1,
				TechnologyID: 2,
				IsPrimary:    true,
				IsRequired:   false,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, jobTech *JobTechnology) {
				mock.ExpectExec(regexp.QuoteMeta(updateJobTechnologyQuery)).
					WithArgs(
						jobTech.IsPrimary,
						jobTech.IsRequired,
						jobTech.ID,
					).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.jobTech)

			err = repo.Update(context.Background(), tt.jobTech)
			tt.checkResults(t, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		id           int
		mockSetup    func(mock pgxmock.PgxPoolIface, id int)
		checkResults func(t *testing.T, err error)
	}{
		{
			name: "successful delete",
			id:   1,
			mockSetup: func(mock pgxmock.PgxPoolIface, id int) {
				mock.ExpectExec(regexp.QuoteMeta(deleteJobTechnologyQuery)).
					WithArgs(id).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			checkResults: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "job technology not found",
			id:   999,
			mockSetup: func(mock pgxmock.PgxPoolIface, id int) {
				mock.ExpectExec(regexp.QuoteMeta(deleteJobTechnologyQuery)).
					WithArgs(id).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			checkResults: func(t *testing.T, err error) {
				assert.Error(t, err)

				var notFoundErr *ErrNotFound
				assert.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name: "database error",
			id:   1,
			mockSetup: func(mock pgxmock.PgxPoolIface, id int) {
				mock.ExpectExec(regexp.QuoteMeta(deleteJobTechnologyQuery)).
					WithArgs(id).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, dbError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.id)

			err = repo.Delete(context.Background(), tt.id)
			tt.checkResults(t, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_ListByJob(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		jobID        int
		mockSetup    func(mock pgxmock.PgxPoolIface, jobID int)
		checkResults func(t *testing.T, results []*JobTechnology, err error)
	}{
		{
			name:  "successful listing with results",
			jobID: 1,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByJobQuery)).
					WithArgs(jobID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", "technology_id", "is_primary", "is_required", "created_at",
					}).AddRow(
						1, jobID, 2, true, true, now,
					).AddRow(
						2, jobID, 3, false, true, now,
					))
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.NoError(t, err)
				assert.Len(t, results, 2)

				assert.Equal(t, 1, results[0].ID)
				assert.Equal(t, 1, results[0].JobID)
				assert.Equal(t, 2, results[0].TechnologyID)
				assert.True(t, results[0].IsPrimary)
				assert.True(t, results[0].IsRequired)

				assert.Equal(t, 2, results[1].ID)
				assert.Equal(t, 1, results[1].JobID)
				assert.Equal(t, 3, results[1].TechnologyID)
				assert.False(t, results[1].IsPrimary)
				assert.True(t, results[1].IsRequired)
			},
		},
		{
			name:  "successful listing with no results",
			jobID: 999,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByJobQuery)).
					WithArgs(jobID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", "technology_id", "is_primary", "is_required", "created_at",
					}))
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.NoError(t, err)
				assert.Empty(t, results)
			},
		},
		{
			name:  "database error",
			jobID: 1,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByJobQuery)).
					WithArgs(jobID).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.Error(t, err)
				assert.Nil(t, results)
				assert.ErrorIs(t, err, dbError)
			},
		},
		{
			name:  "scan error",
			jobID: 1,
			mockSetup: func(mock pgxmock.PgxPoolIface, jobID int) {
				// Return mismatched column count to cause scan error
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByJobQuery)).
					WithArgs(jobID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", // Missing columns to cause scan error
					}).AddRow(
						1, jobID,
					))
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.Error(t, err)
				assert.Nil(t, results)
				assert.Contains(t, err.Error(), "scan")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.jobID)

			results, err := repo.ListByJob(context.Background(), tt.jobID)
			tt.checkResults(t, results, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_ListByTechnology(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		techID       int
		mockSetup    func(mock pgxmock.PgxPoolIface, techID int)
		checkResults func(t *testing.T, results []*JobTechnology, err error)
	}{
		{
			name:   "successful listing with results",
			techID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, techID int) {
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByTechnologyQuery)).
					WithArgs(techID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", "technology_id", "is_primary", "is_required", "created_at",
					}).AddRow(
						1, 1, techID, true, true, now,
					).AddRow(
						3, 2, techID, false, true, now,
					))
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.NoError(t, err)
				assert.Len(t, results, 2)

				assert.Equal(t, 1, results[0].ID)
				assert.Equal(t, 1, results[0].JobID)
				assert.Equal(t, 2, results[0].TechnologyID)
				assert.True(t, results[0].IsPrimary)
				assert.True(t, results[0].IsRequired)

				assert.Equal(t, 3, results[1].ID)
				assert.Equal(t, 2, results[1].JobID)
				assert.Equal(t, 2, results[1].TechnologyID)
				assert.False(t, results[1].IsPrimary)
				assert.True(t, results[1].IsRequired)
			},
		},
		{
			name:   "successful listing with no results",
			techID: 999,
			mockSetup: func(mock pgxmock.PgxPoolIface, techID int) {
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByTechnologyQuery)).
					WithArgs(techID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", "technology_id", "is_primary", "is_required", "created_at",
					}))
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.NoError(t, err)
				assert.Empty(t, results)
			},
		},
		{
			name:   "database error",
			techID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, techID int) {
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByTechnologyQuery)).
					WithArgs(techID).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.Error(t, err)
				assert.Nil(t, results)
				assert.ErrorIs(t, err, dbError)
			},
		},
		{
			name:   "scan error",
			techID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, techID int) {
				// Return mismatched column count to cause scan error
				mock.ExpectQuery(regexp.QuoteMeta(listJobTechnologiesByTechnologyQuery)).
					WithArgs(techID).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "job_id", // Missing columns to cause scan error
					}).AddRow(
						1, 1,
					))
			},
			checkResults: func(t *testing.T, results []*JobTechnology, err error) {
				assert.Error(t, err)
				assert.Nil(t, results)
				assert.Contains(t, err.Error(), "scan")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			repo := NewRepository(mockDB)
			tt.mockSetup(mockDB, tt.techID)

			results, err := repo.ListByTechnology(context.Background(), tt.techID)
			tt.checkResults(t, results, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
