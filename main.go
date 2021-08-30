package main

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/analyzer"
	"gitlabAnalyzer/dto"
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
	projects := gitlab_api.Projects()
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
