package model

import (
	"gorm.io/gorm"
)

type AnalysisResult struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Maven    bool
	Version  uint
	BranchId int
}
