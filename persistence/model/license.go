package model

type License struct {
	Name        string
	ID          string
	Url         string
	Deprecated  bool
	OsiApproved bool
	FsfApproved bool
	Spdx        bool
}
