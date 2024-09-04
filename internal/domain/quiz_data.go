package domain

// Quiz data structure to store, extract serialize actual tests
type QuizData struct {
	Type        string         `db:"type" json:"type"`               // Json datatype tag
	ID          string         `db:"id" json:"id"`                   // Unique test ID
	UUID        string         `db:"uuid" json:"uuid,omitempty"`     // Unique test UID (to store extracted from database)
	Name        string         `db:"name" json:"name"`               // Test name
	Description string         `db:"description" json:"description"` // Test description
	Version     string         `db:"version" json:"version"`         // Test version
	IsActive    bool           `db:"is_active" json:"is_active"`     // Active flag: whether it is available to users
	Data        []QuizQuestion `db:"data" json:"data"`               // Test data array
	Scores      []QuizScore    `db:"scores" json:"scores"`           // Test data array
}

// Quiz question entry
type QuizQuestion struct {
	ID           string `db:"options" json:"id"`                  // Unique question ID
	UUID         string `db:"uuid" json:"uuid,omitempty"`         // Unique question UID (to store extracted from database)
	Text         string `db:"text" json:"text"`                   // Question text
	Comment      string `db:"comment" json:"comment,omitempty"`   // Question comment (i.e. floating help)
	Type         string `db:"type" json:"type"`                   // Json datatype tag: "single_choice", "multiple_choice" or "entry_int"
	CorrectValue string `db:"correct_value" json:"correct_value"` // Correct answers for open answers ("entry_int" etc.)  ////TODO: `json:"correct_value,omitempty"`
	//CorrectOptions []string     `json:"correct_options"`   // Correct answers for closed answers ("single_choice", "multiple_choice") by option ID ////TODO: `json:"correct_options,omitempty"`
	Options []QuizOption `db:"options" json:"options"` // Question answer options ////TODO: `json:"options,omitempty"`
}

// Quiz question option for closed answers ("single_choice", "multiple_choice")
type QuizOption struct {
	ID        string `db:"id" json:"id"`                           // Unique option ID
	UUID      string `db:"uuid" json:"uuid,omitempty"`             // Unique option UID (to store extracted from database)
	Text      string `db:"text,omitempty" json:"text,omitempty"`   // Option text
	Value     string `db:"value,omitempty" json:"value,omitempty"` // Option value (if needed)
	IsCorrect bool   `db:"is_correct" json:"is_correct"`           // Correct options have TRUE
}

// Quiz score entry
type QuizScore struct {
	ID         string `db:"id" json:"id"`                     // Unique score ID
	UUID       string `db:"uuid" json:"uuid,omitempty"`       // Unique score UID (to store extracted from database)
	MinPrecent int    `db:"min_percent" json:"min_percent"`   // Min percent of replies to hold this score
	MaxPrecent int    `db:"max_percent" json:"max_percent"`   // Max percent of replies to hold this score
	Score      int    `db:"score" json:"score"`               // Score to be set
	Pass       bool   `db:"pass" json:"pass"`                 // Quiz pass mark
	Comment    string `db:"comment" json:"comment,omitempty"` // Score comment (i.e. floating help)
}

//Quiz results
type QuizResults struct {
	ID       string `db:"id" json:"id"`               // Unique result ID
	TestID   string `db:"test_id" json:"test_id"`     // Test ID
	ScoreID  string `db:"score_id" json:"score_id"`   // Score ID
	PassTime string `db:"pass_time" json:"pass_time"` // Timestamp
	Result   int    `db:"result" json:"result"`
	Score    int    `db:"score" json:"score"`
	IsPassed bool   `db:"is_passed" json:"is_passed"`
	Replies  string `db:"replies" json:"replies"` // Replies as JSON string
}

// WITH question_data AS (
//     SELECT q.id AS id,
//     q.test_id AS test_id,
//            q.text AS text,
//            JSON_AGG(JSON_BUILD_OBJECT('id', o.id, 'text', o.text, 'is_correct', o.is_correct)) AS options
//     FROM public.questions q
//     LEFT JOIN public.options o ON o.question_id = q.id
//     GROUP BY q.id
// )
// SELECT t.id,
// 	t.ext_id,
// 	t."version",
// 	t.is_active,
// 	t."type",
// 	t."name",
// 	t.description,
//        COALESCE(JSON_AGG(JSON_BUILD_OBJECT('id', q.id, 'text', q.text, 'options', q.options)), '[]') AS questions
// FROM public.tests t
// LEFT JOIN question_data q ON q.test_id = t.id
// GROUP BY t.id, t.ext_id, t."version", t.is_active, t."type", t."name", t.description;
