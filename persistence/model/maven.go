package model

import (
	"gorm.io/gorm"
)

type MavenModule struct {
	gorm.Model
	Path                 string
	ArtifactID           string
	GroupID              string
	Version              string
	ParentArtifactID     string
	ParentGroupID        string
	ParentVersion        string
	Analysis             *AnalysisResult `gorm:"foreignKey:AnalysisId"`
	AnalysisId           uint
	Packaging            string
	Repository           *MavenDistributionManagement `gorm:"foreignKey:RepositoryID"`
	RepositoryID         *uint
	SnapshotRepository   *MavenDistributionManagement `gorm:"foreignKey:SnapshotRepositoryID"`
	SnapshotRepositoryID *uint
}

type MavenDistributionManagement struct {
	gorm.Model
	RepoID string
	Name   string
	URL    string
}

type MavenModuleDependency struct {
	gorm.Model
	Scope        string
	Version      string
	Packaging    string
	Depth        uint
	Dependency   *MavenDependency `gorm:"foreignKey:DependencyID"`
	DependencyID uint
	Module       *MavenModule `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModuleID     uint
	Parent       *MavenModuleDependency `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ParentID     uint
}

type MavenDependency struct {
	gorm.Model
	GroupID    string
	ArtifactID string
}
