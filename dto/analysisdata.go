package dto

import (
	"github.com/go-git/go-git/v5"
	"gitlabAnalyzer/persistence/model"
)

type AnalysisData struct {
	Project *model.Project
	Branch  *model.Branch
	Path    string
	Repo    *git.Repository
}
