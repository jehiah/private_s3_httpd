package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

const VERSION = "1.1.0"

func VCSRevision() string {
	var version string
	var modified string
	buildInfo, _ := debug.ReadBuildInfo()
	for _, s := range buildInfo.Settings {
		if s.Key == "vcs.revision" && s.Value != "" {
			version = s.Value
		}
		if s.Key == "vcs.modified" && s.Value != "" {
			modified = s.Value
		}
	}
	return fmt.Sprintf("%s%s", version, func() string {
		if modified != "" {
			return "; modified"
		}
		return ""
	}())
}

func ShowVersion() string {
	return fmt.Sprintf("private_s3_httpd v%s (%s) (built w/%s)", VERSION, VCSRevision(), runtime.Version())
}
