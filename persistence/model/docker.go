package model

import "gorm.io/gorm"

type Dockerfile struct {
	gorm.Model
	Path         string
	Analysis     *AnalysisResult `gorm:"foreignKey:AnalysisId"`
	AnalysisId   uint
	LatestFrom   *Dockerimage `gorm:"foreignKey:LatestFromId"`
	LatestFromId uint
	Tag          string
}

type Dockerimage struct {
	gorm.Model
	Image     string
	License   *License `gorm:"foreignKey:LicenseID"`
	LicenseID *uint
}
