package tsserv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tinkermode/tsserv/pkg/datasource"
)

func sayHello(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		sendErrorResponse(response, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	response.WriteHeader(http.StatusOK)
	if _, err := response.Write([]byte("Hello\n")); err != nil {
		errLogger.Printf("Failed to write response (%v)\n", err)
	}
}

func getRawDataPoints(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		sendErrorResponse(response, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	paramBegin := request.FormValue("begin")
	begin, err := time.Parse(time.RFC3339, paramBegin)
	if err != nil {
		sendErrorResponse(response, http.StatusBadRequest, "Param 'begin' must be in RFC3339 format")
		return
	}

	paramEnd := request.FormValue("end")
	end, err := time.Parse(time.RFC3339, paramEnd)
	if err != nil {
		sendErrorResponse(response, http.StatusBadRequest, "Param 'end' must be in RFC3339 format")
		return
	}

	ds := datasource.New()

	cur, err := ds.Query(begin, end)
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
			errLogger.Printf("Failed to write response (%v)\n", err)
			break
		}
	}
}
