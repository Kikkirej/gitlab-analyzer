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
	Analysis         AnalysisResult `gorm:"foreignKey:AnalysisId"`
	AnalysisId       uint
}

type MavenModuleDependency struct {
	gorm.Model
	Scope        string
	Version      string
	Dependency   MavenDependency `gorm:"foreignKey:DependencyID"`
	DependencyID uint
	Module       MavenModule
}

type MavenDependency struct {
	gorm.Model
	ArtifactID string
	GroupID    string
}
