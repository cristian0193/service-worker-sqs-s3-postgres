package entity

// FileData represents the entity.
type FileData struct {
	ID      *int64 `gorm:"NULL;TYPE:INT;COLUMN:id" json:"id"`
	Message string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:message" json:"message"`
	Owner   string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:owner" json:"owner"`
	Date    string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:date" json:"date"`
}

// TableName definition name for table .
func (FileData) TableName() string {
	return "filedata"
}
