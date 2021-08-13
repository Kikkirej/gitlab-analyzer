package main

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/git"
	"gitlabAnalyzer/git/gitlab_api"
	"gitlabAnalyzer/persistence"
	"gitlabAnalyzer/settings"
	"log"
	"os"
	"strings"
)

func main() {
	settings.InitSettings()
	projects := gitlab_api.GetProjects()
	log.Println("number of identified projects:", len(projects))
	for _, project := range projects {
		if shouldBeAnalyzed(project) {
			log.Println("analyze project:", project.PathWithNamespace, " (", project.ID, ")")
			handleProject(project)
		} else {
			//log.Println("project does not meet criteria:", project.PathWithNamespace, "(", project.ID, ")")
		}
	}
}

func handleProject(project *gitlab.Project) {
	clonePath, _ := git.Clone(project)
	if clonePath == "error" {
		return
	}
	dbProject := persistence.DBObjectOf(project)
	log.Println(dbProject.PackagesEnabled)
	err := os.RemoveAll(clonePath)
	if err != nil {
		log.Println("deletion not possible:", clonePath)
	}
}

func shouldBeAnalyzed(project *gitlab.Project) bool {
	if settings.Struct.GitlabProjectRoot == "" {
		return true
	}
	if strings.HasPrefix(project.PathWithNamespace, settings.Struct.GitlabProjectRoot) {
		return true
	}
	return false
}
