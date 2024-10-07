package coprocessor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RouterLifecycleRequest struct {
	Stage string `json:"stage"`
	Body  any    `json:"body"`
}

type CommonProperties struct {
	Version int    `json:"version"`
	Stage   string `json:"stage"`
	ID      string `json:"id,omitempty"`
}

type Headers struct {
	Headers http.Header `json:"headers,omitempty"`
}

// The value of "control" coming from the router is always a string value of "continue"
// If you don't modify it, or respond with it, and you don't respond with an object, everything proceeds normally
// The only time you are going to modify it, you're going to set it to an object with a Break property and specific status code
// https://www.apollographql.com/docs/router/customizations/coprocessor/#terminating-a-client-request
type BreakControl struct {
	Break float64 `json:"break,omitempty"`
}

type Context struct {
	Entries map[string]any `json:"entries"`
}

type Body struct {
	Errors        []Error        `json:"errors,omitempty"`
	Query         string         `json:"query,omitempty"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`

	// RouterResponse Stage
	Data any `json:"data,omitempty"`
}

type Error struct {
	Message    string     `json:"message,omitempty"`
	Extensions *Extension `json:"extensions,omitempty"`
}

type Extension struct {
	Code string `json:"code,omitempty"`

	// Adding defaults to ErrorType and ErrorCode means even if there are no errors
	// JSON.Marshal sees this as non-nil and will erroneously add them to all responses
	ErrorType   string `json:"errorType,omitempty"`
	ErrorCode   string `json:"errorCode,omitempty"`
	ServiceName string `json:"serviceName,omitempty"`
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var response []byte

	cr, err := HandleRequest(w, r)
	if err != nil {
		logger.Error(err, "Error handling coprocessor request")
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	response, err = json.Marshal(&cr)
	if err != nil {
		logger.Error(err, "Failed to marshal response")
		http.Error(w, fmt.Sprintf("Failed to marshal response: %s", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		logger.Error(err, "Error writing coprocessor response")
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
}

func NewRequest(r *http.Request) (*[]byte, string, error) {
	var err error
	var cr *RouterLifecycleRequest

	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, "", fmt.Errorf("error reading request body: %w", err)
	}

	// If the Router is configured to send data, this should never be empty
	// If it isn't configured to send data, it shouldn't call the coprocessor
	if len(httpRequestBody) == 0 {
		return nil, "", fmt.Errorf("error empty http request body at /%s", r.URL.Path[1:])
	}

	err = json.Unmarshal(httpRequestBody, &cr)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	return &httpRequestBody, cr.Stage, nil
}
