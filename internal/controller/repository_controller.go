package controller

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/soerenschneider/gollum/internal/github"
	"github.com/soerenschneider/gollum/internal/metrics"
	"github.com/soerenschneider/gollum/internal/requeue"
	"github.com/soerenschneider/gollum/internal/tekton"
	pool "github.com/sourcegraph/conc/pool"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/strings/slices"
	"knative.dev/pkg/apis"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	gollumv1alpha1 "github.com/soerenschneider/gollum/api/v1alpha1"
)

type PipelineRunner interface {
	GetPipeline(ctx context.Context, namespace, name string) (*pipelinev1.Pipeline, error)
	GetPipelineRun(ctx context.Context, namespace, name string) (*pipelinev1.PipelineRun, error)
	CreatePipelineRun(ctx context.Context, req tekton.CreatePipelineRunRequest) (*pipelinev1.PipelineRun, error)
}

type GithubClient interface {
	GetReleases(ctx context.Context, params github.RepoQuery) ([]github.Release, error)
	GetAssets(ctx context.Context, assetQuery github.ArtifactQuery) ([]github.ReleaseAsset, error)
	GetPackages(ctx context.Context, query github.ArtifactQuery) ([]github.Package, error)
}

type Requeue interface {
	Requeue(duration time.Duration) time.Duration
}

type VersionFilter interface {
	Matches(version string) (bool, error)
}

type ReleaseArtifactChecker interface {
	HasValidArtifacts(artifacts *ReleaseArtifacts, artifactType gollumv1alpha1.ArtifactType) (bool, error)
}

// RepositoryReconciler reconciles a Repository object
type RepositoryReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Recorder       record.EventRecorder
	PipelineRunner PipelineRunner
	GithubClient   GithubClient
	Requeue        Requeue

	DefaultRequeueInterval time.Duration
	DefaultJitterPercent   float64
}

// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="tekton.dev",resources=pipelines,verbs=get;list;watch
// +kubebuilder:rbac:groups="tekton.dev",resources=pipelineruns,verbs=create;patch;get;list;watch
// +kubebuilder:rbac:groups=gollum.soeren.cloud,resources=repositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gollum.soeren.cloud,resources=repositories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gollum.soeren.cloud,resources=repositories/finalizers,verbs=update
func (r *RepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	data := &gollumv1alpha1.Repository{}
	if err := r.Get(ctx, req.NamespacedName, data); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	logger := log.FromContext(ctx)
	defer func() {
		if err := r.Status().Update(ctx, data); err != nil {
			logger.Error(err, "could not update status")
		}
	}()

	initStatus(data)

	if err := r.checkIfPipelineExists(ctx, data, req.Namespace); err != nil {
		requeueAfter := requeue.JitterPercentageAdditive(r.Requeue.Requeue(r.DefaultRequeueInterval), r.DefaultJitterPercent)
		metrics.RequeueAfter.WithLabelValues(data.Spec.Owner, data.Spec.Repository).Set(requeueAfter.Seconds())
		logger.Error(err, "could not find desired pipeline, make sure to install it first", "owner", data.Spec.Owner, "repo", data.Spec.Repository, "requeue_after", requeueAfter)
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	r.cleanupRuns(ctx, req.Namespace, data)

	metrics.LastReleaseCheck.WithLabelValues(data.Spec.Owner, data.Spec.Repository).SetToCurrentTime()
	releases, rateLimitReset, err := r.getReleasesForRepository(ctx, data)
	if err != nil {
		requeueAfter := cmp.Or(rateLimitReset, requeue.JitterPercentageDistributed(r.Requeue.Requeue(r.DefaultRequeueInterval), r.DefaultJitterPercent))
		logger.Error(err, "could not get releases from Github", "requeue_after", requeueAfter)
		metrics.RequeueAfter.WithLabelValues(data.Spec.Owner, data.Spec.Repository).Set(requeueAfter.Seconds())
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	filteredReleases, err := r.applyVersionFilter(ctx, data, releases)
	if err != nil {
		logger.Error(err, "filtering releases produced errors, continuing with all releases")
	}
	logger.Info("Found unseen release(s)", "unseen", len(releases), "filtered", len(filteredReleases), "owner", data.Spec.Owner, "repo", data.Spec.Repository)

	releaseArtifacts, rateLimitReset := r.fetchArtifactDataForReleases(ctx, data, filteredReleases)
	releasesWithMissingArtifacts := r.checkReleaseDataForMissingArtifacts(data, releaseArtifacts)
	if len(releasesWithMissingArtifacts) == 0 {
		meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
			Type:    "NoRunsNeeded",
			Status:  metav1.ConditionTrue,
			Message: "No PipelineRun needs to be scheduled",
			Reason:  "NoMissingReleases",
		})

		requeueAfter := cmp.Or(rateLimitReset, requeue.JitterPercentageDistributed(r.Requeue.Requeue(r.DefaultRequeueInterval), r.DefaultJitterPercent))
		logger.Info("no releases with missing artifacts available", "owner", data.Spec.Owner, "repo", data.Spec.Repository, "requeue_after", requeueAfter)
		metrics.RequeueAfter.WithLabelValues(data.Spec.Owner, data.Spec.Repository).Set(requeueAfter.Seconds())
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	backoffDuration, err := r.createPipelineRunsForReleases(ctx, data, releasesWithMissingArtifacts, req.Namespace)
	if err != nil {
		requeueAfter := maxOrDefault(rateLimitReset, backoffDuration)
		logger.Error(err, "errors while creating pipelines for releases", "owner", data.Spec.Owner, "repo", data.Spec.Repository, "requeue_after", requeueAfter)
		metrics.RequeueAfter.WithLabelValues(data.Spec.Owner, data.Spec.Repository).Set(requeueAfter.Seconds())
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	requeueAfter := cmp.Or(rateLimitReset, requeue.JitterPercentageDistributed(r.Requeue.Requeue(r.DefaultRequeueInterval), r.DefaultJitterPercent))
	logger.Info("finished processing repository", "owner", data.Spec.Owner, "repo", data.Spec.Repository, "requeue_after", requeueAfter)
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *RepositoryReconciler) checkIfPipelineExists(ctx context.Context, data *gollumv1alpha1.Repository, namespace string) error {
	var errs error

	for _, pipelineName := range data.Spec.PipelineNames {
		_, err := r.PipelineRunner.GetPipeline(ctx, namespace, pipelineName)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	if errs != nil {
		reason := "Unknown"
		if errors.Is(errs, tekton.ErrTektonPipelineNotFound) {
			reason = "PipelineNotFound"
		} else if errors.Is(errs, tekton.ErrTektonGetPipelineForbidden) {
			reason = "Forbidden"
		}

		meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
			Type:    "TektonPipelineUnavailable",
			Status:  metav1.ConditionFalse,
			Reason:  reason,
			Message: "Requested Tekton Pipeline not available",
		})
	}

	return errs
}

// cleanupRuns cleans up outdated and PipelineRuns that have been deleted from the status field.
//
// It checks the repository's status for each release and the most recent pipeline runs associated with each artifact type.
// If a pipeline run cannot be found (e.g., it has been deleted or is no longer valid), the run is considered outdated and is marked for removal.
// The method ensures that only valid and existing pipeline runs are retained in the repository status.
func (r *RepositoryReconciler) cleanupRuns(ctx context.Context, namespace string, data *gollumv1alpha1.Repository) {
	pipelineRunsToRemove := map[string][]gollumv1alpha1.ArtifactType{}

	for version, release := range data.Status.Releases {
		for artifactType, runs := range release.MostRecentRuns {
			pipelineRun, err := r.PipelineRunner.GetPipelineRun(ctx, namespace, runs.Name)
			isExpired := pipelineRun != nil && isPipelineRunExpired(pipelineRun.CreationTimestamp.Time)
			isNotFound := errors.Is(err, tekton.ErrTektonPipelineNotFound)

			if isExpired || isNotFound {
				_, hasVersion := pipelineRunsToRemove[version]
				if !hasVersion {
					pipelineRunsToRemove[version] = []gollumv1alpha1.ArtifactType{}
				}
				pipelineRunsToRemove[version] = append(pipelineRunsToRemove[version], artifactType)
			}
		}
	}

	for version, artifactTypes := range pipelineRunsToRemove {
		for _, artifactType := range artifactTypes {
			if len(data.Status.Releases[version].MostRecentRuns) <= 1 {
				data.Status.Releases[version].MostRecentRuns = nil
			} else {
				delete(data.Status.Releases[version].MostRecentRuns, artifactType)
			}
		}
	}
}

func (r *RepositoryReconciler) getReleasesForRepository(ctx context.Context, data *gollumv1alpha1.Repository) ([]github.Release, time.Duration, error) {
	ghReleasesRequest := buildReleaseRequest(data)
	releases, err := r.GithubClient.GetReleases(ctx, ghReleasesRequest)
	if err != nil {
		log.FromContext(ctx).Error(err, "could not fetch release info from GitHub")

		var rlErr *github.RateLimitError
		if errors.As(err, &rlErr) {
			meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
				Type:    "FetchReleaseInformationFailed",
				Status:  metav1.ConditionFalse,
				Reason:  "RateLimitExceeded",
				Message: "Fetching release information from GitHub failed",
			})

			r.Recorder.Event(data, v1.EventTypeWarning, "RateLimitExceeded", "Could not fetch missing releases")
			requeueAfter := requeue.JitterFixAdditive(r.Requeue.Requeue(time.Until(rlErr.Info.Reset)), 10)
			return nil, requeueAfter, err
		}

		meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
			Type:    "FetchReleaseInformationFailed",
			Status:  metav1.ConditionFalse,
			Reason:  "Unknown",
			Message: "Fetching release information from GitHub failed",
		})
		r.Recorder.Event(data, v1.EventTypeWarning, "FetchingReleasesFailed", "Could not fetch missing releases")
		return nil, time.Duration(0), err
	}

	return releases, time.Duration(0), nil
}

