package common

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWrapFunc(t *testing.T) {
	testCases := []struct {
		name             string
		function         interface{}
		inputs           []interface{}
		expectedPanic    bool
		expectedRespCode int
		expectedRespBody string
	}{
		{"invalid input parameters", func(int) string { return "" }, []interface{}{1, 2}, true, 0, ""},
		{"invalid output parameters", func() {}, nil, true, 0, ""},
		{"error response", func() (string, error) { return "", assert.AnError }, nil, false, 500, fmt.Sprintf(`{"msg":"%s"}`, assert.AnError)},
		{"function panic", func() (string, error) { panic("some error"); return "", assert.AnError }, nil, false, 500, fmt.Sprintf(`{"msg":"some error"}`)},
		{"success response with one outputs", func() string { return "some msg" }, nil, false, 200, `"some msg"`},
		{"success response with two outputs", func() (string, error) { return "some msg", nil }, nil, false, 200, `"some msg"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedPanic {
				defer func() {
					err := recover()
					assert.NotNil(t, err)
				}()
			}

			e := echo.New()
			w := httptest.NewRecorder()

			h := WrapFunc(tc.function, tc.inputs...)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			c := e.NewContext(req, w)
			err := h(c)
			if err != nil {
				return
			}

			assert.Equal(t, tc.expectedRespCode, w.Code)
			assert.Equal(t, tc.expectedRespBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		data         map[string]interface{}
		key          string
		defaultValue interface{}
		expected     interface{}
	}{
		{"key exists", map[string]interface{}{"key": "value"}, "key", "default", "value"},
		{"key does not exist", map[string]interface{}{"key": "value"}, "key2", "default", "default"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, GetOrDefault(tc.data, tc.key, tc.defaultValue))
		})
	}
}

func TestFormatCost(t *testing.T) {
	testCases := []struct {
		name     string
		cost     float64
		expected string
	}{
		{"zero cost", 0, "$0.000"},
		{"positive cost", 1.234, "$1.234"},
		{"negative cost", -1.234, "$-1.234"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, FormatCost(tc.cost))
		})
	}
}
