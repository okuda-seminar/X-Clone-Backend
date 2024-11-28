package entities

type Block struct {
	SourceUserID string `gorm:"type:varchar(255);not null;index"`
	TargetUserID string `gorm:"type:varchar(255);not null;index"`
}
