package model

import (
	"time"
)

type Project struct {
	ID                                        int `gorm:"primaryKey"`
	Description                               string
	SSHURLToRepo                              string
	HTTPURLToRepo                             string
	WebURL                                    string
	ReadmeURL                                 string
	Name                                      string
	NameWithNamespace                         string
	Path                                      string
	PathWithNamespace                         string
	IssuesEnabled                             bool
	OpenIssuesCount                           int
	MergeRequestsEnabled                      bool
	JobsEnabled                               bool
	WikiEnabled                               bool
	SnippetsEnabled                           bool
	ResolveOutdatedDiffDiscussions            bool
	ContainerRegistryEnabled                  bool
	GitlabCreatedAt                           *time.Time
	LastActivityAt                            *time.Time
	NamespaceName                             string
	NamespacePath                             string
	EmptyRepo                                 bool
	Archived                                  bool
	LicenseURL                                string
	SharedRunnersEnabled                      bool
	ForksCount                                int
	StarCount                                 int
	OnlyAllowMergeIfAllDiscussionsAreResolved bool
	RemoveSourceBranchAfterMerge              bool
	LFSEnabled                                bool
	RequestAccessEnabled                      bool
	ForkedFromProject                         int
	PackagesEnabled                           bool
	AutocloseReferencedIssues                 bool
	SuggestionCommitMessage                   string
}
