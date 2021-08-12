package gitlab_api

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/settings"
	"log"
)

func GetProjects() []*gitlab.Project {
	git, _ := gitlab.NewClient(settings.Struct.GitlabPersonalToken, gitlab.WithBaseURL(settings.Struct.GitlabBaseurl+"api/v4"))
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    0,
		},
	}
	var projects []*gitlab.Project
	projects, header, err := git.Projects.ListProjects(opt)
	log.Println("handled page 1 of", header.TotalPages, "(total pages are not known unless first page is handled)")
	if err != nil {
		log.Fatalln("error while fetching Gitlab-projects", err)
	}
	for page := 2; page <= header.TotalPages; page++ {
		log.Println("handle page", page, "of", header.TotalPages)
		optOtherPages := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		}
		addedProjects, _, err := git.Projects.ListProjects(optOtherPages)
		if err != nil {
			log.Fatalln("error while fetching Gitlab-projects", err)
		}
		projects = append(projects, addedProjects...)
	}
	if err != nil {
		log.Fatalln("error while fetching Gitlab-projects", err)
	}
	return projects
}
