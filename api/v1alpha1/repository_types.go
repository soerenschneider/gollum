package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ArtifactType string

const (
	ArtifactsKeyReleaseAssets     ArtifactType = "assets"
	ArtifactsKeyPackagesContainer ArtifactType = "container"
)

func ArtifactTypes() []ArtifactType {
	return []ArtifactType{
		ArtifactsKeyReleaseAssets,
		ArtifactsKeyPackagesContainer,
	}
}

// RepositorySpec defines the desired state of Repository.
type RepositorySpec struct {
	Owner         string `json:"owner"`
	Repository    string `json:"repo"`
	CloneUsingSsh bool   `json:"cloneUsingSsh"`

	// +kubebuilder:default:=true
	MemorizeReleases bool `json:"memorizeReleases"`

	PipelineRunName string `json:"pipelineRunName"`

	PipelineNames map[ArtifactType]string `json:"PipelineNames"`

	VersionFilter *VersionFilterSpec           `json:"versionFilter,omitempty"`
	Workspaces    map[string]map[string]string `json:"workspaces"`
}

type VersionFilterSpec struct {
	// +kubebuilder:validation:Enum=semver
	Impl string `json:"impl"`
	Arg  string `json:"arg"`
}

// RepositoryStatus defines the observed state of Repository.
type RepositoryStatus struct {
	Ready      bool                `json:"ready"`
	Releases   map[string]*Release `json:"releases"`
	Conditions []metav1.Condition  `json:"conditions,omitempty"`
	LastCheck  *metav1.Time        `json:"lastCheck"`
}

type Release struct {
	MostRecentRuns   map[ArtifactType]*PipelineRun `json:"pipelineRuns,omitempty"`
	MissingArtifacts map[ArtifactType]bool         `json:"missingArtifacts"`
}

type PipelineRun struct {
	Name              string      `json:"name"`
	CreationTimestamp metav1.Time `json:"timestamp,omitempty"`
	RunsCreated       int         `json:"runsCreated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Owner",type=string,JSONPath=`.spec.owner`
// +kubebuilder:printcolumn:name="Repo",type=string,JSONPath=`.spec.repo`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.status)].status`
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec,omitempty"`
	Status RepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RepositoryList contains a list of Repository.
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}

func (r *Repository) GetConditions() *[]metav1.Condition {
	// We only want to keep the most recent condition entry
	r.Status.Conditions = make([]metav1.Condition, 0)
	return &r.Status.Conditions
}
