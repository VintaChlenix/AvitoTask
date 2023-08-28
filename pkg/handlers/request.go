package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnmarshalJSON(httpRequest *http.Request, request interface{}) error {
	body, err := io.ReadAll(httpRequest.Body)
	if err != nil {
		return fmt.Errorf("failed to read http body: %w", err)
	}
	if len(body) == 0 {
		return fmt.Errorf("failed to read http body: body is nil")
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return nil
}
