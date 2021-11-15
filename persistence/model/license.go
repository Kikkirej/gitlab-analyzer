package model

import "gorm.io/gorm"

type License struct {
	gorm.Model
	Name        string
	LicenseID   string
	Url         string
	Deprecated  bool
	OsiApproved bool
	Spdx        bool
}
