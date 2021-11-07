package main

import (
	"github.com/kikkirej/gitlab-analyzer/analyzer"
	"github.com/kikkirej/gitlab-analyzer/dto"
	"github.com/kikkirej/gitlab-analyzer/git"
	"github.com/kikkirej/gitlab-analyzer/git/gitlab_api"
	"github.com/kikkirej/gitlab-analyzer/persistence"
	"github.com/kikkirej/gitlab-analyzer/settings"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	settings.InitSettings()
	projects := gitlab_api.Projects()
	log.Println("number of identified projects:", len(projects))

	maxAnalyzedProjectsParrallel := 5
	guard := make(chan struct{}, maxAnalyzedProjectsParrallel)

	var wg sync.WaitGroup
	for _, project := range projects {
		if shouldBeAnalyzed(project) {
			log.Println("analyze project:", project.PathWithNamespace, " (", project.ID, ")")
			guard <- struct{}{}
			wg.Add(1)
			go func(project *gitlab.Project) {
				defer wg.Done()
				handleProject(project)
				<-guard
			}(project)
		} else {
			//log.Println("project does not meet criteria:", project.PathWithNamespace, "(", project.ID, ")")
		}
	}
	wg.Wait()
}

func handleProject(project *gitlab.Project) {
	clonePath, repo := git.Clone(project)
	if clonePath == "error" {
		return
	}
	dbProject := persistence.DBObjectOfProject(project)
	branches := gitlab_api.Branches(project)
	for _, branch := range branches {
		if branchShouldBeAnalyzed(branch) {
			branchDb := persistence.DBObjectOfBranch(dbProject, branch, true)
			err := git.CheckoutBranch(repo, branch)
			if err != nil {
				continue
			}
			analyzer.AnalyzeBranch(dto.AnalysisData{Project: dbProject, Branch: branchDb, Path: clonePath, Repo: repo})
		} else {
			persistence.DBObjectOfBranch(dbProject, branch, false)
		}

	}
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

func branchShouldBeAnalyzed(branch *gitlab.Branch) bool {
	if branch.Default {
		log.Println("default branch is always analyzed(", branch.Name, ")")
		return true
	}
	for _, configuredBranch := range settings.Struct.BranchesToAnalyze {
		if strings.Contains(branch.Name, configuredBranch) {
			log.Println("Branch", branch.Name, "is analyzed, because of pattern:", configuredBranch)
			return true
		}
	}
	return false
}
