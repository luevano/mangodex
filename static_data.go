package mangodex

const APIVersion = "5.10.0"

// Publication demographic

type Demographic string

const (
	DemographicShounen Demographic = "shounen"
	DemographicShoujo  Demographic = "shoujo"
	DemographicJosei   Demographic = "josei"
	DemographicSeinen  Demographic = "seinen"
)

// Manga publication status

type PublicationStatus string

const (
	PublicationStatusOngoing   PublicationStatus = "ongoing"
	PublicationStatusCompleted PublicationStatus = "completed"
	PublicationStatusHiatus    PublicationStatus = "hiatus"
	PublicationStatusCancelled PublicationStatus = "cancelled"
)

// Manga reading status

type ReadingStatus string

const (
	ReadingStatusReading    ReadingStatus = "reading"
	ReadingStatusOnHold     ReadingStatus = "on_hold"
	ReadingStatusPlanToRead ReadingStatus = "plan_to_read"
	ReadingStatusDropped    ReadingStatus = "dropped"
	ReadingStatusReReading  ReadingStatus = "re_reading"
	ReadingStatusCompleted  ReadingStatus = "completed"
)

// Manga content rating

type ContentRating string

const (
	ContentRatingSafe       ContentRating = "safe"
	ContentRatingSuggestive ContentRating = "suggestive"
	ContentRatingErotica    ContentRating = "erotica"
	ContentRatingPorn       ContentRating = "pornographic"
)

// Relationship types. Useful for reference expansions

type RelationshipType string

const (
	RelationshipTypeManga           RelationshipType = "manga"
	RelationshipTypeChapter         RelationshipType = "chapter"
	RelationshipTypeCoverArt        RelationshipType = "cover_art"
	RelationshipTypeAuthor          RelationshipType = "author"
	RelationshipTypeArtist          RelationshipType = "artist"
	RelationshipTypeScanlationGroup RelationshipType = "scanlation_group"
	RelationshipTypeTag             RelationshipType = "tag"
	RelationshipTypeUser            RelationshipType = "user"
	RelationshipTypeCustomList      RelationshipType = "custom_list"
)

type MangaRelation string

const (
	// MangaRelationMonochrome : A Monochrome variant of this manga
	MangaRelationMonochrome MangaRelation = "monochrome"
	// MangaRelationColored : A Colored variant of this manga
	MangaRelationColored MangaRelation = "colored"
	// MangaRelationPreserialization : The original version of this manga before its official serialization
	MangaRelationPreserialization MangaRelation = "preserialization"
	// MangaRelationSerialization : The official serialization of this manga
	MangaRelationSerialization MangaRelation = "serialization"
	// MangaRelationPrequel : The previous entry in the same series
	MangaRelationPrequel MangaRelation = "prequel"
	// MangaRelationSequel : The next entry in the same series
	MangaRelationSequel MangaRelation = "sequel"
	// MangaRelationMainStory : The original narrative this manga is based on
	MangaRelationMainStory MangaRelation = "main_story"
	// MangaRelationSideStory : A side work contemporaneous with the narrative of this manga
	MangaRelationSideStory MangaRelation = "side_story"
	// MangaRelationAdaptedFrom : The original work this spin-off manga has been adapted from
	MangaRelationAdaptedFrom MangaRelation = "adapted_from"
	// MangaRelationSpinOff : An official derivative work based on this manga
	MangaRelationSpinOff MangaRelation = "spin_off"
	// MangaRelationBasedOn : The original work this self-published derivative manga is based on
	MangaRelationBasedOn MangaRelation = "based_on"
	// MangaRelationDoujinshi : A self-published derivative work based on this manga
	MangaRelationDoujinshi MangaRelation = "doujinshi"
	// MangaRelationSameFranchise : A manga based on the same intellectual property as this manga
	MangaRelationSameFranchise MangaRelation = "same_franchise"
	// MangaRelationSharedUniverse A manga taking place in the same fictional world as this manga
	MangaRelationSharedUniverse MangaRelation = "shared_universe"
	// MangaRelationAlternateStory : An alternative take of the story in this manga
	MangaRelationAlternateStory MangaRelation = "alternate_story"
	// MangaRelationAlternateVersion : A different version of this manga with no other specific distinction
	MangaRelationAlternateVersion MangaRelation = "alternate_version"
)

type OrderEnum string

const (
	OrderAscending  = "asc"
	OrderDescending = "desc"
)

type GetOrder struct {
	Name          OrderEnum `json:"name,omitempty"`
	CreatedAt     OrderEnum `json:"created_at,omitempty"`
	UpdatedAt     OrderEnum `json:"updated_at,omitempty"`
	FollowedCount OrderEnum `json:"followed_count,omitempty"`
	Relevance     OrderEnum `json:"relevance,omitempty"`
}
