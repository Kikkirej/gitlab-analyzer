package settings

import (
	"log"
	"os"
	"strings"
)

const GitlabBaseUrlEnvName = "GITLAB_BASE_URL"
const GitlabPersonalTokenEnvName = "GITLAB_PERSONAL_TOKEN"
const GitlabProjectRootEnvName = "GITLAB_PROJECT_ROOT"
const WorkingDirEnvName = "WORKING_DIR"

type SettingsStruct struct {
	GitlabBaseurl       string
	GitlabPersonalToken string
	GitlabProjectRoot   string
	WorkingDir          string
}

var Struct = SettingsStruct{}

func InitSettings() {
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
}
