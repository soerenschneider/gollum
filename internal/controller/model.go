package controller

import (
	"strings"

	"github.com/soerenschneider/gollum/internal/github"
)

type ReleaseArtifacts struct {
	Release  github.Release
	Packages []github.Package
	Assets   []github.ReleaseAsset
}

func (r *ReleaseArtifacts) IsEmpty() bool {
	return strings.TrimSpace(r.Release.TagName) == ""
}
