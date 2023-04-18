package model

import (
	errs "github.com/canbo-x/port-service/internal/error"
)

// ValidatePortID validates the port id.
// Just for demonstration purposes
// Please read the "Validation of Port ID" title
// under "Personal Thoughts and Notes" section in the README.md
func ValidatePortID(id string) error {
	if len(id) == 0 || len(id) > 8 {
		return errs.ErrInvalidPortID
	}

	return nil
}
