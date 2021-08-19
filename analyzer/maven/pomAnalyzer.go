package maven

import (
	"encoding/xml"
	"github.com/creekorful/mvnparser"
	"gitlabAnalyzer/dto"
	"gitlabAnalyzer/persistence"
	"gitlabAnalyzer/persistence/model"
	"io/ioutil"
	"log"
)

func processPomFiles(data dto.AnalysisData, modulePaths []string, result *model.AnalysisResult) []model.MavenModule {
	var mavenModules []model.MavenModule

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
		mavenModule := model.MavenModule{Path: modulePath, ArtifactID: project.ArtifactId}
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

		mavenModule.Analysis = *result
		persistence.SaveMavenModule(&mavenModule)
		mavenModules = append(mavenModules, mavenModule)
	}
	return mavenModules
}
