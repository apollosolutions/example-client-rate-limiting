package coprocessor

import (
	"fmt"
	"net/http"
)

// Examples can be found:
// https://www.apollographql.com/docs/router/customizations/coprocessor/#coprocessor-request-format
func HandleRequest(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	httpRequestBody, stage, err := NewRequest(r)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling httpRequestBody: %w", err)
	}

	// Add more stages here as needed as per the docs: https://www.apollographql.com/docs/router/customizations/coprocessor/#stage
	switch stage {
	case "RouterRequest":
		return handleRouterRequest(httpRequestBody)
	case "RouterResponse":
		return handleRouterResponse(httpRequestBody)
	case "":
		// This shouldn't happen, everything should have a Stage
		return nil, fmt.Errorf("no stage for request %+v", httpRequestBody)
	default:
		return nil, fmt.Errorf("unhandled coprocessor request stage of type: %T", stage)
	}
}
