package form

// Comment represents a comment form.
type Comment struct {
	Content  string `json:"content"`
	Type     string `json:"type"`     // text, audio, video
	MediaURL string `json:"mediaURL"`
}
