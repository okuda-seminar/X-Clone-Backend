package entities

type Mute struct {
	SourceUserID string `gorm:"type:varchar(255);not null;primaryKey"`
	TargetUserID string `gorm:"type:varchar(255);not null;primaryKey"`
}
