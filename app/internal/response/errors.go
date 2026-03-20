package response

import (
	"net/http"
)

type ErrorCode string

type ApplicationError struct {
	ErrorCode    ErrorCode
	ErrorMessage string
	HttpCode     int
}

func (e *ApplicationError) Error() string {
	return e.ErrorMessage
}

const (
	SomethingWentWrong ErrorCode = "SWW01"
	InvalidParams      ErrorCode = "INV01"
	FetchingTeams      ErrorCode = "FT01"
	FetchingTeamRoles  ErrorCode = "FTR01"
)

var (
	ErrSomethingWentWrong = &ApplicationError{
		ErrorCode:    SomethingWentWrong,
		ErrorMessage: "something went wrong",
		HttpCode:     http.StatusInternalServerError,
	}
	ErrInvalidParams = &ApplicationError{
		ErrorCode:    InvalidParams,
		ErrorMessage: "invalid params",
		HttpCode:     http.StatusBadRequest,
	}
	ErrFetchingTeams = &ApplicationError{
		ErrorCode:    FetchingTeams,
		ErrorMessage: "failed to fetch teams",
		HttpCode:     http.StatusInternalServerError,
	}
	ErrFetchingTeamRoles = &ApplicationError{
		ErrorCode:    FetchingTeamRoles,
		ErrorMessage: "failed to fetch team roles",
		HttpCode:     http.StatusInternalServerError,
	}
)
