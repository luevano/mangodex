package mangodex

import "github.com/google/uuid"

// Tag: Struct containing information on a tag.
type Tag struct {
	ID            uuid.UUID        `json:"id"`
	Type          RelationshipType `json:"type"`
	Attributes    TagAttributes    `json:"attributes"`
	Relationships []*Relationship  `json:"relationships"`
}

// TagAttributes: Attributes for a tag.
type TagAttributes struct {
	Name        LocalisedStrings `json:"name"`
	Description LocalisedStrings `json:"description"`
	Group       TagGroup         `json:"group"`
	Version     int              `json:"version"`
}

// GetName: Get name of the tag.
//
// If the requested language code tag name is not found and fallback is true,
// the first available value is returned, else an empty string.
func (t *Tag) GetName(langCode string, fallback bool) string {
	return t.Attributes.Name.GetLocalString(langCode, fallback)
}
