package gitlab_api

import (
	"github.com/xanzy/go-gitlab"
	"log"
)

func Projects() []*gitlab.Project {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    0,
		},
	}
	var projects []*gitlab.Project
	projects, header, err := getGitlabClient().Projects.ListProjects(opt)
	if err != nil {
		log.Fatalln("error while fetching Gitlab-projects", err)
	}
	log.Println("handled page 1 of", header.TotalPages, "(total pages are not known unless first page is handled)")
	for page := 2; page <= header.TotalPages; page++ {
		log.Println("handle page", page, "of", header.TotalPages)
		optOtherPages := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		}
		addedProjects, _, err := getGitlabClient().Projects.ListProjects(optOtherPages)
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
