package tsserv

import (
	"fmt"
	"github.com/tinkermode/tsserv/pkg/datasource"
	"github.com/tinkermode/tsserv/pkg/logger"
	"net/http"
	"reflect"
	"time"
)

type RequestParams struct {
	Begin time.Time `form:"begin"`
	End   time.Time `form:"end"`
}

func parseRequestParams(request *http.Request, params interface{}) error {
	v := reflect.ValueOf(params).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("form")

		value := request.FormValue(tag)
		if value == "" {
			return fmt.Errorf("missing required parameter: %s", tag)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			// handle integer fields if needed
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				parsedTime, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return fmt.Errorf("Param '%s' must be in RFC3339 format", tag)
				}
				field.Set(reflect.ValueOf(parsedTime))
			}
		default:
			return fmt.Errorf("unhandled default case")
		}
	}
	return nil
}

func sayHello(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		sendErrorResponse(response, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	response.WriteHeader(http.StatusOK)
	if _, err := response.Write([]byte("Hello\n")); err != nil {
		logger.ErrorLogger.Printf("Failed to write response (%v)\n", err)
	} else {
		logger.InfoLogger.Printf("Responded to /hello")
	}
}

func getRawDataPoints(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		sendErrorResponse(response, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var params RequestParams
	if err := parseRequestParams(request, &params); err != nil {
		sendErrorResponse(response, http.StatusBadRequest, err.Error())
		return
	}

	ds := datasource.New()

	cur, err := ds.Query(params.Begin, params.End)
	if err != nil {
		sendErrorResponse(response, http.StatusForbidden, fmt.Sprintf("Failed to fetch data: %v", err))
		return
	}

	response.WriteHeader(http.StatusOK)

	for {
		dp, ok := cur.Next()
		if !ok {
			break
		}

		if _, err := response.Write([]byte(fmt.Sprintf("%s %8.4f\n", dp.Timestamp.Format(time.RFC3339), dp.Value))); err != nil {
			logger.ErrorLogger.Printf("Failed to write response (%v)\n", err)
			break
		}
	}
}
