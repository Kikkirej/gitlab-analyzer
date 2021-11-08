package model

import "gorm.io/gorm"

type Dockerfile struct {
	gorm.Model
	Path       string
	Analysis   *AnalysisResult `gorm:"foreignKey:AnalysisId"`
	AnalysisId uint
	LatestFrom string
}
