package domain

type Post struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}
