package transport

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// BuildMiddlewareLogger return the middleware logger
func BuildMiddlewareLogger() echo.MiddlewareFunc {

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:        true,
		LogURI:           true,
		LogError:         true,
		LogStatus:        true,
		LogLatency:       true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			c.Logger().Info(fmt.Sprintf("REQUEST %v [%v]%v  %v, request[%v] response[%v]", v.Status, v.Method, v.URI, v.Latency, v.ContentLength, v.ResponseSize))
			if v.Error != nil {
				c.Logger().Error(fmt.Sprintf("\n    %v\n\n", v.Error))
			}
			return nil
		},
	})
}
