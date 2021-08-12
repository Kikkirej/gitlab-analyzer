package main

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/settings"
	"log"
	"strings"
)
import "gitlabAnalyzer/gitlab_api"

// TODO https://git-scm.com/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-go-git

func main() {
	settings.InitSettings()
	projects := gitlab_api.GetProjects()
	log.Println("number of identified projects:", len(projects))
	for _, project := range projects {
		if shouldBeAnalyzed(project) {
			log.Println("analyze Project:", project.PathWithNamespace, " (", project.ID, ")")

		} else {
			log.Println("project does not meet criteria:", project.PathWithNamespace, "(", project.ID, ")")
		}
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
