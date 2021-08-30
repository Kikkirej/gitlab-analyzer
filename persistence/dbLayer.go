package persistence

import (
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/persistence/model"
	"gitlabAnalyzer/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db = initDb()

func initDb() *gorm.DB {
	settings.InitSettings()
	dsn := "host=" + settings.Struct.PostgresHost +
		" user=" + settings.Struct.PostgresUser +
		" password=" + settings.Struct.PostgresPassword +
		" dbname=" + settings.Struct.PostgresDbname +
		" port=" + settings.Struct.PostgresPort +
		" sslmode=" + settings.Struct.PostgresSslmode +
		" TimeZone=Europe/Berlin"
	dbCon, errOpen := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errOpen != nil {
		log.Fatalln("error while connecting to database:", errOpen)
	}
	errAutoMigrateProject := dbCon.AutoMigrate(&model.Project{})
	if errAutoMigrateProject != nil {
		log.Fatalln("error while initializing to database:", errAutoMigrateProject)
	}
	errAutoMigrateBranch := dbCon.AutoMigrate(&model.Branch{})
	if errAutoMigrateBranch != nil {
		log.Fatalln("errAutoMigrateBranchor while initializing to database:", errAutoMigrateBranch)
	}
	errAutoMigrateMavenModule := dbCon.AutoMigrate(&model.MavenModuleDependency{})
	if errAutoMigrateMavenModule != nil {
		log.Fatalln("errAutoMigrateMavenModuleor while initializing to database:", errAutoMigrateMavenModule)
	}
	return dbCon
}

func DBObjectOfProject(apiInformation *gitlab.Project) *model.Project {
	var result *model.Project = nil
	db.First(&result, apiInformation.ID)
	if result.ID == 0 {
		result = &model.Project{ID: apiInformation.ID}
		updateProjectFields(result, apiInformation)
		db.Create(result)
	} else {
		updateProjectFields(result, apiInformation)
		db.Save(result)
	}
	return result
}

func DBObjectOfBranch(project *model.Project, branch *gitlab.Branch, analyzed bool) *model.Branch {
	var branches []*model.Branch
	db.Where("Name = ? AND project_id = ?", branch.Name, project.ID).Find(&branches)
	if len(branches) > 0 {
		if len(branches) > 1 {
			log.Println("For the project", project.Name, "(", project.ID, ") exist multiple objects for branch", branch.Name, "! This must not be!")
		}
		dbBranch := branches[0]
		dbBranch.Protected = branch.Protected
		dbBranch.DevelopersCanMerge = branch.DevelopersCanMerge
		dbBranch.DevelopersCanPush = branch.DevelopersCanPush
		dbBranch.Analyzed = analyzed
		dbBranch.DefaultBranch = branch.Default
		dbBranch.Merged = branch.Merged
		dbBranch.LastCommitTime = branch.Commit.CommittedDate
		dbBranch.LastCommitShortID = branch.Commit.ShortID
		db.Save(dbBranch)
		return dbBranch
	} else {
		branch := &model.Branch{
			Project:            project,
			Analyzed:           analyzed,
			Name:               branch.Name,
			WebUrl:             branch.WebURL,
			DevelopersCanMerge: branch.DevelopersCanMerge,
			DevelopersCanPush:  branch.DevelopersCanPush,
			Protected:          branch.Protected,
			DefaultBranch:      branch.Default,
			Merged:             branch.Merged,
			LastCommitTime:     branch.Commit.CommittedDate,
			LastCommitShortID:  branch.Commit.ShortID,
		}
		if branch.CurrentAnalysis != nil && branch.CurrentAnalysis.ID == 0 {
			branch.CurrentAnalysis = nil
		}
		db.Create(branch)
		return branch
	}
}

func updateProjectFields(project *model.Project, apiInformation *gitlab.Project) {
	project.Description = apiInformation.Description
	project.SSHURLToRepo = apiInformation.SSHURLToRepo
	project.HTTPURLToRepo = apiInformation.HTTPURLToRepo
	project.WebURL = apiInformation.WebURL
	project.ReadmeURL = apiInformation.WebURL
	project.Name = apiInformation.Name
	project.NameWithNamespace = apiInformation.NameWithNamespace
	project.Path = apiInformation.Path
	project.PathWithNamespace = apiInformation.PathWithNamespace
	project.IssuesEnabled = apiInformation.IssuesEnabled
	project.OpenIssuesCount = apiInformation.OpenIssuesCount
	project.MergeRequestsEnabled = apiInformation.MergeRequestsEnabled
	project.JobsEnabled = apiInformation.JobsEnabled
	project.WikiEnabled = apiInformation.WikiEnabled
	project.SnippetsEnabled = apiInformation.SnippetsEnabled
	project.ResolveOutdatedDiffDiscussions = apiInformation.ResolveOutdatedDiffDiscussions
	project.ContainerRegistryEnabled = apiInformation.ContainerRegistryEnabled
	project.GitlabCreatedAt = apiInformation.CreatedAt
	project.LastActivityAt = apiInformation.LastActivityAt
	project.NamespaceName = apiInformation.Namespace.Name
	project.NamespacePath = apiInformation.Namespace.Path
	project.Archived = apiInformation.Archived
	project.LicenseURL = apiInformation.LicenseURL
	project.SharedRunnersEnabled = apiInformation.SharedRunnersEnabled
	project.ForksCount = apiInformation.ForksCount
	project.StarCount = apiInformation.StarCount
	project.OnlyAllowMergeIfAllDiscussionsAreResolved = apiInformation.OnlyAllowMergeIfAllDiscussionsAreResolved
	project.RemoveSourceBranchAfterMerge = apiInformation.RemoveSourceBranchAfterMerge
	project.LFSEnabled = apiInformation.LFSEnabled
	project.RequestAccessEnabled = apiInformation.RequestAccessEnabled
	if apiInformation.ForkedFromProject != nil {
		project.ForkedFromProject = apiInformation.ForkedFromProject.ID
	}
	project.PackagesEnabled = apiInformation.PackagesEnabled
	project.AutocloseReferencedIssues = apiInformation.AutocloseReferencedIssues
	project.SuggestionCommitMessage = apiInformation.SuggestionCommitMessage
}

func CreateAnalysisAndConnectToBranch(branch *model.Branch) *model.AnalysisResult {
	if branch.CurrentAnalysisId != 0 {
		var result *model.AnalysisResult = nil
		db.First(&result, branch.CurrentAnalysisId)
		if result.Version != settings.CurrentVersion() {
			result.Version = settings.CurrentVersion()
			db.Save(result)
		}
		return result
	}
	analysisResult := &model.AnalysisResult{}
	analysisResult.Version = settings.CurrentVersion()
	db.Create(analysisResult)
	branch.CurrentAnalysis = analysisResult
	db.Save(branch)
	return analysisResult
}

func SaveAnalysisResult(result *model.AnalysisResult) {
	db.Save(result)
}

func MavenModulesForAnalysis(mavenmodules *[]model.MavenModule, result *model.AnalysisResult) {
	db.Where("analysis_id=?", result.ID).Find(&mavenmodules)
}

func SaveMavenModule(module *model.MavenModule) {
	db.Save(module)
}

func DeleteMavenModule(mavenModule model.MavenModule) {
	db.Delete(&mavenModule)
}

func GetMavenDependency(groupId string, artifactId string) *model.MavenDependency {
	var result *model.MavenDependency = nil
	db.Where("group_id=? and artifact_id=?", groupId, artifactId).Find(&result)
	if result.ID == 0 {
		result = &model.MavenDependency{GroupID: groupId, ArtifactID: artifactId}
		db.Create(&result)
	}
	return result
}

func GetMavenModuleDependency(dependency *model.MavenDependency, module model.MavenModule, parent *model.MavenModuleDependency, version string, scope string, packaging string, depth uint) *model.MavenModuleDependency {
	result := searchMavenModuleDependency(dependency, module, parent, packaging, depth)
	if result.ID == 0 {
		result = &model.MavenModuleDependency{Dependency: dependency, Module: &module, Scope: scope, Version: version, Packaging: packaging, Depth: depth}
		if parent.ID != 0 {
			result.Parent = parent
		}
		db.Create(&result)
	} else if result.Version != version || result.Scope != scope {
		result.Version = version
		result.Scope = scope
		db.Save(&result)
	}
	return result
}

func searchMavenModuleDependency(dependency *model.MavenDependency, module model.MavenModule, parent *model.MavenModuleDependency, packaging string, depth uint) *model.MavenModuleDependency {
	where := db.Where("dependency_id=? and module_id=? and depth=? and packaging=?", dependency.ID, module.ID, depth, packaging)
	if parent.ID != 0 {
		where.Where("parent_id=?", parent.ID)
	} else {
		where.Where("scope=?", "")
	}
	var result *model.MavenModuleDependency = nil
	where.Find(&result)
	return result
}

func GetDependenciesForModule(module model.MavenModule, result *[]model.MavenModuleDependency) {
	db.Where("module_id=?", module.ID).Find(&result)
}

func DeleteMavenModuleDependency(tobedeleted model.MavenModuleDependency) {
	db.Delete(&tobedeleted)
}
