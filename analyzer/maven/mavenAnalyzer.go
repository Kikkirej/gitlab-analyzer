package maven

import (
	"github.com/kikkirej/gitlab-analyzer/dto"
	"github.com/kikkirej/gitlab-analyzer/persistence"
	"github.com/kikkirej/gitlab-analyzer/persistence/model"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Maven struct{}

const moduleMaxSearchDepth = 3

func (m Maven) ShouldApply(data dto.AnalysisData) bool {
	log.Println("for checking, wether", data.Path, "contains a Maven-Project a pom-File is look for")
	return pathHasPom(data.Path, 0)
}

func pathHasPom(path string, depth uint) bool {
	if depth >= moduleMaxSearchDepth {
		return false
	}
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("error while checking whether Maven Analyzer Should Apply")
		return false
	}
	for _, info := range dirContent {
		if info.IsDir() && info.Name() != ".git" {
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
	modulePaths := mavenModulesInPath(data.Path, string(os.PathSeparator), 0, []string{})
	mavenModules := processPomFiles(data, modulePaths, result)
	processDependencies(data, mavenModules)
}

func processDependencies(data dto.AnalysisData, modules []model.MavenModule) {
	var wg sync.WaitGroup
	for _, module := range modules {
		wg.Add(1)
		go func(module model.MavenModule, data dto.AnalysisData) {

			defer wg.Done()
			getAndCreateDependenciesFor(module, data)
		}(module, data)
	}
	wg.Wait()
}

func mavenModulesInPath(basePath string, searchPath string, depth uint, modules []string) []string {
	if depth >= moduleMaxSearchDepth {
		return modules
	}
	dirContent, err := ioutil.ReadDir(basePath + searchPath)
	if err != nil {
		return modules
	}
	for _, info := range dirContent {
		if info.IsDir() && info.Name() != ".git" {
			modules = mavenModulesInPath(basePath, searchPath+info.Name()+string(os.PathSeparator), depth+1, modules)
		} else if info.Name() == "pom.xml" {
			modules = append(modules, searchPath)
		}
	}
	return modules
}

func (m Maven) NotApplied(_ dto.AnalysisData, result *model.AnalysisResult) {
	result.Maven = false
	persistence.SaveAnalysisResult(result)
}

func (m Maven) Name() string {
	return "Maven"
}
