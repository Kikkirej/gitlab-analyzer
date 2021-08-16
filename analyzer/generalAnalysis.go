package analyzer

import (
	"gitlabAnalyzer/analyzer/maven"
	"gitlabAnalyzer/dto"
	"gitlabAnalyzer/persistence"
	"gitlabAnalyzer/persistence/model"
	"log"
)

var (
	analyzers = []Analyzer{maven.Maven{}}
)

func AnalyzeBranch(data dto.AnalysisData) {
	analysisResult := persistence.CreateAnalysisAndConnectToBranch(data.Branch)
	println(analysisResult) //TODO delete me
	for _, analyzer := range analyzers {
		if analyzer.ShouldApply(data) {
			log.Println("Project:", data.Project.Name, "- Branch:", data.Branch.Name, "- Analyzer:", analyzer.Name())
			analyzer.Apply(data, analysisResult)
		} else {
			analyzer.NotApplied(data, analysisResult)
		}
	}
}

type Analyzer interface {
	ShouldApply(data dto.AnalysisData) bool
	Apply(data dto.AnalysisData, result *model.AnalysisResult)
	// NotApplied in the case something needs to set to false or something
	NotApplied(data dto.AnalysisData, result *model.AnalysisResult)
	Name() string
}
