package docker

import (
	"github.com/kikkirej/gitlab-analyzer/dto"
	"github.com/kikkirej/gitlab-analyzer/persistence/model"
)

type DockerFile struct{}

func (d DockerFile) ShouldApply(data dto.AnalysisData) bool {
	//panic("implement me")
	return false
}

func (d DockerFile) Apply(data dto.AnalysisData, result *model.AnalysisResult) {
	//panic("implement me")
}

func (d DockerFile) NotApplied(data dto.AnalysisData, result *model.AnalysisResult) {
	//panic("implement me")
}

func (d DockerFile) Name() string {
	return "DockerFile"
}
