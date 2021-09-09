package maven

import (
	"encoding/xml"
	"github.com/kikkirej/gitlab-analyzer/dto"
	"github.com/kikkirej/gitlab-analyzer/persistence"
	"github.com/kikkirej/gitlab-analyzer/persistence/model"
	"github.com/kikkirej/mvnparser"
	"io/ioutil"
	"log"
)

func processPomFiles(data dto.AnalysisData, modulePaths []string, result *model.AnalysisResult) []model.MavenModule {
	var mavenModules []model.MavenModule
	var mavenModulesDB []model.MavenModule
	persistence.MavenModulesForAnalysis(&mavenModulesDB, result)
	//delete not found modules
	for _, dbModule := range mavenModulesDB {
		foundModulePath := false
		for _, modulePath := range modulePaths {
			if dbModule.Path == modulePath {
				foundModulePath = true
				continue
			}
		}
		if foundModulePath == false {
			persistence.DeleteMavenModule(dbModule)
		}
	}
	for _, modulePath := range modulePaths {
		pomFile, err := ioutil.ReadFile(data.Path + modulePath + "pom.xml")
		if err != nil {
			log.Println(data.Path+modulePath, ": could not be handled:", err)
			continue
		}
		var project mvnparser.MavenProject
		if err := xml.Unmarshal([]byte(pomFile), &project); err != nil {
			log.Println(data.Path+modulePath, ": could not be unmarshalled:", err)
			continue
		}

		mavenModule := GetMavenModule(modulePath, mavenModulesDB)
		mavenModule.ArtifactID = project.ArtifactId
		if project.GroupId == "" {
			mavenModule.GroupID = project.Parent.GroupId
		} else {
			mavenModule.GroupID = project.GroupId
		}

		if project.Version == "" {
			mavenModule.Version = project.Parent.Version
		} else {
			mavenModule.Version = project.Version
		}

		if project.Parent.GroupId != "" && project.Parent.ArtifactId != "" && project.Parent.Version != "" {
			mavenModule.ParentVersion = project.Parent.Version
			mavenModule.ParentGroupID = project.Parent.GroupId
			mavenModule.ParentArtifactID = project.Parent.ArtifactId
		}

		if mavenModule.Repository != nil && project.Repository.URL == "" {
			mavenModule.Repository = nil
		} else if mavenModule.Repository == nil && project.Repository.URL != "" {
			mavenModule.Repository = &model.MavenDistributionManagement{
				RepoID: project.Repository.Id,
				Name:   project.Repository.Name,
				URL:    project.Repository.URL,
			}
		} else if mavenModule.Repository != nil &&
			(mavenModule.Repository.RepoID != project.Repository.Id || mavenModule.Repository.Name != project.Repository.Name || mavenModule.Repository.URL != project.Repository.URL) {
			mavenModule.Repository.RepoID = project.Repository.Id
			mavenModule.Repository.Name = project.Repository.Name
			mavenModule.Repository.URL = project.Repository.URL
		}

		if mavenModule.SnapshotRepository != nil && project.SnapshotRepository.URL == "" {
			mavenModule.SnapshotRepository = nil
		} else if mavenModule.SnapshotRepository == nil && project.SnapshotRepository.URL != "" {
			mavenModule.SnapshotRepository = &model.MavenDistributionManagement{
				RepoID: project.SnapshotRepository.Id,
				Name:   project.SnapshotRepository.Name,
				URL:    project.SnapshotRepository.URL,
			}
		} else if mavenModule.SnapshotRepository != nil &&
			(mavenModule.SnapshotRepository.RepoID != project.SnapshotRepository.Id || mavenModule.SnapshotRepository.Name != project.SnapshotRepository.Name || mavenModule.SnapshotRepository.URL != project.SnapshotRepository.URL) {
			mavenModule.SnapshotRepository.RepoID = project.SnapshotRepository.Id
			mavenModule.SnapshotRepository.Name = project.SnapshotRepository.Name
			mavenModule.SnapshotRepository.URL = project.SnapshotRepository.URL
		}

		mavenModule.Packaging = project.Packaging

		mavenModule.Analysis = result
		persistence.SaveMavenModule(&mavenModule)
		mavenModules = append(mavenModules, mavenModule)
	}
	return mavenModules
}

func GetMavenModule(path string, allMavenModules []model.MavenModule) model.MavenModule {
	for _, module := range allMavenModules {
		if module.Path == path {
			return module
		}
	}
	return model.MavenModule{Path: path}
}
