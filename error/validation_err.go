package cerror

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iamrz1/rutils"
)

var (
	InvalidInputErr = fmt.Errorf("input data is not valid")
	NotFoundErr     = fmt.Errorf("could not find data")
)

// ValidationError ...
type ValidationError map[string]string

func (err ValidationError) Error() string {
	buf, _ := json.Marshal(err)
	return string(buf)
}

// Add ...
func (err ValidationError) Add(key, msg string) {
	err[key] = msg
}

var validationErrorTexts = []string{
	"UNIQUE constraint failed",
}

func GetValidationErr(err error) error {
	found, _ := containsAny(err.Error())
	if !found {
		return err
	}

	return rutils.NewValidationError("", err)
}

func containsAny(input string) (bool, string) {
	for _, str := range validationErrorTexts {
		if strings.Contains(input, str) {
			return true, str
		}
	}

	return false, ""
}