func (r *RepositoryReconciler) checkReleaseDataForMissingArtifacts(data *gollumv1alpha1.Repository, releases []ReleaseArtifacts) []ReleaseArtifacts {
	releasesWithMissingArtifacts := make([]ReleaseArtifacts, 0, len(releases))
	var releaseAssetChecker ReleaseArtifactChecker = &DefaultReleaseArtifactChecker{}

	for _, release := range releases {
		tagName := release.Release.TagName
		_, found := data.Status.Releases[tagName]
		if !found {
			data.Status.Releases[tagName] = &gollumv1alpha1.Release{
				MissingArtifacts: make(map[gollumv1alpha1.ArtifactType]bool),
			}
		}

		for _, artifactType := range gollumv1alpha1.ArtifactTypes() {
			_, hasPipelineDefined := data.Spec.PipelineNames[artifactType]
			if !hasPipelineDefined {
				delete(data.Status.Releases[tagName].MissingArtifacts, artifactType)
			} else {
				validArtifacts, _ := releaseAssetChecker.HasValidArtifacts(&release, artifactType)
				data.Status.Releases[tagName].MissingArtifacts[artifactType] = !validArtifacts
			}

			if data.Status.Releases[tagName].MissingArtifacts[artifactType] {
				releasesWithMissingArtifacts = append(releasesWithMissingArtifacts, release)
			}
		}
	}

	return releasesWithMissingArtifacts
}

func (r *RepositoryReconciler) createRunsForRelease(ctx context.Context, namespace string, data *gollumv1alpha1.Repository, rel github.Release) (int, error) {
	var errs error
	startedRuns := 0

	for _, artType := range gollumv1alpha1.ArtifactTypes() {
		created, err := r.createRun(ctx, namespace, data, rel, artType)
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			startedRuns += created
		}
	}

	return startedRuns, errs
}

