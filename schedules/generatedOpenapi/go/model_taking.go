// Code generatedOpenapi by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Schedules API
 *
 * API для работы с расписаниями
 *
 * API version: 1.0.0
 */

package openapi

type Taking struct {
	Name string `json:"name"`

	Time string `json:"time"`
}

// AssertTakingRequired checks if the required fields are not zero-ed
func AssertTakingRequired(obj Taking) error {
	elements := map[string]interface{}{
		"name": obj.Name,
		"time": obj.Time,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertTakingConstraints checks if the values respects the defined constraints
func AssertTakingConstraints(obj Taking) error {
	return nil
}
