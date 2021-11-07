package maven

import (
	"github.com/kikkirej/gitlab-analyzer/dto"
	"github.com/kikkirej/gitlab-analyzer/persistence"
	"github.com/kikkirej/gitlab-analyzer/persistence/model"
	"github.com/kikkirej/gitlab-analyzer/settings"
	"github.com/kikkirej/go-tgf"
	"github.com/kikkirej/go-tgf/ast"
	"log"
	"os/exec"
	"strings"
	"sync"
)

var mutexTGF sync.Mutex

var mutexDependencyAnalysis sync.Mutex

func getAndCreateDependenciesFor(module model.MavenModule, data dto.AnalysisData) []model.MavenModuleDependency {
	dependenciesTgfPath := data.Path + module.Path + "dependencies.tgf"
	command := exec.Command(settings.Struct.MavenCommand, "-pl", ":"+module.ArtifactID, "dependency:tree", "-DoutputEncoding=utf-8", "-DoutputFile="+dependenciesTgfPath, "-DoutputType=tgf")
	command.Dir = data.Path + module.Path
	output, err := command.CombinedOutput()
	if err != nil {
		log.Println("error while getting dependencies:", err, "\nOutput:", string(output[:]))
		return []model.MavenModuleDependency{}
	}
	mutexTGF.Lock()
	_, rootNodes, allNodesUUID, err := tgf.ParseFile(dependenciesTgfPath)
	mutexTGF.Unlock()
	if err != nil {
		log.Println("error, while processing tgf:", err)
		return []model.MavenModuleDependency{}
	}
	if len(rootNodes) == 0 {
		return []model.MavenModuleDependency{}
	}

	mutexDependencyAnalysis.Lock()
	dependencyObjects := dependencyObjectsFromTree(getRootNode(rootNodes, module), allNodesUUID, []model.MavenModuleDependency{}, &model.MavenModuleDependency{}, module, 0)
	mutexDependencyAnalysis.Unlock()
	deleteUnusedDependencies(module, dependencyObjects)
	return dependencyObjects
}

func deleteUnusedDependencies(module model.MavenModule, dependenciesObjectsActual []model.MavenModuleDependency) {
	var dependenciesObjectsDB []model.MavenModuleDependency
	persistence.GetDependenciesForModule(module, &dependenciesObjectsDB)
	if len(dependenciesObjectsDB) == len(dependenciesObjectsActual) {
		return
	}
	for _, dependencyDB := range dependenciesObjectsDB {
		if existsActually(dependencyDB, dependenciesObjectsActual) == false {
			persistence.DeleteMavenModuleDependency(dependencyDB)
		}
	}
}

func existsActually(dependencyDB model.MavenModuleDependency, dependenciesObjectsActual []model.MavenModuleDependency) bool {
	for _, dependencyActual := range dependenciesObjectsActual {
		if dependencyActual.ParentID == dependencyDB.ParentID &&
			dependencyActual.DependencyID == dependencyDB.DependencyID &&
			dependencyActual.Packaging == dependencyDB.Packaging {
			return true
		}
	}
	return false
}

func getRootNode(nodes []ast.Node, module model.MavenModule) ast.Node {
	for _, node := range nodes {
		if strings.Contains(node.Label, module.GroupID+":"+module.ArtifactID+":") {
			return node
		}
	}
	return nodes[0]
}

func dependencyObjectsFromTree(root ast.Node, allNodes map[string]ast.Edge, result []model.MavenModuleDependency, parent *model.MavenModuleDependency, module model.MavenModule, depth uint) []model.MavenModuleDependency {
	dependency := getMavenModuleDependency(root, module, parent, depth)
	for _, edgeId := range root.OutboundEdgeIds {
		result = dependencyObjectsFromTree(allNodes[edgeId].OutboundNode, allNodes, result, dependency, module, depth+1)
	}
	dependencies := append(result, *dependency)
	return dependencies
}

func getMavenModuleDependency(root ast.Node, module model.MavenModule, parent *model.MavenModuleDependency, depth uint) *model.MavenModuleDependency {
	splitRoot := strings.Split(root.Label, ":")
	groupId := splitRoot[0]
	artifactId := splitRoot[1]
	packaging := splitRoot[2]
	version := splitRoot[3]
	scope := ""
	if len(splitRoot) == 5 {
		scope = splitRoot[4]
	}
	dependency := persistence.GetMavenDependency(groupId, artifactId)
	moduleDependency := persistence.GetMavenModuleDependency(dependency, module, parent, version, scope, packaging, depth)
	return moduleDependency
}
