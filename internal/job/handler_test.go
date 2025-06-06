package job

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SearchJobs(t *testing.T) {
	t.Parallel()

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		query        map[string]string
		mockSetup    func(t *testing.T, m *MockJobRepository)
		checkResults func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "successful search with basic query",
			query: map[string]string{
				"q":      "golang developer",
				"limit":  "20",
				"offset": "0",
			},
			mockSetup: func(t *testing.T, m *MockJobRepository) {
				t.Helper()
				expectedParams := &SearchParams{
					Query:  "golang developer",
					Limit:  20, // default limit
					Offset: 0,  // default offset
				}
				jobs := []*Job{
					{
						ID:              1,
						CompanyID:       100,
						Title:           "Senior Golang Developer",
						Description:     "We're looking for a senior golang developer",
						ExperienceLevel: "Senior",
						EmploymentType:  "Full-time",
						Location:        "Costa Rica",
						WorkMode:        "Remote",
						ApplicationURL:  "https://example.com/apply",
						IsActive:        true,
						Signature:       "job-sig-1",
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
					{
						ID:              2,
						CompanyID:       101,
						Title:           "Golang Backend Engineer",
						Description:     "Join our team as a golang backend engineer",
						ExperienceLevel: "Mid-level",
						EmploymentType:  "Full-time",
						Location:        "LATAM",
						WorkMode:        "Remote",
						ApplicationURL:  "https://example.com/apply2",
						IsActive:        true,
						Signature:       "job-sig-2",
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
				}
				m.On("Search", context.Background(), expectedParams).Return(jobs, nil)
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, w.Code)

				var response SearchResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Len(t, response.Data, 2)
				assert.Equal(t, "Senior Golang Developer", response.Data[0].Title)
				assert.Equal(t, "Golang Backend Engineer", response.Data[1].Title)

				assert.Equal(t, 2, response.Pagination.Total)
				assert.Equal(t, 20, response.Pagination.Limit)
				assert.Equal(t, 0, response.Pagination.Offset)
				assert.False(t, response.Pagination.HasMore)
			},
		},
		{
			name: "search with all filters",
			query: map[string]string{
				"q":                "java developer",
				"limit":            "50",
				"offset":           "10",
				"experience_level": "Senior",
				"employment_type":  "Full-time",
				"location":         "Costa Rica",
				"work_mode":        "Remote",
				"date_from":        "2024-01-01",
				"date_to":          "2024-12-31",
			},
			mockSetup: func(t *testing.T, m *MockJobRepository) {
				t.Helper()
				experienceLevel := "Senior"
				employmentType := "Full-time"
				location := "Costa Rica"
				workMode := "Remote"
				dateFrom, _ := time.Parse("2006-01-02", "2024-01-01")
				dateTo, _ := time.Parse("2006-01-02", "2024-12-31")

				expectedParams := &SearchParams{
					Query:           "java developer",
					Limit:           50,
					Offset:          10,
					ExperienceLevel: &experienceLevel,
					EmploymentType:  &employmentType,
					Location:        &location,
					WorkMode:        &workMode,
					DateFrom:        &dateFrom,
					DateTo:          &dateTo,
				}
				jobs := []*Job{
					{
						ID:              3,
						CompanyID:       102,
						Title:           "Senior Java Developer",
						Description:     "Senior Java Developer position",
						ExperienceLevel: "Senior",
						EmploymentType:  "Full-time",
						Location:        "Costa Rica",
						WorkMode:        "Remote",
						ApplicationURL:  "https://example.com/apply3",
						IsActive:        true,
						Signature:       "job-sig-3",
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
				}
				m.On("Search", context.Background(), expectedParams).Return(jobs, nil)
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, w.Code)

				var response SearchResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Len(t, response.Data, 1)
				assert.Equal(t, "Senior Java Developer", response.Data[0].Title)
				assert.Equal(t, 1, response.Pagination.Total)
				assert.Equal(t, 50, response.Pagination.Limit)
				assert.Equal(t, 10, response.Pagination.Offset)
				assert.False(t, response.Pagination.HasMore)
			},
		},
		{
			name: "search results equal to limit indicates more results",
			query: map[string]string{
				"q":     "python developer",
				"limit": "2",
			},
			mockSetup: func(t *testing.T, m *MockJobRepository) {
				t.Helper()
				expectedParams := &SearchParams{
					Query:  "python developer",
					Limit:  2,
					Offset: 0,
				}
				jobs := []*Job{
					{
						ID:              4,
						CompanyID:       103,
						Title:           "Python Developer 1",
						Description:     "Python Developer position",
						ExperienceLevel: "Mid-level",
						EmploymentType:  "Full-time",
						Location:        "LATAM",
						WorkMode:        "Remote",
						ApplicationURL:  "https://example.com/apply4",
						IsActive:        true,
						Signature:       "job-sig-4",
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
					{
						ID:              5,
						CompanyID:       104,
						Title:           "Python Developer 2",
						Description:     "Another Python Developer position",
						ExperienceLevel: "Senior",
						EmploymentType:  "Full-time",
						Location:        "LATAM",
						WorkMode:        "Remote",
						ApplicationURL:  "https://example.com/apply5",
						IsActive:        true,
						Signature:       "job-sig-5",
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
				}
				m.On("Search", context.Background(), expectedParams).Return(jobs, nil)
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, w.Code)

				var response SearchResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Len(t, response.Data, 2)
				assert.Equal(t, 2, response.Pagination.Total)
				assert.Equal(t, 2, response.Pagination.Limit)
				assert.True(t, response.Pagination.HasMore, "HasMore should be true when results equal limit")
			},
		},
		{
			name: "empty search results",
			query: map[string]string{
				"q":      "nonexistent technology",
				"limit":  "20",
				"offset": "0",
			},
			mockSetup: func(t *testing.T, m *MockJobRepository) {
				t.Helper()
				expectedParams := &SearchParams{
					Query:  "nonexistent technology",
					Limit:  20,
					Offset: 0,
				}

				m.On("Search", context.Background(), expectedParams).Return([]*Job{}, nil)
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, w.Code)

				var response SearchResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Empty(t, response.Data)
				assert.Equal(t, 0, response.Pagination.Total)
				assert.Equal(t, 20, response.Pagination.Limit)
				assert.Equal(t, 0, response.Pagination.Offset)
				assert.False(t, response.Pagination.HasMore)
			},
		},
		{
			name:  "missing required query parameter",
			query: map[string]string{}, // No 'q' parameter
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - request should fail before reaching repository
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeInvalidRequest, response.Error.Code)
				assert.Equal(t, "Invalid request parameters", response.Error.Message)
				assert.NotEmpty(t, response.Error.Details)
			},
		},
		{
			name: "invalid experience level",
			query: map[string]string{
				"q":                "developer",
				"experience_level": "Invalid Level",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "experience_level")
			},
		},
		{
			name: "invalid employment type",
			query: map[string]string{
				"q":               "developer",
				"employment_type": "Invalid Type",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "employment_type")
			},
		},
		{
			name: "invalid location",
			query: map[string]string{
				"q":        "developer",
				"location": "Invalid Location",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "location")
			},
		},
		{
			name: "invalid work mode",
			query: map[string]string{
				"q":         "developer",
				"work_mode": "Invalid Mode",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "work_mode")
			},
		},
		{
			name: "missing date_to when date_from is provided",
			query: map[string]string{
				"q":         "developer",
				"date_from": "2024-01-01",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "both date_from and date_to must be provided together")
			},
		},
		{
			name: "missing date_from when date_to is provided",
			query: map[string]string{
				"q":       "developer",
				"date_to": "2024-12-31",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "both date_from and date_to must be provided together")
			},
		},
		{
			name: "invalid date format",
			query: map[string]string{
				"q":         "developer",
				"date_from": "invalid-date",
				"date_to":   "2024-12-31",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "date_from must be in YYYY-MM-DD format")
			},
		},
		{
			name: "date_from after date_to",
			query: map[string]string{
				"q":         "developer",
				"date_from": "2024-12-31",
				"date_to":   "2024-01-01",
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.Contains(t, response.Error.Details[0], "date_from cannot be after date_to")
			},
		},
		{
			name: "repository search error",
			query: map[string]string{
				"q":      "developer",
				"limit":  "20",
				"offset": "0",
			},
			mockSetup: func(t *testing.T, m *MockJobRepository) {
				t.Helper()
				expectedParams := &SearchParams{
					Query:  "developer",
					Limit:  20,
					Offset: 0,
				}

				m.On("Search", context.Background(), expectedParams).Return([]*Job(nil), errors.New("database connection error"))
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusInternalServerError, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeSearchError, response.Error.Code)
				assert.Equal(t, "Failed to search jobs", response.Error.Message)
				assert.Equal(t, "database connection error", response.Error.Details[0])
			},
		},
		{
			name: "search with partial filters",
			query: map[string]string{
				"q":                "developer",
				"experience_level": "Senior",
				"location":         "Costa Rica",
				"limit":            "20",
				"offset":           "0",
				// employment_type and work_mode not provided
			},
			mockSetup: func(t *testing.T, m *MockJobRepository) {
				t.Helper()
				experienceLevel := "Senior"
				location := "Costa Rica"

				expectedParams := &SearchParams{
					Query:           "developer",
					Limit:           20,
					Offset:          0,
					ExperienceLevel: &experienceLevel,
					Location:        &location,
					// EmploymentType and WorkMode should be nil
				}
				jobs := []*Job{
					{
						ID:              6,
						CompanyID:       105,
						Title:           "Senior Developer",
						Description:     "Senior Developer position in Costa Rica",
						ExperienceLevel: "Senior",
						EmploymentType:  "Full-time",
						Location:        "Costa Rica",
						WorkMode:        "Hybrid",
						ApplicationURL:  "https://example.com/apply6",
						IsActive:        true,
						Signature:       "job-sig-6",
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
				}
				m.On("Search", context.Background(), expectedParams).Return(jobs, nil)
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, w.Code)

				var response SearchResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Len(t, response.Data, 1)
				assert.Equal(t, "Senior Developer", response.Data[0].Title)
				assert.Equal(t, "Costa Rica", response.Data[0].Location)
			},
		},
		{
			name: "multiple validation errors",
			query: map[string]string{
				"q":                "developer",
				"experience_level": "Invalid Level",
				"employment_type":  "Invalid Type",
				"location":         "Invalid Location",
				"work_mode":        "Invalid Mode",
				"date_from":        "2024-12-31",
				"date_to":          "2024-01-01", // date_from after date_to
			},
			mockSetup: func(t *testing.T, _ *MockJobRepository) {
				t.Helper()
				// No mock setup needed - validation should fail
			},
			checkResults: func(t *testing.T, w *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, w.Code)

				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, ErrCodeValidationError, response.Error.Code)
				assert.Equal(t, "Invalid search parameters", response.Error.Message)
				assert.GreaterOrEqual(t, len(response.Error.Details), 5, "Should have at least 5 validation errors")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock repository
			mockRepo := NewMockJobRepository(t)
			tt.mockSetup(t, mockRepo)

			// Create handler
			handler := NewHandler(mockRepo)

			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Build query string
			queryValues := url.Values{}
			for k, v := range tt.query {
				queryValues.Add(k, v)
			}

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/jobs?"+queryValues.Encode(), http.NoBody)
			require.NoError(t, err)
			req = req.WithContext(context.Background())
			c.Request = req

			// Call handler
			handler.SearchJobs(c)

			// Check results
			tt.checkResults(t, w)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
