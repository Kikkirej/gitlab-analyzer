package model

import "gorm.io/gorm"

type MavenModule struct {
	gorm.Model
	Path             string
	ArtifactID       string
	GroupID          string
	Version          string
	ParentArtifactID string
	ParentGroupID    string
	ParentVersion    string
	Analysis         *AnalysisResult `gorm:"foreignKey:AnalysisId"`
	AnalysisId       uint
	Packaging        string
}

type MavenModuleDependency struct {
	gorm.Model
	Scope        string
	Version      string
	Packaging    string
	Depth        uint
	Dependency   *MavenDependency `gorm:"foreignKey:DependencyID"`
	DependencyID uint
	Module       *MavenModule `gorm:"foreignKey:ModuleID"`
	ModuleID     uint
	Parent       *MavenModuleDependency `gorm:"foreignKey:ParentID"`
	ParentID     uint
}

type MavenDependency struct {
	gorm.Model
	GroupID    string
	ArtifactID string
}
