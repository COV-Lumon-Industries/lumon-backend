package schemas

import "time"

type DocumentResponse struct {
	ID             string     `json:"id"`
	ContentSummary string     `json:"content_summary"`
	Type           string     `json:"type"`
	UploadedAt     *time.Time `json:"uploaded_at"`
	UserID         string     `json:"user_id"`
}
