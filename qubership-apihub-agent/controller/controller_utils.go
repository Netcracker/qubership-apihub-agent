package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Netcracker/qubership-apihub-agent/exception"
	"github.com/gorilla/mux"
)

func getStringParam(r *http.Request, p string) string {
	params := mux.Vars(r)
	return params[p]
}

func getUnescapedStringParam(r *http.Request, p string) (string, error) {
	params := mux.Vars(r)
	return url.PathUnescape(params[p])
}

func getFailOnErrorQueryParam(r *http.Request) (bool, *exception.CustomError) {
	if r.URL.Query().Get("failOnError") != "" {
		val, err := strconv.ParseBool(r.URL.Query().Get("failOnError"))
		if err != nil {
			return false, &exception.CustomError{
				Status:  http.StatusBadRequest,
				Code:    exception.IncorrectParamType,
				Message: exception.IncorrectParamTypeMsg,
				Params:  map[string]interface{}{"param": "failOnError", "type": "bool"},
				Debug:   err.Error(),
			}
		}
		return val, nil
	}
	return false, nil
}
