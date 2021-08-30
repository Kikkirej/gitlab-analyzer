package model

import (
	"gorm.io/gorm"
)

type AnalysisResult struct {
	gorm.Model
	Maven   bool
	Version uint
}
