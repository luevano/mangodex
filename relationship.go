package mangodex

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Relationship: Struct containing relationships, with optional attributes for the relation.
type Relationship struct {
	ID         uuid.UUID        `json:"id"`
	Type       RelationshipType `json:"type"`
	Attributes interface{}      `json:"attributes"`
}

func (a *Relationship) UnmarshalJSON(data []byte) error {
	// Check for the type of the relationship, then unmarshal accordingly.
	typ := struct {
		ID         uuid.UUID        `json:"id"`
		Type       RelationshipType `json:"type"`
		Attributes json.RawMessage  `json:"attributes"`
	}{}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	var err error
	switch typ.Type {
	case RelationshipTypeManga:
		a.Attributes = &MangaAttributes{}
	case RelationshipTypeAuthor:
		a.Attributes = &AuthorAttributes{}
	case RelationshipTypeScanlationGroup:
		a.Attributes = &ScanlationGroupAttributes{}
	case RelationshipTypeCoverArt:
		a.Attributes = &CoverAttributes{}
	default:
		a.Attributes = &json.RawMessage{}
	}

	a.ID = typ.ID
	a.Type = typ.Type
	if typ.Attributes != nil {
		if err = json.Unmarshal(typ.Attributes, a.Attributes); err != nil {
			return fmt.Errorf("error unmarshalling relationship of type %s: %s, %s",
				typ.Type, err.Error(), string(data))
		}
	}
	return err
}
