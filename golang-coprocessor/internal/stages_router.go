package coprocessor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type RouterStageBody struct {
	Body json.RawMessage `json:"body,omitempty"`
}

type RouterRequest struct {
	// Control properties
	*CommonProperties
	Control any `json:"control,omitempty"`

	// Data properties
	*Headers
	// The dash here tells Go's JSON library "don't do anything with this field"
	// We manually set it in the custom UnmarshalJSON method below
	Body    *Body    `json:"-"`
	Context *Context `json:"context,omitempty"`
	SDL     string   `json:"sdl,omitempty"`
	Path    string   `json:"path,omitempty"`
	Method  string   `json:"method,omitempty"`
}

// UnmarshalJSON implements a custom unmarshal function for the RouterRequest.  This is necessary
// because the Body field, represented by the Body struct, contains a struct with nested query
// content that must be scrubbed of extra escape characters.
// This is used any time we call `json.Unmarshal` into a RouterRequest
func (rr *RouterRequest) UnmarshalJSON(data []byte) error {

	// Use a custom type so we don't trigger an infinite Unmarshalling loop
	type TmpType RouterRequest

	var body *Body
	var err error
	var strBody string
	var tmpRequest TmpType

	// 1. Unmarshal everything but Body
	if err = json.Unmarshal(data, &tmpRequest); err != nil {

		return fmt.Errorf("failed to unmarshal router request: %v", err)
	}

	// 2. Unmarshal Body as a raw json string
	rrb := &RouterStageBody{}

	// We only log an error here, as an invalid body will result in an Unquote error below.
	// Additionally, I'm not sure there's a way to cause an error in unit tests when decoding to a
	// json.RawMessage value.
	if err = json.Unmarshal(data, &rrb); err != nil {
		logger.Error(err, "Failed to unmarshal body")
	}

	if len(rrb.Body) > 0 {

		// 3. Unquote/unescape body
		if strBody, err = strconv.Unquote(string(rrb.Body)); err != nil {

			return fmt.Errorf("failed to unquote: %v", err)
		}

		// There's some inconsistency between when the router sends a totally empty string and an escaped one ("\"\"")
		// For safety's sake, we check both here - if the body is empty, no need to create it
		if strBody != `""` && strBody != "" {
			// 4. Unmarshal into Body struct
			if err = json.Unmarshal([]byte(strBody), &body); err != nil {

				return fmt.Errorf("failed to unmarshal coprocessor request string into struct: %v", err)
			}

			tmpRequest.Body = body
		}
	}

	// 5. Create a new RouterRequest with the appropriately destructured fields
	*rr = RouterRequest(tmpRequest)

	return nil
}

func (rr *RouterRequest) MarshalJSON() ([]byte, error) {
	var err error

	type TmpType RouterRequest
	type RouterRequestBodyOverride struct {
		Body string `json:"body,omitempty"`
		*TmpType
	}

	routerRequestBodyOverride := &RouterRequestBodyOverride{
		TmpType: (*TmpType)(rr),
	}

	if rr.Body != nil {
		var jsonBody []byte
		jsonBody, err = json.Marshal(rr.Body)

		if err != nil {
			logger.Error(err, "Failed to marshal coprocessor request struct: %v")
		}

		routerRequestBodyOverride.Body = string(jsonBody)
	}

	return json.Marshal(routerRequestBodyOverride)
}

type RouterResponse struct {
	// Control properties
	*CommonProperties
	Control any `json:"control,omitempty"`

	// Data properties
	*Headers
	// The dash here tells Go's JSON library "don't do anything with this field"
	// We manually set it in the custom UnmarshalJSON method below
	Body       *Body    `json:"-"`
	Context    *Context `json:"context,omitempty"`
	StatusCode float64  `json:"statusCode,omitempty"`
	SDL        string   `json:"sdl,omitempty"`
}

// This is used any time we call `json.Unmarshal` into a RouterRequest
func (rr *RouterResponse) UnmarshalJSON(data []byte) error {

	// Use a custom type so we don't trigger an infinite Unmarshalling loop
	type TmpType RouterResponse

	var body *Body
	var err error
	var strBody string
	var tmpResponse TmpType

	// 1. Unmarshal everything but Body
	if err = json.Unmarshal(data, &tmpResponse); err != nil {

		return fmt.Errorf("failed to unmarshal router response: %v", err)
	}

	// 2. Unmarshal Body as a raw json string
	rrb := &RouterStageBody{}

	// We only log an error here, as an invalid body will result in an Unquote error below.
	// Additionally, I'm not sure there's a way to cause an error in unit tests when decoding to a
	// json.RawMessage value.
	if err = json.Unmarshal(data, &rrb); err != nil {
		logger.Error(err, "Failed to unmarshal body")
	}

	if len(rrb.Body) > 0 {

		// 3. Unquote/unescape body
		if strBody, err = strconv.Unquote(string(rrb.Body)); err != nil {

			return fmt.Errorf("failed to unquote: %v", err)
		}

		// There's some inconsistency between when the router sends a totally empty string and an escaped one ("\"\"")
		// For safety's sake, we check both here - if the body is empty, no need to create it
		if strBody != `""` && strBody != "" {
			// 4. Unmarshal into Body struct
			if err = json.Unmarshal([]byte(strBody), &body); err != nil {

				return fmt.Errorf("failed to unmarshal coprocessor request string into struct: %v", err)
			}

			tmpResponse.Body = body
		}

	}

	// 5. Create a new RouterRequest with the appropriately destructured fields
	*rr = RouterResponse(tmpResponse)

	return nil
}

func (rr *RouterResponse) MarshalJSON() ([]byte, error) {
	var err error

	type TmpType RouterResponse
	type RouterRequestBodyOverride struct {
		Body string `json:"body,omitempty"`
		*TmpType
	}

	routerRequestBodyOverride := &RouterRequestBodyOverride{
		TmpType: (*TmpType)(rr),
	}

	if rr.Body != nil {
		var jsonBody []byte
		jsonBody, err = json.Marshal(rr.Body)

		if err != nil {
			logger.Error(err, "Failed to marshal coprocessor request struct: %v")
		}

		routerRequestBodyOverride.Body = string(jsonBody)
	}

	return json.Marshal(routerRequestBodyOverride)
}

func handleRouterRequest(httpRequestBody *[]byte) (*RouterRequest, error) {
	cr, err := NewRouterRequest(httpRequestBody)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling httpRequestBody: %w", err)
	}

	// Normally you won't Marshal like this, this is simply to get a stringified
	// representation of the request in the console
	requestBody, err := json.Marshal(cr)

	if err != nil {
		logger.Error(err, "Failed to marshal coprocessor request struct: %v")
	}

	// This is the object sent by the Router that you can act upon to update headers, context, auth claims, etc
	// If you update the "control" property from "Continue" to something like { "break": 400 }, it will terminate the request and return the specified HTTP error
	// See: https://www.apollographql.com/docs/router/customizations/coprocessor/
	fmt.Println(string(requestBody))

	return cr, nil
}

func handleRouterResponse(httpRequestBody *[]byte) (*RouterResponse, error) {
	cr, err := NewRouterResponse(httpRequestBody)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling httpRequestBody: %w", err)
	}
	return cr, nil
}

func NewRouterRequest(httpRequestBody *[]byte) (*RouterRequest, error) {
	var err error
	var cr *RouterRequest
	err = json.Unmarshal(*httpRequestBody, &cr)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func NewRouterResponse(httpRequestBody *[]byte) (*RouterResponse, error) {
	var err error
	var cr *RouterResponse
	err = json.Unmarshal(*httpRequestBody, &cr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cr, nil
}
