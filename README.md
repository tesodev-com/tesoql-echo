# tesoql-echo
*tesoql-echo* is a package for integrating *tesoql* with the *Echo framework*. It simplifies the process of handling database queries that supported by *tesoql* within Echo-based web applications, providing flexible query-building capabilities and error handling.
#### Features
- **Seamless Echo Integration:** Easily integrate tesoql's query building into your Echo routes.
- **Customizable Error Handling:** Override default error behavior with custom error functions.
- **Multi-Engine Query Support:** Leverage tesoql to manage SQL and MongoDB queries.
- **Pagination and Response Models:** Built-in support for pagination and custom response models.

------------

### Installation
To install *tesoql-echo*, use the following `go get` command:

`go get github.com/tesodev-com/tesoql-echo`

------------

### Usage
Here’s an example of how to use tesoql-echo in an Echo application:

#### Step 1: Import the required packages

```go
import (
	"github.com/labstack/echo/v4"
	"tesodev-generic-get/tesoql"
	tesoql_echo "tesodev-generic-get/tesoql-echo"
)
```

#### Step 2: Initialize tesoql-echo

```go
cfg := tesoql.Config{ /* configuration setup */ }
tesoQLEcho := tesoql_echo.NewTesoQLEcho(cfg)
```
See [github.com/tesodev-com/tesoql](https://github.com/tesodev-com/tesoql?tab=readme-ov-file#1-config-struct "https://github.com/tesodev-com/tesoql") README for detailed explanation of configuration at the corresponding part.

##### Step 2.1: Giving an optional custom error response function
To customize error handling in tesoql-echo, you can pass your own error function, if you don’t tesoql-echo automatically uses built-in *tesoql_echo.NewHTTPError:*

```go
func NewHTTPError(c echo.Context, errCode int, key string, msg string, tesoQlCode int) error {
   err := &HttpError{
      Code:          errCode,
      Key:           key,
      Message:       msg,
      TesoQlErrCode: tesoQlCode,
   }
   return c.JSON(err.Code, err)
}
```

Custom error function should have this format ; 
`func(c echo.Context, errCode int, key string, msg string, tesoQlCode int) error
`

##### Step 3.1: Set up a route for handling queries (Way 1)

```go
cfg := tesoql.Config{ /* configuration setup */ }
tesoQLEcho := tesoql_echo.NewTesoQLEcho(cfg)

e := echo.New()

tesoQLEcho.Route(e, "/tesoql")

e.Start(":8080")
```

`tesoQLEcho.Route` method directly registers the handler function of the *tesoql-echo* as a handler function for *echo*. It creates a POST endpoint with the given route. 

##### Step 3.2: Register tesoql-echo as a handler function (Way 2)

```go
cfg := tesoql.Config{ /* configuration setup */ }
tesoQLEcho := tesoql_echo.NewTesoQLEcho(cfg)

e := echo.New()

e.POST("/tesoql", tesoqlEcho.TesoQLHandler)

e.Start(":8080")
```

You can directly register the handler method exposed by the *tesoql-echo* as an *Echo framework handler method* by yourself.

------------

### Request Payload
An instance of the request payload in *JSON* format is given below. The payload is exclusively explained on the link, and the *JSON* equivalent is derived from the *JsonMap* object on the link;
[github.com/tesodev-com/tesoql](https://github.com/tesodev-com/tesoql?tab=readme-ov-file#building-a-jsonmap-payload "github.com/tesodev-com/tesoql")

```json
{
    "search": {
        "productName": ["piz", "slic"]
    },
    "projectionFields": [
        "id",
        "productName",
        "productId",
        "order_status",
        "orderId",
        "quantity",
        "amount",
        "userCount",
        "remainingStock",
        "order_date"
    ],
    "sortConditions": [
        {
            "field": "quantity",
            "sortCondition": "ASC"
        },
        {
            "field": "amount",
            "sortCondition": "DESC"
        }
    ],
    "conditions": {
        "remainingStock": {
            "greaterThan": 25
        },
        "order_status": {
            "valuesToExclude": ["inactive"]
        },
        "quantity": {
            "greaterThan": 2,
            "lowerThan": 10
        },
        "order_date": {
            "greaterOrEqual": "2023-01-01T00:00:00",
            "lowerOrEqual": "2024-09-10T00:00:00"
        },
        "userCount": {
            "valuesToExactMatch": [2]
        }
    },
    "pagination": {
        "limit": 48,
        "offset": 2
    },
    "totalCount": true,
    "suppressDataResponse": false
}

```

------------

### Response Payload 

```json
{
   "size": 0,
   "totalCount": 0,
   "data": [] 
}

```

------------

### Error Response Example

```json
{
   "code": 400,
   "error": "BINDING_ERROR",
   "message": "Error encountered while binding the request payload!",
   "tesoQlErrCode": 400000
}

```

------------


## Contributing
Contributions are welcome! Please open an issue or submit a pull request if you have any improvements or bug fixes.
