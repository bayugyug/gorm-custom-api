package configs

import (
	"fmt"
)

var (
	// Application the app name
	Application = "Building API Service"
	// BuildTime is a time label of the moment when the binary was built
	BuildTime = ""
	// Commit is a last commit hash at the moment when the binary was built
	Commit = ""
	// Release is a semantic version of current build
	Release = ""
	// APIVersion is the app ver string
	APIVersion = ""
)

func init() {
	APIVersion = fmt.Sprintf("Ver: %s", Release)
}
