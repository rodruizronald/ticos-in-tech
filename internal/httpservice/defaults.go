package httpservice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Constants for error codes and messages
const (
	ErrCodeInternalError   = "INTERNAL_ERROR"
	ErrCodeInvalidRequest  = "INVALID_REQUEST"
	ErrCodeValidationError = "VALIDATION_ERROR"
	ErrCodeSearchError     = "SEARCH_ERROR"
)

// DefaultRequestParser - GENERIC IMPLEMENTATION that consumers can use
type DefaultRequestParser[T SearchRequest] struct {
	createRequest func() T // Factory function to create new request instance
}

func NewDefaultRequestParser[T SearchRequest](createRequest func() T) RequestParser[T] {
	return &DefaultRequestParser[T]{createRequest: createRequest}
}

func (p *DefaultRequestParser[T]) ParseSearchRequest(c *gin.Context) (T, error) {
	req := p.createRequest()
	if err := c.ShouldBindQuery(req); err != nil {
		var zero T
		return zero, &RequestParseError{Err: err}
	}
	return req, nil
}

// DefaultResponseBuilder - GENERIC IMPLEMENTATION that consumers can use
type DefaultResponseBuilder[TResult SearchResult, TParams SearchParams] struct{}

func NewDefaultResponseBuilder[TResult SearchResult, TParams SearchParams]() ResponseBuilder[TResult, TParams] {
	return &DefaultResponseBuilder[TResult, TParams]{}
}

func (b *DefaultResponseBuilder[TResult, TParams]) BuildSearchResponse(results TResult, total int, params TParams) SearchResponse {
	hasMore := params.GetOffset()+len(results.GetItems()) < total

	return SearchResponse{
		Data: results.GetItems(),
		Pagination: PaginationDetails{
			Total:   total,
			Limit:   params.GetLimit(),
			Offset:  params.GetOffset(),
			HasMore: hasMore,
		},
	}
}

func (b *DefaultResponseBuilder[TResult, TParams]) BuildErrorResponse(err error) (int, ErrorResponse) {
	switch e := err.(type) {
	case *RequestParseError:
		// BAD REQUEST - client error
		return http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeInvalidRequest,
				Message: "Invalid request parameters",
				Details: []string{e.Error()},
			},
		}
	case *ValidationError:
		// BAD REQUEST - client error
		return http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeValidationError,
				Message: "Invalid search parameters",
				Details: e.Errors,
			},
		}
	case *SearchError:
		// INTERNAL SERVER ERROR - server/infrastructure error
		return http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeSearchError,
				Message: fmt.Sprintf("Failed to %s", e.Operation),
				Details: []string{e.Error()},
			},
		}
	case *ConversionError:
		// BAD REQUEST - client sent data that can't be converted
		return http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeValidationError,
				Message: "Invalid search parameters",
				Details: []string{e.Error()},
			},
		}
	default:
		// INTERNAL SERVER ERROR - unexpected error
		return http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrCodeInternalError,
				Message: "Internal server error",
				Details: []string{err.Error()},
			},
		}
	}
}
