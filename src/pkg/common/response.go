package common

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"net/http"
	"think-before-you-speak/pkg/log"
)

func NewResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, data)
}

func ResponseSuccess(c echo.Context, data interface{}) error {
	return NewResponse(c, http.StatusOK, data)
}

func ResponseSuccessStream(c echo.Context, server *sse.Server) error {
	if server == nil {
		return fmt.Errorf("server is nil")
	}
	l := log.Get()
	go func() {
		<-EchoToGoContext(c).Done() // Received Client Side Disconnection eg. Browser Disconnection
		l.Info().Msgf("The client is disconnected: %v", c.RealIP())
		return
	}()
	server.ServeHTTP(c.Response(), c.Request())
	return nil
}

func ResponseFailed(c echo.Context, code int, err error) error {
	if code == 0 {
		code = http.StatusInternalServerError
	}

	var msg string
	if err != nil {
		msg = err.Error()
		var url string
		if c.Request() != nil {
			url = c.Request().URL.String()
		}
		l := log.Get()
		l.Warn().Msgf("url: %s, error: %v", url, msg)
	}
	return NewResponse(c, code, struct {
		Msg string `json:"msg"`
	}{
		Msg: msg,
	})
}
