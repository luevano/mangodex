package mangodex

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
)

// LocalisedStrings: A struct wrapping around a map containing each localised string.
type LocalisedStrings struct {
	// Values: First value for each language.
	Values map[string]string
	// List: All values in the order the API sent them (alt titles can have many per language).
	List []LocalisedString
}

// LocalisedString: A single value and its language code.
type LocalisedString struct {
	Lang  string
	Value string
}

func (l *LocalisedStrings) UnmarshalJSON(data []byte) error {
	l.Values = map[string]string{}
	l.List = nil

	// Either a single object or an array of objects (e.g. alt titles).
	var locals []map[string]string
	var local map[string]string
	if err := json.Unmarshal(data, &local); err == nil {
		locals = []map[string]string{local}
	} else if err := json.Unmarshal(data, &locals); err != nil {
		return fmt.Errorf("decode localised strings: %w", err)
	}

	for _, entry := range locals {
		// a JSON object has no order, sort its languages so it's always the same
		langs := slices.Sorted(maps.Keys(entry))
		for _, lang := range langs {
			value := entry[lang]
			l.List = append(l.List, LocalisedString{lang, value})
			// the API lists the official/most common value first
			if _, ok := l.Values[lang]; !ok {
				l.Values[lang] = value
			}
		}
	}
	return nil
}

// GetAll: Get all the localised strings for a particular language code, in order.
func (l *LocalisedStrings) GetAll(langCode string) []string {
	var values []string
	for _, s := range l.List {
		if s.Lang == langCode {
			values = append(values, s.Value)
		}
	}
	return values
}

// GetLocalString: Get the localised string for a particular language code.
//
// If the required string is not found and fallback is true, it will return the first entry,
// or an empty string otherwise.
func (l *LocalisedStrings) GetLocalString(langCode string, fallback bool) string {
	s, found := l.Values[langCode]
	if found {
		return s
	}

	if fallback && len(l.List) != 0 {
		return l.List[0].Value
	}

	return ""
}
