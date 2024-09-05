package tesoql_echo

import (
	"github.com/labstack/echo/v4"
	"github.com/tesodev-com/tesoql"
	"net/http"
)

// TesoQLEcho provides an interface to integrate TesoQL with the Echo framework.
// It encapsulates a TesoQL instance, field mappings, pagination configuration,
// and a custom error handling function.
type TesoQLEcho struct {
	tesoQL          *tesoql.TesoQL
	cfg             *tesoql.Config
	customErrorFunc func(c echo.Context, errCode int, key string, msg string, tesoQlCode int) error
}

// NewTesoQLEcho initializes a new TesoQLEcho instance with the provided TesoQL configuration.
// Optionally, a custom error handling function can be provided to override the default error behavior.
//
// Example usage:
//
//	cfg := tesoql.Config{ /* TesoQL configuration setup */ }
//	tesoQLEcho := tesoql_echo.NewTesoQLEcho(cfg, customErrorFunc)
//
// Returns:
//
// - *TesoQLEcho: A pointer to the initialized TesoQLEcho struct.
func NewTesoQLEcho(
	cfg tesoql.Config,
	customErrorFn ...func(c echo.Context, errCode int, key string, msg string, tesoQlCode int) error) *TesoQLEcho {

	tesoQL := cfg.NewTesoQL()

	var customErrorFunc func(c echo.Context, errCode int, key string, msg string, tesoQlCode int) error

	if len(customErrorFn) > 0 {
		customErrorFunc = customErrorFn[0]
	} else {
		customErrorFunc = NewHTTPError
	}

	return &TesoQLEcho{
		tesoQL:          tesoQL,
		cfg:             &cfg,
		customErrorFunc: customErrorFunc,
	}
}

// Route sets up a POST route in the Echo framework to handle TesoQL queries.
// The specified path will be mapped to the TesoQLHandler function.
//
// Example usage:
//
//	e := echo.New()
//	tesoQLEcho := tesoql_echo.NewTesoQLEcho(cfg)
//	tesoQLEcho.Route(e, "/tesoql")
//
// Returns:
//
// - none
func (te *TesoQLEcho) Route(e *echo.Echo, path string) {
	e.POST(path, te.TesoQLHandler)
}

// TesoQLHandler is the main HTTP handler function for processing TesoQL queries.
// It binds the incoming request to a JsonMap, validates it, and then processes
// the query using the TesoQL service. The results are returned as a JSON response.
//
// Example usage:
//
//	// Route setup
//	e.POST("/tesoql", tesoQLEcho.TesoQLHandler)
//
// Returns:
//
// - error: An HTTP error if binding, validation, or query processing fails.
func (te *TesoQLEcho) TesoQLHandler(c echo.Context) error {
	jsonMap := tesoql.JsonMap{}
	err := c.Bind(&jsonMap)
	if err != nil {
		return te.customErrorFunc(
			c,
			http.StatusBadRequest,
			tesoql.BINDING_ERR,
			"Error encountered while binding the request payload!",
			tesoql.BINDING_ERR_CODE)
	}

	validationErr := jsonMap.Validate(te.cfg)
	if validationErr != nil {
		return te.customErrorFunc(
			c,
			validationErr.ErrorCode/1000,
			validationErr.ErrorType,
			validationErr.ErrorMsg,
			validationErr.ErrorCode)
	}

	service := *te.tesoQL.Service
	data, totalCount, length, serviceErr := service.Get(&jsonMap)
	if serviceErr != nil {
		return te.customErrorFunc(
			c,
			serviceErr.ErrorCode/1000,
			serviceErr.ErrorType,
			serviceErr.ErrorMsg,
			serviceErr.ErrorCode)
	}

	resp := newTesoQlResponseModel(totalCount, length, data)
	return c.JSON(http.StatusOK, resp)
}
