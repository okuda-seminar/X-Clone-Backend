package entities

type Followship struct {
	SourceUserID string `gorm:"type:varchar(255);not null;primaryKey"`
	TargetUserID string `gorm:"type:varchar(255);not null;primaryKey"`
}
