package versionfilter

import (
	"github.com/Masterminds/semver/v3"
)

type NoFilter struct {
}

func (f *NoFilter) Matches(_ string) (bool, error) {
	return false, nil
}

type SemverReleaseFilter struct {
	constraint *semver.Constraints
}

func NewSemver(constraint string) (*SemverReleaseFilter, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return nil, err
	}

	return &SemverReleaseFilter{
		constraint: c,
	}, nil
}

func (f *SemverReleaseFilter) Matches(tag string) (bool, error) {
	version, err := semver.NewVersion(tag)
	if err != nil {
		return false, err
	}

	return f.constraint.Check(version), nil
}
