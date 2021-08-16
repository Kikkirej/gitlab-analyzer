package gitlab_api

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/settings"
	"log"
)

var gitlabClient = initGitlabClient()

func initGitlabClient() *gitlab.Client {
	settings.InitSettings()
	gitlabClient, err := gitlab.NewClient(settings.Struct.GitlabPersonalToken, gitlab.WithBaseURL(settings.Struct.GitlabBaseurl+"api/v4"))
	if err != nil {
		log.Fatalln("couldn't connect to gitlab:", err)
	}
	return gitlabClient
}

func getGitlabClient() *gitlab.Client {
	if gitlabClient == nil {
		return initGitlabClient()
	}
	return gitlabClient
}
