package tesoql_echo

import (
	"github.com/labstack/echo/v4"
)

// GenericResponseModel represents a standard response structure for TesoQL queries.
// It includes the size of the data returned, the total count of matching records,
// and the actual data as a slice of maps.
type GenericResponseModel struct {
	Size       int                      `json:"size"`       // The number of records in the current response.
	TotalCount int                      `json:"totalCount"` // The total number of records that match the query.
	Data       []map[string]interface{} `json:"data"`       // The data returned by the query.
}

func newTesoQlResponseModel(totalCount int, size int, data []map[string]interface{}) *GenericResponseModel {
	model := GenericResponseModel{
		Size:       size,
		TotalCount: totalCount,
		Data:       data,
	}
	return &model
}

// HttpError represents a standardized error structure used in HTTP responses.
// It includes an error code, a key identifying the error, a descriptive message,
// and a specific TesoQL error code.
type HttpError struct {
	Code          int    `json:"code"`          // The HTTP status code for the error.
	Key           string `json:"error"`         // A key identifying the type of error.
	Message       string `json:"message"`       // A descriptive message explaining the error.
	TesoQlErrCode int    `json:"tesoQlErrCode"` // A specific error code related to TesoQL.
}

// NewHTTPError creates a new HttpError and sends it as a JSON response to the client.
// This function is typically used to handle errors encountered during request processing.
//
// Example usage:
//
//	// Handling an error during request processing
//	err := NewHTTPError(c, 400, "BINDING_ERROR", "Failed to bind request payload.", 400000)
//
//	// Using NewHTTPError as a custom error function in TesoQLEcho
//	tesoQLEcho := tesoql_echo.NewTesoQLEcho(cfg, tesoql_echo.NewHTTPError)
//	tesoQLEcho.Route(e, "/query")
//
// Returns:
//
// - error: An error object that can be returned by an Echo handler.
func NewHTTPError(c echo.Context, errCode int, key string, msg string, tesoQlCode int) error {
	err := &HttpError{
		Code:          errCode,
		Key:           key,
		Message:       msg,
		TesoQlErrCode: tesoQlCode,
	}
	return c.JSON(err.Code, err)
}

func (e *HttpError) Error() string {
	return e.Key + ": " + e.Message
}
