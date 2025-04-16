package controller

import (
	"cmp"
	"errors"
	"time"

	gollumv1alpha1 "github.com/soerenschneider/gollum/api/v1alpha1"
	"github.com/soerenschneider/gollum/internal/github"
	"github.com/soerenschneider/gollum/internal/versionfilter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getVersionFilter(filter *gollumv1alpha1.VersionFilterSpec) (VersionFilter, error) {
	if filter == nil {
		return nil, errors.New("filter empty")
	}

	switch filter.Impl {
	case "semver":
		return versionfilter.NewSemver(filter.Arg)
	default:
		return &versionfilter.NoFilter{}, nil
	}
}

func buildReleaseRequest(data *gollumv1alpha1.Repository) github.RepoQuery {
	var since *time.Time
	if data.Status.LastCheck != nil {
		since = &data.Status.LastCheck.Time
	}

	var successfullyBuiltReleases []string
	if data.Spec.MemorizeReleases {
		successfullyBuiltReleases = getSatisfiedReleases(data)
	}

	return github.RepoQuery{
		Owner:          data.Spec.Owner,
		Repo:           data.Spec.Repository,
		Since:          since,
		IgnoreReleases: successfullyBuiltReleases,
	}
}

// a release is satisfied, if there are no missing artifacts
func getSatisfiedReleases(repo *gollumv1alpha1.Repository) []string {
	var ret []string

	hasContainerPipelineDefined := len(repo.Spec.PipelineNames[gollumv1alpha1.ArtifactsKeyPackagesContainer]) > 0
	hasReleaseAssetsPipelineDefined := len(repo.Spec.PipelineNames[gollumv1alpha1.ArtifactsKeyReleaseAssets]) > 0

	for version, run := range repo.Status.Releases {
		satisfiesContainerArtifacts := !hasContainerPipelineDefined || hasContainerPipelineDefined && !run.MissingArtifacts[gollumv1alpha1.ArtifactsKeyPackagesContainer]
		satisfiesReleaseAssetArtifacts := !hasReleaseAssetsPipelineDefined || hasReleaseAssetsPipelineDefined && !run.MissingArtifacts[gollumv1alpha1.ArtifactsKeyReleaseAssets]

		if satisfiesContainerArtifacts && satisfiesReleaseAssetArtifacts {
			ret = append(ret, version)
		}
	}

	return ret
}

func buildArtifactQuery(data *gollumv1alpha1.Repository, release github.Release) github.ArtifactQuery {
	return github.ArtifactQuery{
		Owner:   data.Spec.Owner,
		Repo:    data.Spec.Repository,
		Release: release,
	}
}

func isPipelineRunExpired(creationDate time.Time) bool {
	// TODO: make configurable
	expiry := time.Now().Add(-14 * 24 * time.Hour)
	return creationDate.Before(expiry)
}

func maxOrDefault(a, b time.Duration) time.Duration {
	// both values are zero valued, return the default
	if a == 0 && b == 0 {
		return time.Hour
	}

	// both values are not zero value, return the larger
	if a != 0 && b != 0 {
		if a > b {
			return a
		}
		return b
	}

	// one value is the zero value, return the first non-zero value
	return cmp.Or(a, b)
}

func initStatus(data *gollumv1alpha1.Repository) {
	if data == nil {
		return
	}

	if data.Status.Releases == nil {
		data.Status.Releases = map[string]*gollumv1alpha1.Release{}
	}

	if data.Status.Conditions == nil {
		data.Status.Conditions = []metav1.Condition{}
	}

	data.Status.LastCheck = &metav1.Time{Time: time.Now()}
}
