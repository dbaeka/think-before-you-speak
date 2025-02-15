package common

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrelationIDContext(t *testing.T) {
	w := httptest.NewRecorder()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	c := e.NewContext(req, w)

	assert.Contains(t, GetCorrelationID(c.Request().Context()), "gen_")

	newCtx := SetCorrelationID(c.Request().Context(), "test")
	assert.Equal(t, "test", GetCorrelationID(newCtx))
}

func TestEchoToGoContext(t *testing.T) {
	w := httptest.NewRecorder()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	c := e.NewContext(req, w)

	assert.Equal(t, c.Request().Context(), EchoToGoContext(c))
}
