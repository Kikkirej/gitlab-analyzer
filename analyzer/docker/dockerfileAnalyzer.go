package docker

import (
	"github.com/kikkirej/gitlab-analyzer/dto"
	"github.com/kikkirej/gitlab-analyzer/persistence"
	"github.com/kikkirej/gitlab-analyzer/persistence/model"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DockerFile struct{}

const maxSearchDepth = 5

func (d DockerFile) ShouldApply(data dto.AnalysisData) bool {
	return pathHasDockerfile(data.Path, 0)
}

func pathHasDockerfile(path string, depth uint) bool {
	if depth >= maxSearchDepth {
		return false
	}
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("error while checking whether Maven Analyzer Should Apply")
		return false
	}
	for _, info := range dirContent {
		if info.IsDir() && info.Name() != ".git" {
			hasDockerfile := pathHasDockerfile(path+string(os.PathSeparator)+info.Name(), depth+1)
			if hasDockerfile {
				return true
			}
		} else if strings.Contains(strings.ToLower(info.Name()), "dockerfile") {
			return true
		}
	}
	return false
}

func (d DockerFile) Apply(data dto.AnalysisData, result *model.AnalysisResult) {
	result.Docker = true
	persistence.SaveAnalysisResult(result)
	dockerfilePaths := dockerfilesInPath(data.Path, string(os.PathSeparator), 0, []string{})
	processDockerFiles(dockerfilePaths, data, result)
}

func processDockerFiles(dockerfilePaths []string, data dto.AnalysisData, result *model.AnalysisResult) {
	for _, dockerfilePath := range dockerfilePaths {
		dockerfile := getDockerfile(dockerfilePath, result)
		persistence.SaveDockerfile(dockerfile)
	}
}

func getDockerfile(path string, result *model.AnalysisResult) *model.Dockerfile {
	var dockerfiles []model.Dockerfile
	persistence.GetDockerfile(path, result, &dockerfiles)
	if len(dockerfiles) >= 1 {
		return &dockerfiles[0]
	} else {
		return &model.Dockerfile{Path: path, Analysis: result}
	}
}

func dockerfilesInPath(basePath string, searchPath string, depth uint, files []string) []string {
	if depth >= maxSearchDepth {
		return files
	}
	dirContent, err := ioutil.ReadDir(basePath + searchPath)
	if err != nil {
		return files
	}
	for _, info := range dirContent {
		if info.IsDir() && info.Name() != ".git" {
			files = dockerfilesInPath(basePath, searchPath+info.Name()+string(os.PathSeparator), depth+1, files)
		} else if strings.Contains(strings.ToLower(info.Name()), "dockerfile") {
			files = append(files, searchPath+info.Name())
		}
	}
	return files
}

func (d DockerFile) NotApplied(data dto.AnalysisData, result *model.AnalysisResult) {
	result.Docker = false
	persistence.SaveAnalysisResult(result)
}

func (d DockerFile) Name() string {
	return "DockerFile"
}
