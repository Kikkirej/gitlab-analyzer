package model

import (
	"gorm.io/gorm"
	"time"
)

type Branch struct {
	gorm.Model
	Analyzed          bool
	Name              string
	Project           *Project `gorm:"foreignKey:ProjectId"`
	ProjectId         int
	CurrentAnalysis   *AnalysisResult `gorm:"foreignKey:CurrentAnalysisId"`
	CurrentAnalysisId *uint

	WebUrl             string
	DevelopersCanMerge bool
	DevelopersCanPush  bool
	Protected          bool
	DefaultBranch      bool
	Merged             bool

	LastCommitTime    *time.Time
	LastCommitShortID string
}
