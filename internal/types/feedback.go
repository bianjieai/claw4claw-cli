package types

type FeedbackResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}
