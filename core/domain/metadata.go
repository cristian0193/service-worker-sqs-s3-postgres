package domain

// MetaData represents the dto.
type MetaData struct {
	TrackID  string `json:"trackid"`
	Bucket   string `json:"bucket"`
	FileName string `json:"filename"`
	Key      string `json:"key"`
	Size     int64  `json:"size"`
}
