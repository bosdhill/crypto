package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"text/template"
)

// ValidateTemplate validates a text template results in valid JSON
// when it's executed with empty template data. If template execution
// results in invalid JSON, the template is invalid. When the template
// is valid, it can be used safely. A valid template can still result
// in invalid JSON when non-empty template data is provided.
func ValidateTemplate(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// get the default supported functions
	var failMessage string
	funcMap := GetFuncMap(&failMessage)

	// prepare the template with our template functions
	_, err := template.New("template").Funcs(funcMap).Parse(string(data))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	return nil
}

// ValidateTemplateData validates that template data is
// valid JSON.
func ValidateTemplateData(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if ok := json.Valid(data); !ok {
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("invalid JSON: %w", enrichJSONError(err))
		}

		// json.Valid() returns NOK, but decoding doesn't result in error with trailing brace.
		// It results in `map[subject:<nil>]`, instead. The Valid() function checks the entire JSON;
		// Decode() does not and sees the trailing brace as the final closing one, and thus stops
		// decoding.
		return errors.New("invalid JSON: early decoder termination")
	}

	return nil
}

// enrichJSONError tries to extract more information about the cause of
// an error related to a malformed JSON template and adds this to the
// error message.
func enrichJSONError(err error) error {
	var (
		syntaxError *json.SyntaxError
	)
	// TODO(hs): extracting additional info doesn't always work as expected, as the provided template is
	// first transformed by executing it. After transformation, the offsets in the error are not the offsets
	// for the original, user-provided template. If we want this to work, we should revert the transformation
	// somehow and then find the correct offset to use. This doesn't seem trivial to do.
	switch {
	case errors.As(err, &syntaxError):
		//return fmt.Errorf("%s at offset %d", err.Error(), syntaxError.Offset)
		return err
	default:
		return err
	}
}
