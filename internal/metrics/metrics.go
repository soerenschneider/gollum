package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace       = "gollum"
	subsystemGitHub = "github"
	subsystemTekton = "tekton"
)

var (
	FilteredReleasesTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystemGitHub,
		Name:      "filtered_releases_available_total",
		Help:      "The total amount of releases after filtering",
	}, []string{"owner", "repo"})

	ReleasesAvailableTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystemGitHub,
		Name:      "releases_available_total",
		Help:      "The total amount of releases for a repository",
	}, []string{"owner", "repo"})

	GithubRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemGitHub,
		Name:      "requests_total",
		Help:      "The total amount of GitHub requests",
	}, []string{"owner", "repo"})

	GithubRequestErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemGitHub,
		Name:      "request_errors_total",
		Help:      "The total amount of failed GitHub requests",
	}, []string{"owner", "repo", "url"})

	PipelineRunCreationErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemTekton,
		Name:      "pipelineruns_creation_errors_total",
		Help:      "The total amount of errors while trying to create pipeline runs",
	}, []string{"owner", "repo", "ref"})

	PipelineRunsCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemTekton,
		Name:      "pipelineruns_created_total",
		Help:      "The total amount of pipeline runs created",
	}, []string{"owner", "repo", "ref"})
)
