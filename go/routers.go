/*
 * ExCommerce
 *
 * ExCommerce is an example commerce system.
 *
 * API version: beta
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// A Route defines the parameters for an api endpoint
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes are a collection of defined api endpoints
type Routes []Route

// Router defines the required methods for retrieving api routes
type Router interface {
	Routes() Routes
}

// NewRouter creates a new router for any number of api routers
func NewRouter(routers ...Router) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, api := range routers {
		for _, route := range api.Routes() {
			var handler http.Handler
			handler = route.HandlerFunc
			handler = logger(handler, route.Name)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}

	return router
}

func invalidInput(message, details string, w http.ResponseWriter) {
	status := http.StatusBadRequest // 400
	err := EncodeJSONResponse(map[string]string{
		"message": message,
		"details": details,
	}, &status, w)
	if err != nil {
		panic(err)
	}
}

func invalidJSON(jsonErr error, w http.ResponseWriter) {
	invalidInput("Invalid JSON in request body.", jsonErr.Error(), w)
}

func failValidation(message, pointer string, w http.ResponseWriter) {
	status := http.StatusUnprocessableEntity // 422
	err := EncodeJSONResponse(MalformedInputError{
		Message: message,
		Pointer: pointer,
	}, &status, w)
	if err != nil {
		panic(err)
	}
}

func unexpectedError(_ error, w http.ResponseWriter) {
	status := http.StatusInternalServerError // 500
	err := EncodeJSONResponse(map[string]string{
		"message": "Unexpected error.",
	}, &status, w)
	if err != nil {
		panic(err)
	}
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code
func EncodeJSONResponse(i interface{}, status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return json.NewEncoder(w).Encode(i)
}

// ReadFormFileToTempFile reads file data from a request form and writes it to a temporary file
func ReadFormFileToTempFile(r *http.Request, key string) (*os.File, error) {
	r.ParseForm()
	formFile, _, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}

	defer formFile.Close()
	file, err := ioutil.TempFile("tmp", key)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(formFile)
	if err != nil {
		return nil, err
	}

	file.Write(fileBytes)
	return file, nil
}

// parseIntParameter parses a sting parameter to an int64
func parseIntParameter(param string) (int64, error) {
	return strconv.ParseInt(param, 10, 64)
}
