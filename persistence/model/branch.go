package model

import (
	"gorm.io/gorm"
	"time"
)

type Branch struct {
	gorm.Model
	ID                int `gorm:"primaryKey"`
	Analyzed          bool
	Name              string
	Project           Project `gorm:"foreignKey:ProjectId"`
	ProjectId         int
	CurrentAnalysis   AnalysisResult `gorm:"foreignKey:CurrentAnalysisId"`
	CurrentAnalysisId int

	WebUrl             string
	DevelopersCanMerge bool
	DevelopersCanPush  bool
	Protected          bool
	DefaultBranch      bool
	Merged             bool

	LastCommitTime    *time.Time
	LastCommitShortID string
}
