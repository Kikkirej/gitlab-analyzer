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

func initDb() *gorm.DB { //TODO Felder in Variablen
	settings.InitSettings()
	dsn := "host=" + settings.Struct.PostgresHost +
		" user=" + settings.Struct.PostgresUser +
		" password=" + settings.Struct.PostgresPassword +
		" dbname=" + settings.Struct.PostgresDbname +
		" port=" + settings.Struct.PostgresPort +
		" sslmode=" + settings.Struct.PostgresSslmode +
		" TimeZone=Europe/Berlin"
	dbCon, errAutoMigrate := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errAutoMigrate != nil {
		log.Fatalln("error while connecting to database:", errAutoMigrate)
	}
	errAutoMigrate = dbCon.AutoMigrate(&model.Project{})
	if errAutoMigrate != nil {
		log.Fatalln("error while initalizing to database:", errAutoMigrate)
	}
	return dbCon
}

func DBObjectOf(apiInformation *gitlab.Project) *model.Project {
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

func updateProjectFields(project *model.Project, apiInformation *gitlab.Project) {
	project.Description = apiInformation.Description
	project.DefaultBranch = apiInformation.DefaultBranch
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
	//if apiInformation.ForkedFromProject != nil {
	//	project.ForkedFromProject = apiInformation.ForkedFromProject.ID
	//}
	project.PackagesEnabled = apiInformation.PackagesEnabled
	project.AutocloseReferencedIssues = apiInformation.AutocloseReferencedIssues
	project.SuggestionCommitMessage = apiInformation.SuggestionCommitMessage
}
