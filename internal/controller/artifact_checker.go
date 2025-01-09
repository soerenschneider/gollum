package controller

import (
	"fmt"

	gollumv1alpha1 "github.com/soerenschneider/gollum/api/v1alpha1"
)

type DefaultReleaseArtifactChecker struct{}

func (c *DefaultReleaseArtifactChecker) HasValidArtifacts(artifacts ReleaseArtifacts, artifactType gollumv1alpha1.ArtifactType) (bool, error) {
	switch artifactType {
	case gollumv1alpha1.ArtifactsKeyReleaseAssets:
		return len(artifacts.Assets) > 0, nil
	case gollumv1alpha1.ArtifactsKeyPackagesContainer:
		return len(artifacts.Packages) > 0, nil
	}

	return false, fmt.Errorf("no such artifactType %q", artifactType)
}
