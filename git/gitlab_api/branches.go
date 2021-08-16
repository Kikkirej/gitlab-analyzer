package gitlab_api

import (
	"github.com/xanzy/go-gitlab"
	"log"
)

func Branches(project *gitlab.Project) []*gitlab.Branch {
	branches, response, _ := getGitlabClient().Branches.ListBranches(project.ID, &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	})
	if response.TotalPages > 1 {
		log.Println("WARNING: The project has ", project.ID, "100+ Branches. (", response.TotalItems, "Branches ). This is currently not supported. WTF? Why so many branches?")
	}
	return branches
}
