package company

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
	dbError := errors.New("database error")
	tests := []struct {
		name         string
		company      *Company
		mockSetup    func(mock pgxmock.PgxPoolIface, company *Company)
		checkResults func(t *testing.T, result *Company, err error)
	}{
		{
			name: "successful creation",
			company: &Company{
				Name:     "Test Company",
				LogoURL:  "https://testcompany.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(createCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, result.ID)
			},
		},
		{
			name: "duplicate company name",
			company: &Company{
				Name:     "Duplicate Company",
				LogoURL:  "https://duplicate.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(createCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				var actualErr *ErrDuplicate
				assert.Error(t, err)
				assert.ErrorAs(t, err, &actualErr)
			},
		},
		{
			name: "database error",
			company: &Company{
				Name:     "Error Company",
				LogoURL:  "https://error.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(createCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
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
			tt.mockSetup(mockDB, tt.company)

			err = repo.Create(context.Background(), tt.company)
			tt.checkResults(t, tt.company, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetByName(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")
	tests := []struct {
		name         string
		companyName  string
		mockSetup    func(mock pgxmock.PgxPoolIface, companyName string)
		checkResults func(t *testing.T, result *Company, err error)
	}{
		{
			name:        "company found",
			companyName: "Test Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}).AddRow(
						1, companyName, "https://testcompany.com/logo.png", true, now, now,
					))
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, "Test Company", result.Name)
				assert.Equal(t, "https://testcompany.com/logo.png", result.LogoURL)
				assert.True(t, result.IsActive)
				assert.Equal(t, now, result.CreatedAt)
				assert.Equal(t, now, result.UpdatedAt)
			},
		},
		{
			name:        "company not found",
			companyName: "Nonexistent Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)

				var notFoundErr *ErrNotFound
				assert.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, "Nonexistent Company", notFoundErr.Name)
			},
		},
		{
			name:        "database error",
			companyName: "Error Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
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
			tt.mockSetup(mockDB, tt.companyName)

			result, err := repo.GetByName(context.Background(), tt.companyName)
			tt.checkResults(t, result, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Update(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		company      *Company
		mockSetup    func(mock pgxmock.PgxPoolIface, company *Company)
		checkResults func(t *testing.T, result *Company, err error)
	}{
		{
			name: "successful update",
			company: &Company{
				ID:       1,
				Name:     "Updated Company",
				LogoURL:  "https://updated.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive, company.ID).
					WillReturnRows(pgxmock.NewRows([]string{"updated_at"}).AddRow(now))
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.NoError(t, err)
				assert.Equal(t, now, result.UpdatedAt)
			},
		},
		{
			name: "company not found",
			company: &Company{
				ID:       999,
				Name:     "Nonexistent Company",
				LogoURL:  "https://nonexistent.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive, company.ID).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.Error(t, err)

				var notFoundErr *ErrNotFound
				assert.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name: "duplicate company name",
			company: &Company{
				ID:       2,
				Name:     "Duplicate Company",
				LogoURL:  "https://duplicate.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				pgErr := &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "companies_name_key",
				}
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive, company.ID).
					WillReturnError(pgErr)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.Error(t, err)

				var duplicateErr *ErrDuplicate
				assert.ErrorAs(t, err, &duplicateErr)
				assert.Equal(t, "Duplicate Company", duplicateErr.Name)
			},
		},
		{
			name: "database error",
			company: &Company{
				ID:       3,
				Name:     "Error Company",
				LogoURL:  "https://error.com/logo.png",
				IsActive: true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.IsActive, company.ID).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
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
			tt.mockSetup(mockDB, tt.company)

			err = repo.Update(context.Background(), tt.company)
			tt.checkResults(t, tt.company, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		companyID    int
		mockSetup    func(mock pgxmock.PgxPoolIface, companyID int)
		checkResults func(t *testing.T, err error)
	}{
		{
			name:      "successful deletion",
			companyID: 1,
			mockSetup: func(mock pgxmock.PgxPoolIface, companyID int) {
				mock.ExpectExec(regexp.QuoteMeta(deleteCompanyQuery)).
					WithArgs(companyID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			checkResults: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "company not found",
			companyID: 999,
			mockSetup: func(mock pgxmock.PgxPoolIface, companyID int) {
				mock.ExpectExec(regexp.QuoteMeta(deleteCompanyQuery)).
					WithArgs(companyID).
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
			name:      "database error",
			companyID: 2,
			mockSetup: func(mock pgxmock.PgxPoolIface, companyID int) {
				mock.ExpectExec(regexp.QuoteMeta(deleteCompanyQuery)).
					WithArgs(companyID).
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
			tt.mockSetup(mockDB, tt.companyID)

			err = repo.Delete(context.Background(), tt.companyID)
			tt.checkResults(t, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_List(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		mockSetup    func(mock pgxmock.PgxPoolIface)
		checkResults func(t *testing.T, companies []*Company, err error)
	}{
		{
			name: "successful listing with results",
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(listCompaniesQuery)).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}).AddRow(
						1, "Company A", "https://example.com/logo1.png", true, now, now,
					).AddRow(
						2, "Company B", "https://example.com/logo2.png", false, now, now,
					))
			},
			checkResults: func(t *testing.T, companies []*Company, err error) {
				assert.NoError(t, err)
				assert.Len(t, companies, 2)

				assert.Equal(t, 1, companies[0].ID)
				assert.Equal(t, "Company A", companies[0].Name)
				assert.Equal(t, "https://example.com/logo1.png", companies[0].LogoURL)
				assert.True(t, companies[0].IsActive)

				assert.Equal(t, 2, companies[1].ID)
				assert.Equal(t, "Company B", companies[1].Name)
				assert.Equal(t, "https://example.com/logo2.png", companies[1].LogoURL)
				assert.False(t, companies[1].IsActive)
			},
		},
		{
			name: "successful listing with no results",
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(listCompaniesQuery)).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}))
			},
			checkResults: func(t *testing.T, companies []*Company, err error) {
				assert.NoError(t, err)
				assert.Empty(t, companies)
			},
		},
		{
			name: "database error",
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(listCompaniesQuery)).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, companies []*Company, err error) {
				assert.Error(t, err)
				assert.Nil(t, companies)
				assert.ErrorIs(t, err, dbError)
			},
		},
		{
			name: "scan error",
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// Return mismatched column count to cause scan error
				mock.ExpectQuery(regexp.QuoteMeta(listCompaniesQuery)).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", // Missing columns to cause scan error
					}).AddRow(
						1, "Company A",
					))
			},
			checkResults: func(t *testing.T, companies []*Company, err error) {
				assert.Error(t, err)
				assert.Nil(t, companies)
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
			tt.mockSetup(mockDB)

			companies, err := repo.List(context.Background())
			tt.checkResults(t, companies, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetWithJobs(t *testing.T) {
	now := time.Now()
	dbError := errors.New("database error")

	tests := []struct {
		name         string
		companyName  string
		mockSetup    func(mock pgxmock.PgxPoolIface, companyName string)
		checkResults func(t *testing.T, company *Company, err error)
	}{
		{
			name:        "successful retrieval with jobs",
			companyName: "Test Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				// First query to get the company
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}).AddRow(
						1, companyName, "https://example.com/logo.png", true, now, now,
					))

				// Second query to get the jobs
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyJobsQuery)).
					WithArgs(1).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}).AddRow(
						101, 1, "Software Engineer", "Job description", "Mid-Level", "Full-Time",
						"San Francisco", "Remote", "https://example.com/apply", true, "job-signature-1", now, now,
					).AddRow(
						102, 1, "Product Manager", "Another description", "Senior", "Full-Time",
						"New York", "Hybrid", "https://example.com/apply2", true, "job-signature-2", now, now,
					))
			},
			checkResults: func(t *testing.T, company *Company, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, company)
				assert.Equal(t, 1, company.ID)
				assert.Equal(t, "Test Company", company.Name)
				assert.Equal(t, "https://example.com/logo.png", company.LogoURL)
				assert.True(t, company.IsActive)

				// Check jobs
				assert.Len(t, company.Jobs, 2)
				assert.Equal(t, 101, company.Jobs[0].ID)
				assert.Equal(t, "Software Engineer", company.Jobs[0].Title)
				assert.Equal(t, 102, company.Jobs[1].ID)
				assert.Equal(t, "Product Manager", company.Jobs[1].Title)
			},
		},
		{
			name:        "company not found",
			companyName: "Nonexistent Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, company *Company, err error) {
				assert.Error(t, err)
				assert.Nil(t, company)

				var notFoundErr *ErrNotFound
				assert.ErrorAs(t, err, &notFoundErr)
				assert.Equal(t, "Nonexistent Company", notFoundErr.Name)
			},
		},
		{
			name:        "company found but error fetching jobs",
			companyName: "Test Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				// First query to get the company
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}).AddRow(
						1, companyName, "https://example.com/logo.png", true, now, now,
					))

				// Second query to get jobs returns error
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyJobsQuery)).
					WithArgs(1).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, company *Company, err error) {
				assert.Error(t, err)
				assert.Nil(t, company)
				assert.ErrorIs(t, err, dbError)
			},
		},
		{
			name:        "company found with no jobs",
			companyName: "Test Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				// First query to get the company
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}).AddRow(
						1, companyName, "https://example.com/logo.png", true, now, now,
					))

				// Second query to get jobs returns empty result
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyJobsQuery)).
					WithArgs(1).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", "description", "experience_level", "employment_type",
						"location", "work_mode", "application_url", "is_active", "signature", "created_at", "updated_at",
					}))
			},
			checkResults: func(t *testing.T, company *Company, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, company)
				assert.Equal(t, 1, company.ID)
				assert.Equal(t, "Test Company", company.Name)
				assert.Empty(t, company.Jobs)
			},
		},
		{
			name:        "scan error in jobs",
			companyName: "Test Company",
			mockSetup: func(mock pgxmock.PgxPoolIface, companyName string) {
				// First query to get the company
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyByNameQuery)).
					WithArgs(companyName).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "name", "logo_url", "active", "created_at", "updated_at",
					}).AddRow(
						1, companyName, "https://example.com/logo.png", true, now, now,
					))

				// Second query returns mismatched columns to cause scan error
				mock.ExpectQuery(regexp.QuoteMeta(getCompanyJobsQuery)).
					WithArgs(1).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "company_id", "title", // Missing columns to cause scan error
					}).AddRow(
						101, 1, "Software Engineer",
					))
			},
			checkResults: func(t *testing.T, company *Company, err error) {
				assert.Error(t, err)
				assert.Nil(t, company)
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
			tt.mockSetup(mockDB, tt.companyName)

			company, err := repo.GetWithJobs(context.Background(), tt.companyName)
			tt.checkResults(t, company, err)

			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