func (r *RepositoryReconciler) createRun(ctx context.Context, namespace string, data *gollumv1alpha1.Repository, rel github.Release, artifactType gollumv1alpha1.ArtifactType) (int, error) {
	logger := log.FromContext(ctx)

	pipelineRunRequest := tekton.BuildRunRequest(rel.TagName, namespace, data, artifactType)
	if pipelineRunRequest == nil {
		return 0, nil
	}

	// try not to spawn new pipelineruns for a pipeline if it hasn't finished
	createdRuns, found := data.Status.Releases[rel.TagName]
	if found && createdRuns != nil && createdRuns.MostRecentRuns != nil && createdRuns.MostRecentRuns[artifactType] != nil {
		pipelineRunName := createdRuns.MostRecentRuns[artifactType].Name
		pipelineRun, err := r.PipelineRunner.GetPipelineRun(ctx, namespace, pipelineRunName)
		if err != nil {
			logger.Info("Could not get PipelineRun", "pipelinerun", pipelineRunName, "error", err)
		} else {
			hasStarted := pipelineRun.Status.StartTime != nil
			hasCompleted := hasStarted && pipelineRun.Status.CompletionTime != nil
			hasStartedRecently := time.Since(pipelineRun.Status.StartTime.Time) < 60*time.Minute
			isSucceeded := false
			for _, condition := range pipelineRun.Status.Status.Conditions {
				if condition.Type == apis.ConditionSucceeded && condition.IsTrue() {
					isSucceeded = true
				}
			}

			if hasStarted && !hasCompleted && hasStartedRecently {
				logger.Info("Found previous PipelineRun that is not completed, yet", "run", pipelineRunName)
				return 0, nil
			} else if hasCompleted {
				logger.Info("Found previous PipelineRun that has completed some time ago but produced no artifacts, starting new PipelineRun", "run", pipelineRunName, "success", isSucceeded)
			}
		}
	}

	logger.Info("Creating a PipelineRun request for release", "release", rel.TagName)
	run, err := r.PipelineRunner.CreatePipelineRun(ctx, *pipelineRunRequest)
	if err != nil {
		metrics.PipelineRunCreationErrors.WithLabelValues(data.Spec.Owner, data.Spec.Repository, rel.TagName).Inc()
		return 0, err
	}

	metrics.PipelineRunsCreated.WithLabelValues(data.Spec.Owner, data.Spec.Repository, rel.TagName).Inc()
	statusRun, found := data.Status.Releases[rel.TagName]
	if !found {
		data.Status.Releases[rel.TagName] = &gollumv1alpha1.Release{
			MostRecentRuns: make(map[gollumv1alpha1.ArtifactType]*gollumv1alpha1.PipelineRun),
		}
	} else if statusRun.MostRecentRuns == nil {
		statusRun.MostRecentRuns = make(map[gollumv1alpha1.ArtifactType]*gollumv1alpha1.PipelineRun)
	}

	statusRun = data.Status.Releases[rel.TagName]
	if statusRun.MostRecentRuns[artifactType] == nil {
		statusRun.MostRecentRuns[artifactType] = &gollumv1alpha1.PipelineRun{}
	}
	statusRun.MostRecentRuns[artifactType].RunsCreated += 1
	statusRun.MostRecentRuns[artifactType].Name = run.Name
	statusRun.MostRecentRuns[artifactType].CreationTimestamp = metav1.Time{Time: time.Now()}

	r.Recorder.Event(data, v1.EventTypeNormal, "PipelineRunScheduled", fmt.Sprintf("Scheduled PipelineRun %s (#%d) for tag %s", run.Name, statusRun.MostRecentRuns[artifactType].RunsCreated, rel.TagName))
	return statusRun.MostRecentRuns[artifactType].RunsCreated, nil
}

func (r *RepositoryReconciler) fetchArtifactDataForReleases(ctx context.Context, data *gollumv1alpha1.Repository, releases []github.Release) ([]ReleaseArtifacts, time.Duration) {
	p := pool.NewWithResults[ReleaseArtifacts]().WithContext(ctx).WithMaxGoroutines(3)

	for _, release := range releases {
		p.Go(func(ctx context.Context) (ReleaseArtifacts, error) {
			ret, err := r.fetchArtifactDataForRelease(ctx, data, release)
			if err != nil {
				log.FromContext(ctx).Error(err, "could not fetch artifact for release")
			}
			if ret != nil {
				return *ret, err
			}
			return ReleaseArtifacts{}, err
		})
	}

	results, err := p.Wait()

	if errors.Is(err, github.ErrUnauthorized) {
		r.Recorder.Event(data, v1.EventTypeWarning, "UnauthorizedRequests", "Could not fetch data from GitHub Packages API")
	}

	var rateLimitErr *github.RateLimitError
	var requeueAfter time.Duration
	if errors.As(err, &rateLimitErr) {
		requeueAfter = requeue.JitterFixAdditive(r.Requeue.Requeue(time.Until(rateLimitErr.Info.Reset)), 60)
		r.Recorder.Event(data, v1.EventTypeWarning, "RateLimitExceeded", "Could not fetch artifacts from GitHub API")
	}

	meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
		Type:    "FetchReleaseArtifactsFailed",
		Status:  metav1.ConditionFalse,
		Reason:  "ApiErrors",
		Message: "Fetching release artifacts from GitHub failed",
	})

	if err != nil && !errors.Is(err, github.ErrUnauthorized) && !errors.As(err, &rateLimitErr) {
		r.Recorder.Event(data, v1.EventTypeWarning, "FetchArtifactsError", "Received non-fatal errors while fetching artifact data")
		log.FromContext(ctx).Error(err, "fetching artifact data produced error(s)")
	}

	ret := make([]ReleaseArtifacts, 0, len(results))
	for _, result := range results {
		if !result.IsEmpty() {
			ret = append(ret, result)
		}
	}

	return ret, requeueAfter
}

