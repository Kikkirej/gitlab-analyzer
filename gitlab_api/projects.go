package gitlab_api

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/settings"
	"log"
)

func GetProjects() []*gitlab.Project {
	git, _ := gitlab.NewClient(settings.Struct.GitlabPersonalToken, gitlab.WithBaseURL(settings.Struct.GitlabBaseurl + "api/v4"))
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage:  100,
			Page: 0,
		},
	}
	var projects []*gitlab.Project
	projects, header, err := git.Projects.ListProjects(opt)
	if err != nil {
		log.Fatalln("Error while fetching Gitlab-Projects", err)
	}
	for page := 1; page <= header.TotalPages; page++ {
		log.Println("Verarbeite Seite ", page, " von ", header.TotalPages)
		optOtherPages := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage:  100,
				Page: page,
			},
		}
		addedProjects, _ , err := git.Projects.ListProjects(optOtherPages)
		if err != nil {
			log.Fatalln("Error while fetching Gitlab-Projects", err)
		}
		projects = append(projects, addedProjects...)
	}
	if err != nil {
		log.Fatalln("Error while fetching Gitlab-Projects", err)
	}
	return projects
}
