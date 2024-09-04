package domain

// Page data structure
type QuizPageData struct {
	Code        string
	Title       string
	Description string
	Questions   []QuizQuestion
}

// Page data storage
type PageData struct {
	Title       string
	TableHeader string
	Columns     []string
	Rows        []QuizRowData
}

// Row description
type QuizRowData struct {
	Type        string `db:"type" json:"type"`               // Json datatype tag
	ID          string `db:"id" json:"id"`                   // Unique quiz ID
	UUID        string `db:"uuid" json:"uuid,omitempty"`     // Unique quiz UID (to store extracted from database)
	Name        string `db:"name" json:"name"`               // Quiz name
	Description string `db:"description" json:"description"` // Quiz description
	Version     string `db:"version" json:"version"`         // Quiz version
	IsActive    bool   `db:"is_active" json:"is_active"`     // Active = users can take the quiz
	Value       string
	Link        string
}
