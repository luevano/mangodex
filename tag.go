package mangodex

import "github.com/google/uuid"

// Tag : Struct containing information on a tag.
type Tag struct {
	ID            uuid.UUID        `json:"id"`
	Type          RelationshipType `json:"type"`
	Attributes    TagAttributes    `json:"attributes"`
	Relationships []*Relationship  `json:"relationships"`
}

// TagAttributes : Attributes for a Tag.
type TagAttributes struct {
	Name        LocalisedStrings `json:"name"`
	Description LocalisedStrings `json:"description"`
	Group       string           `json:"group"`
	Version     int              `json:"version"`
}

// GetName : Get name of the tag.
func (t *Tag) GetName(langCode string) string {
	return t.Attributes.Name.GetLocalString(langCode)
}
