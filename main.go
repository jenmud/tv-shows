package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort int = 8000

// HTTPJSONParseError is a error which is raise if any of the
// request handlers have issues parsing JSON formatted data.
type HTTPJSONParseError struct {
	Error string `json:"error"`
}

var logger = log.New(os.Stdout, "", log.LstdFlags)

// getIntFromEnvvar gets and converts the environment variable into a integer.
func getIntFromEnvvar(key string) (int, error) {
	value, ok := os.LookupEnv(key)

	if ok {
		port, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return port, nil
	}

	return 0, fmt.Errorf("Could not find environment variable %q", key)
}

// GetListenPort returns a port number which can be used for the server as
// a listening port to accept client connections.
// It will first check and use the environment variable `PORT` if
// found else it will return the default number 8000.
func GetListenPort() int {
	port, err := getIntFromEnvvar("PORT")
	if err != nil {
		log.Printf("%s, using default port %d\n", err, defaultPort)
		return defaultPort
	}
	return port
}

// TvShowJSONHandler is a handler which parses a JSON formatted
// request body and filters for tv shows which has DRM enabled and
// has one or more episodes.
func TvShowJSONHandler(w http.ResponseWriter, req *http.Request) {
	httpError, _ := json.Marshal(
		HTTPJSONParseError{
			Error: "Could not decode request: JSON parsing failed",
		},
	)

	// Because we are only dealing with JSON data, lets set the
	// connect headers.
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(
			w,
			string(httpError),
			http.StatusBadRequest,
		)
	} else {
		shows, err := ParseRequestJSON(body)

		if err != nil {
			logger.Printf("Original Body: %s", body)
			logger.Printf("Request JSON parsing error: %s", err)

			http.Error(
				w,
				string(httpError),
				http.StatusBadRequest,
			)

		} else {

			shows = FilterTVShowsForDRM(shows)
			shows = FilterTVShowsWithEpisodes(shows)

			response, err := MakeResponseJSON(shows)
			if err != nil {
				logger.Printf("Original Body: %s", body)
				logger.Printf(
					"Request JSON parsing error: %s",
					err,
				)

				http.Error(
					w,
					string(httpError),
					http.StatusBadRequest,
				)

			} else {
				w.Write(response)
			}
		}
	}
}

// RunServer starts the server and accepts connections.
func RunServer(port int) error {
	http.HandleFunc("/", TvShowJSONHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// main is the entry point.
func main() {
	port := GetListenPort()
	logger.Printf("Listening on %d\n", port)
	logger.Fatal(RunServer(port))
}
