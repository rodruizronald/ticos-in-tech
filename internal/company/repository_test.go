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
				Name:    "Test Company",
				LogoURL: "https://testcompany.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(createCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active).
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
				Name:    "Duplicate Company",
				LogoURL: "https://duplicate.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(createCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				var actualErr *ErrDuplicate
				assert.Error(t, err)
				assert.True(t, errors.As(err, &actualErr))
			},
		},
		{
			name: "database error",
			company: &Company{
				Name:    "Error Company",
				LogoURL: "https://error.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(createCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.True(t, errors.Is(err, dbError))
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
				assert.True(t, result.Active)
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
				assert.True(t, errors.As(err, &notFoundErr))
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
				assert.True(t, errors.Is(err, dbError))
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
				ID:      1,
				Name:    "Updated Company",
				LogoURL: "https://updated.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active, company.ID).
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
				ID:      999,
				Name:    "Nonexistent Company",
				LogoURL: "https://nonexistent.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active, company.ID).
					WillReturnError(pgx.ErrNoRows)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.Error(t, err)

				var notFoundErr *ErrNotFound
				assert.True(t, errors.As(err, &notFoundErr))
				assert.Equal(t, 999, notFoundErr.ID)
			},
		},
		{
			name: "duplicate company name",
			company: &Company{
				ID:      2,
				Name:    "Duplicate Company",
				LogoURL: "https://duplicate.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				pgErr := &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "companies_name_key",
				}
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active, company.ID).
					WillReturnError(pgErr)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.Error(t, err)

				var duplicateErr *ErrDuplicate
				assert.True(t, errors.As(err, &duplicateErr))
				assert.Equal(t, "Duplicate Company", duplicateErr.Name)
			},
		},
		{
			name: "database error",
			company: &Company{
				ID:      3,
				Name:    "Error Company",
				LogoURL: "https://error.com/logo.png",
				Active:  true,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, company *Company) {
				mock.ExpectQuery(regexp.QuoteMeta(updateCompanyQuery)).
					WithArgs(company.Name, company.LogoURL, company.Active, company.ID).
					WillReturnError(dbError)
			},
			checkResults: func(t *testing.T, result *Company, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, dbError))
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
				assert.True(t, errors.As(err, &notFoundErr))
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
				assert.True(t, errors.Is(err, dbError))
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
