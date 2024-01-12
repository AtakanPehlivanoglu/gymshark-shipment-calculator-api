package response

import (
	"fmt"
)

func ErrInvalidRequest(err error) []byte {
	return []byte(fmt.Sprintf("Invalid Request - %v", err))
}

func ErrInternalServer(err error) []byte {
	return []byte(fmt.Sprintf("Internal server error - %v", err))
}
