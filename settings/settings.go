package settings

import (
	"log"
	"os"
	"strings"
	"sync"
)

const GitlabBaseUrlEnvName = "GITLAB_BASE_URL"
const GitlabPersonalTokenEnvName = "GITLAB_PERSONAL_TOKEN"
const GitlabProjectRootEnvName = "GITLAB_PROJECT_ROOT"
const WorkingDirEnvName = "WORKING_DIR"
const BranchesToAnalyzeEnvName = "BRANCHES_TO_ANALYZE"

const PostgresHostEnvName = "POSTGRES_HOST"
const PostgresUserEnvName = "POSTGRES_USER"
const PostgresPasswordEnvName = "POSTGRES_PASSWORD"
const PostgresDbnameEnvName = "POSTGRES_DBNAME"
const PostgresPortEnvName = "POSTGRES_PORT"
const PostgresSslmodeEnvName = "POSTGRES_SSLMODE"

const VERSION uint = 1

func CurrentVersion() uint {
	return VERSION
}

type SettingsStruct struct {
	Initialized         bool
	GitlabBaseurl       string
	GitlabPersonalToken string
	GitlabProjectRoot   string
	WorkingDir          string
	BranchesToAnalyze   []string

	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDbname   string
	PostgresPort     string
	PostgresSslmode  string
}

var Struct = SettingsStruct{Initialized: false}

var settingsInit sync.Mutex

func InitSettings() {
	settingsInit.Lock()
	if Struct.Initialized {
		settingsInit.Unlock()
		return
	}
	Struct.Initialized = true
	gitlabBaseUrl := os.Getenv(GitlabBaseUrlEnvName)

	if gitlabBaseUrl == "" {
		log.Fatal(GitlabBaseUrlEnvName, " not set. Exiting.")
	} else if !strings.HasSuffix(gitlabBaseUrl, "/") {
		gitlabBaseUrl = gitlabBaseUrl + "/"
	}

	log.Println("Gitlab Base-Url: ", gitlabBaseUrl)
	Struct.GitlabBaseurl = gitlabBaseUrl

	gitlabPersonalToken := os.Getenv(GitlabPersonalTokenEnvName)
	if gitlabPersonalToken == "" {
		log.Fatal(GitlabPersonalTokenEnvName, " not set. Exiting.")
	}
	Struct.GitlabPersonalToken = gitlabPersonalToken

	gitlabProjectRoot := os.Getenv(GitlabProjectRootEnvName)
	if gitlabProjectRoot != "" {
		log.Println("Gitlab-Project-Root: ", gitlabProjectRoot)
		Struct.GitlabProjectRoot = gitlabProjectRoot
	} else {
		log.Println("Gitlab-Project-Root not set")
	}

	workingDir := os.Getenv(WorkingDirEnvName)
	if workingDir == "" {
		log.Println(WorkingDirEnvName, "not set. Using default value: /tmp/")
		Struct.WorkingDir = "/tmp/"
	} else {
		if !strings.HasSuffix(workingDir, "/") && !strings.HasSuffix(workingDir, "\\") {
			builder := strings.Builder{}
			builder.WriteString(workingDir)
			builder.WriteRune(os.PathSeparator)
			Struct.WorkingDir = builder.String()
		} else {
			Struct.WorkingDir = workingDir
		}
	}

	branchesToAnalyze := os.Getenv(BranchesToAnalyzeEnvName)
	if branchesToAnalyze != "" {
		branchesToAnalyze = strings.ReplaceAll(branchesToAnalyze, ",", ";")
		Struct.BranchesToAnalyze = strings.Split(branchesToAnalyze, ";")
	} else {
		Struct.BranchesToAnalyze = []string{}
	}

	initPostgresSettings()
	settingsInit.Unlock()
}

func initPostgresSettings() {
	Struct.PostgresHost = os.Getenv(PostgresHostEnvName)
	if Struct.PostgresHost == "" {
		Struct.PostgresHost = "localhost"
	}
	Struct.PostgresPort = os.Getenv(PostgresPortEnvName)
	if Struct.PostgresPort == "" {
		Struct.PostgresPort = "5432"
	}
	Struct.PostgresUser = os.Getenv(PostgresUserEnvName)
	if Struct.PostgresUser == "" {
		Struct.PostgresUser = "postgres"
	}
	Struct.PostgresPassword = os.Getenv(PostgresPasswordEnvName)
	Struct.PostgresDbname = os.Getenv(PostgresDbnameEnvName)
	if Struct.PostgresDbname == "" {
		Struct.PostgresDbname = "postgres"
	}
	Struct.PostgresSslmode = os.Getenv(PostgresSslmodeEnvName)
	if Struct.PostgresSslmode == "" {
		Struct.PostgresSslmode = "disable"
	}
}
