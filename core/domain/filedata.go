package domain

// FileData represents the dto.
type FileData struct {
	ID      *int64 `json:"id"`
	Message string `json:"message"`
	Owner   string `json:"owner"`
	Date    string `json:"date"`
}
