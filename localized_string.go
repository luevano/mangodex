package mangodex

import (
	"encoding/json"
	"fmt"
)

// LocalisedStrings: A struct wrapping around a map containing each localised string.
type LocalisedStrings struct {
	Values map[string]string
}

func (l *LocalisedStrings) UnmarshalJSON(data []byte) error {
	l.Values = map[string]string{}

	// First try if can unmarshal directly.
	if err := json.Unmarshal(data, &l.Values); err == nil {
		return nil
	}

	// If fail, try to unmarshal to array of maps.
	var locals []map[string]string
	if err := json.Unmarshal(data, &locals); err != nil {
		return fmt.Errorf("error unmarshalling localisedstring: %s", err.Error())
	}

	// If pass, then add each item in the array to flatten to one map.
	for _, entry := range locals {
		for key, value := range entry {
			l.Values[key] = value
		}
	}
	return nil
}

// GetLocalString: Get the localised string for a particular language code.
// If the required string is not found, it will return the first entry, or an empty string otherwise.
func (l *LocalisedStrings) GetLocalString(langCode string) string {
	// If we cannot find the required code, then return first value.
	if s, ok := l.Values[langCode]; !ok {
		var v string
		for _, value := range l.Values {
			v = value
			break
		}
		return v
	} else {
		return s
	}
}