func (r *RepositoryReconciler) fetchArtifactDataForRelease(ctx context.Context, data *gollumv1alpha1.Repository, release github.Release) (*ReleaseArtifacts, error) {
	relWithArtifacts := &ReleaseArtifacts{
		Release: release,
	}

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	query := buildArtifactQuery(data, release)

	fatalErrChan := make(chan error)
	var err error
	_, found := data.Spec.PipelineNames[gollumv1alpha1.ArtifactsKeyReleaseAssets]
	if found {
		wg.Add(1)
		go func() {
			defer wg.Done()
			relWithArtifacts.Assets, err = r.GithubClient.GetAssets(ctx, query)
			if err != nil {
				fatalErrChan <- err
			}
		}()
	}

	_, found = data.Spec.PipelineNames[gollumv1alpha1.ArtifactsKeyPackagesContainer]
	if found {
		wg.Add(1)
		go func() {
			defer wg.Done()
			relWithArtifacts.Packages, err = r.GithubClient.GetPackages(ctx, query)
			if err != nil {
				fatalErrChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(fatalErrChan)
	}()

	err = <-fatalErrChan
	// throw away result for this release on a fatal error
	if err != nil {
		cancel()
		return nil, err
	}

	return relWithArtifacts, nil
}

func (r *RepositoryReconciler) applyVersionFilter(ctx context.Context, data *gollumv1alpha1.Repository, releases []github.Release) ([]github.Release, error) {
	if data.Spec.VersionFilter == nil && len(data.Spec.OmitVersions) == 0 {
		return releases, nil
	}
	metrics.ReleasesAvailableTotal.WithLabelValues(data.Spec.Owner, data.Spec.Repository).Set(float64(len(releases)))

	filter, err := getVersionFilter(data.Spec.VersionFilter)
	if err != nil {
		return nil, err
	}

	var filteredReleases []github.Release //nolint prealloc
	var errs error
	for _, rel := range releases {
		matches, err := filter.Matches(rel.TagName)
		if err != nil || !matches {
			errs = multierror.Append(errs, err)
			continue
		}

		if slices.Contains(data.Spec.OmitVersions, rel.TagName) {
			continue
		}

		filteredReleases = append(filteredReleases, rel)
	}

	filteredReleaseCount := len(releases)
	if len(releases) != filteredReleaseCount {
		GollumReasonFilteredReleases := "ReleasesFiltered"
		log.FromContext(ctx).Info("Filtered releases", "amount-filtered", filteredReleaseCount)
		r.Recorder.Event(data, v1.EventTypeNormal, GollumReasonFilteredReleases, fmt.Sprintf("Got %d releases, %d releases do not match the release filter", len(filteredReleases), len(releases)-filteredReleaseCount))
	}
	metrics.FilteredReleasesTotal.WithLabelValues(data.Spec.Owner, data.Spec.Repository).Set(float64(filteredReleaseCount))
	return filteredReleases, errs
}

func (r *RepositoryReconciler) createPipelineRunsForReleases(ctx context.Context, data *gollumv1alpha1.Repository, releases []ReleaseArtifacts, namespace string) (time.Duration, error) {
	var err error
	var maxPreviouslyCreatedRuns int
	var runsCreated int

	for _, rel := range releases {
		totalRunsCreatedForVersion, createRunErr := r.createRunsForRelease(ctx, namespace, data, rel.Release)
		if createRunErr != nil {
			err = multierror.Append(err, createRunErr)
		} else {
			runsCreated++
			if totalRunsCreatedForVersion > maxPreviouslyCreatedRuns {
				maxPreviouslyCreatedRuns = totalRunsCreatedForVersion
			}
		}
	}

	if runsCreated > 0 {
		// Some runs created successfully
		message := "Created PipelineRuns"
		if err != nil {
			message += ", produced error(s)"
		}

		meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
			Type:    "PipelineRunsCreated",
			Status:  metav1.ConditionTrue,
			Message: message,
			Reason:  "MissingReleasesFound",
		})

		// Even if err != nil, as some runs could be created, it's likely a transient problem, let's requeue the object
		// normally
		return time.Duration(0), nil
	}

	if runsCreated == 0 && err != nil {
		meta.SetStatusCondition(data.GetConditions(), metav1.Condition{
			Type:    "PipelineRunsCreationFailed",
			Status:  metav1.ConditionFalse,
			Message: "No PipelineRun(s) could be created",
			Reason:  "Unknown",
		})
	}

	// No runs could be created, let's use the amount of maxPreviouslyCreatedRuns to calculate an exponential backoff duration
	return requeue.JitterFixAdditive(r.Requeue.Requeue(requeue.ExponentialBackoff(maxPreviouslyCreatedRuns, 5)), 120), err
}

// SetupWithManager sets up the controller with the Manager.
func (r *RepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Recorder = mgr.GetEventRecorderFor("repository-controller")
	return ctrl.NewControllerManagedBy(mgr).
		For(&gollumv1alpha1.Repository{}).
		Named("repository").
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
