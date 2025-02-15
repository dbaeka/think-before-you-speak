package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lithammer/shortuuid/v3"
	"think-before-you-speak/pkg/common"
	"think-before-you-speak/pkg/config"
	"think-before-you-speak/pkg/log"
)

func useMiddlewares(e *echo.Echo, conf *config.Config) {
	e.Use(
		middleware.Recover(),
		middleware.CORS(),
		middleware.Secure(),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return shortuuid.New()
			},
		}),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:       true,
			LogRequestID: true,
			LogRemoteIP:  true,
			LogUserAgent: true,
			LogStatus:    true,
			LogMethod:    true,
			LogLatency:   true,
			LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
				log.FromContext(c.Request().Context()).Info().Fields(
					map[string]interface{}{
						"URI":        values.URI,
						"request_id": values.RequestID,
						"remote_ip":  values.RemoteIP,
						"user_agent": values.UserAgent,
						"status":     values.Status,
						"method":     values.Method,
						"duration":   values.Latency.String(),
					}).Err(values.Error).Msg("Request done")

				return nil
			},
		}),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				ctx := req.Context()

				reqCorrelationID := req.Header.Get(common.CorrelationIDHttpHeader)
				if reqCorrelationID == "" {
					reqCorrelationID = shortuuid.New()
				}

				ctx = log.ToContext(ctx, map[string]interface{}{common.CorrelationIDKey: reqCorrelationID})
				ctx = common.SetCorrelationID(ctx, reqCorrelationID)

				c.SetRequest(req.WithContext(ctx))
				c.Response().Header().Set(common.CorrelationIDHttpHeader, reqCorrelationID)

				return next(c)
			}
		},
	)
}
