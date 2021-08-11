package main

import (
	"gitlabAnalyzer/settings"
	"log"
)
import "gitlabAnalyzer/gitlab_api"

// TODO https://git-scm.com/book/en/v2/Appendix-B:-Embedding-Git-in-your-Applications-go-git

func main() {
	settings.InitSettings()
	projects := gitlab_api.GetProjects()
	log.Println("number of identified projects: ", len(projects) )
}
