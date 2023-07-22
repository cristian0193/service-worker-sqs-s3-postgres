package entity

// MetaData represents the entity.
type MetaData struct {
	TrackID  string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:trackid" json:"trackid"`
	Bucket   string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:bucket" json:"bucket"`
	FileName string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:filename" json:"filename"`
	Key      string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:key" json:"key"`
	Size     int64  `gorm:"NULL;TYPE:INT;COLUMN:size" json:"size"`
}

// TableName definition name for table .
func (MetaData) TableName() string {
	return "metadata"
}
