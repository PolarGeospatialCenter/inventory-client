package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PolarGeospatialCenter/inventory/pkg/lambdautils"
	"gopkg.in/resty.v1"
)

var (
	ErrConflict = errors.New("conflict")
)

func UnmarshalApiResponse(r *resty.Response, obj interface{}) error {
	success := r.StatusCode() >= 200 && r.StatusCode() < 300

	if success && obj != nil {
		err := json.Unmarshal(r.Body(), obj)
		if err != nil {
			return fmt.Errorf("unable to unmarshal response: %v", err)
		}
	} else if !success {
		errorResponse := &lambdautils.ErrorResponse{}
		err := json.Unmarshal(r.Body(), errorResponse)
		if err != nil {
			return fmt.Errorf("unable to unmarshal error response: %v", err)
		}
		return fmt.Errorf("%s: %s", errorResponse.Status, errorResponse.ErrorMessage)
	}
	return nil
}
