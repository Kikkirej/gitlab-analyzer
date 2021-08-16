package model

import (
	"gorm.io/gorm"
	"time"
)

type Branch struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	Analyzed  bool
	Name      string
	Project   Project `gorm:"foreignKey:ProjectId"`
	ProjectId int
	CreatedAt time.Time
	UpdatedAt time.Time

	WebUrl             string
	DevelopersCanMerge bool
	DevelopersCanPush  bool
	Protected          bool
	DefaultBranch      bool
	Merged             bool

	LastCommitTime    *time.Time
	LastCommitShortID string
}
