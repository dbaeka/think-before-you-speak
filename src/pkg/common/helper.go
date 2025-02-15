package common

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"think-before-you-speak/pkg/log"
)

func WrapFunc(f interface{}, args ...interface{}) echo.HandlerFunc {
	fn := reflect.ValueOf(f)
	if fn.Type().NumIn() != len(args) {
		panic(fmt.Sprintf("invalid input parameters of function %v", fn.Type()))
	}

	outNum := fn.Type().NumOut()
	if outNum == 0 {
		panic(fmt.Sprintf("invalid output parameters of function %v, at least one, but got %d", fn.Type(), outNum))
	}

	inputs := make([]reflect.Value, len(args))
	for i, in := range args {
		inputs[i] = reflect.ValueOf(in)
	}

	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				l := log.Get()
				l.Warn().Msgf("panic: %v", err)
				ResponseFailed(c, http.StatusInternalServerError, fmt.Errorf("%v", err))
				return
			}
		}()

		outputs := fn.Call(inputs)
		if len(outputs) > 1 {
			err, ok := outputs[len(outputs)-1].Interface().(error)
			if ok && err != nil {
				return ResponseFailed(c, http.StatusInternalServerError, err)
			}
		}

		return c.JSON(http.StatusOK, outputs[0].Interface())
	}
}

func GetOrDefault[T any](data map[string]interface{}, key string, defaultValue T) T {
	if value, exists := data[key]; exists {
		result, ok := value.(T)
		if ok {
			return result
		}
	}
	return defaultValue
}

// FormatCost formats a cost in dollar amount.
//
// Parameters
//
//	-cost - The value to format
//
// Returns
// A string of the cost with 3 decimal places and $ prefixed
func FormatCost(cost float64) string {
	return fmt.Sprintf("$%.3f", cost)
}
