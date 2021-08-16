package maven

import (
	"gitlabAnalyzer/dto"
	"gitlabAnalyzer/persistence"
	"gitlabAnalyzer/persistence/model"
	"io/ioutil"
	"log"
	"os"
)

type Maven struct{}

func (m Maven) ShouldApply(data dto.AnalysisData) bool {
	log.Println("for checking, wether", data.Path, "contains a Maven-Project a pom-File is look for")
	return pathHasPom(data.Path, 0)
}

func pathHasPom(path string, depth uint) bool {
	if depth > 2 {
		return false
	}
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("error while checking wether Maven Analyzer Should Apply")
		return false
	}
	for _, info := range dirContent {
		if info.IsDir() {
			hasPom := pathHasPom(path+string(os.PathSeparator)+info.Name(), depth+1)
			if hasPom {
				return true
			}
		} else if info.Name() == "pom.xml" {
			return true
		}
	}
	return false
}

func (m Maven) Apply(data dto.AnalysisData, result *model.AnalysisResult) {
	result.Maven = true
	persistence.SaveAnalysisResult(result)
	//TODO Maven-Module ermitteln und Dependencies abfragen
}

func (m Maven) NotApplied(data dto.AnalysisData, result *model.AnalysisResult) {
	result.Maven = false
	persistence.SaveAnalysisResult(result)
}

func (m Maven) Name() string {
	return "Maven"
}
