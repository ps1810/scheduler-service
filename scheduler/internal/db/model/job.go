package model

type CronJob struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"type:text;not null"`
	CronExpression string `gorm:"type:text;not null"`
	Enabled        bool   `gorm:"not null;default:true"`
	Table          string `gorm:"type:text;not null"`
	Field          string `gorm:"type:text;not null"`
	Aggregation    string `gorm:"type:text;not null"`
	Duration       string `gorm:"type:text;not null"`
	DurationFilter string `gorm:"type:text;not null"`
	CreatedAt      string `gorm:"type:datetime"`
	LastRun        string `gorm:"type:datetime"`
	NextRun        string `gorm:"type:datetime"`
}
