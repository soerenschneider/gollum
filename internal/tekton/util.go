package tekton

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	gollumv1alpha1 "github.com/soerenschneider/gollum/api/v1alpha1"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	DefaultPipelineName = "build-gh-release"
	ArgCloneUrl         = "clone-url"
	ArgOwner            = "owner"
	ArgRepo             = "repository"
	ArgRevision         = "revision"
	DefaultRevision     = ""
)

func BuildRunRequest(tag string, namespace string, data *gollumv1alpha1.Repository, artifactType gollumv1alpha1.ArtifactType) *CreatePipelineRunRequest {
	pipelineName, found := data.Spec.PipelineNames[artifactType]
	if !found || len(pipelineName) == 0 {
		return nil
	}

	pipelineRunName := cmp.Or(data.Spec.PipelineRunName, fmt.Sprintf("gollum-%s-%s-%s", safeSlice(data.Spec.Owner, 5), safeSlice(data.Spec.Repository, 5), tag))
	return &CreatePipelineRunRequest{
		Namespace:       namespace,
		PipelineRunName: pipelineRunName,
		PipelineName:    pipelineName,
		Params: map[string]string{
			ArgCloneUrl: GetRepoUrl(data.Spec.CloneUsingSsh, data.Spec.Owner, data.Spec.Repository),
			ArgRevision: tag,
			ArgOwner:    data.Spec.Owner,
			ArgRepo:     data.Spec.Repository,
		},
		WorkspaceBindings: data.Spec.Workspaces,
	}
}

func safeSlice(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func GetRepoUrl(sshCheckout bool, owner, repo string) string {
	if sshCheckout {
		return fmt.Sprintf("git@github.com:%s/%s.git", owner, repo)
	}
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

func getPipelineRunSpec(req CreatePipelineRunRequest) (*pipelinev1.PipelineRunSpec, error) {
	workspaces, err := getWorkspaceBindings(req)
	if err != nil {
		return nil, fmt.Errorf("could not build workspaces spec: %w", err)
	}

	return &pipelinev1.PipelineRunSpec{
		PipelineRef: &pipelinev1.PipelineRef{
			Name: req.PipelineName,
		},
		Workspaces: workspaces,
		Params:     getParams(req),
	}, nil
}

func getWorkspaceBindings(req CreatePipelineRunRequest) ([]pipelinev1.WorkspaceBinding, error) {
	ret := make([]pipelinev1.WorkspaceBinding, 0, len(req.WorkspaceBindings))
	var errs error

	for key, val := range req.WorkspaceBindings {
		binding := pipelinev1.WorkspaceBinding{
			Name: key,
		}

		switch value := val[keyType]; value {
		case valSecret:
			secretName := val[keySecretName]
			if strings.TrimSpace(secretName) == "" {
				errs = multierror.Append(errs, fmt.Errorf("type secret is missing %q", keySecretName))
			}
			binding.Secret = &v1.SecretVolumeSource{
				SecretName: secretName,
			}
		case valVolume:
			var storageClassName *string
			if val[keyStorageClassName] != "" {
				name := val[keyStorageClassName]
				storageClassName = &name
			}

			binding.VolumeClaimTemplate = &v1.PersistentVolumeClaim{
				Spec: v1.PersistentVolumeClaimSpec{
					StorageClassName: storageClassName,
					Resources: v1.VolumeResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceStorage: *resource.NewQuantity(1, resource.BinarySI),
						},
					},
					AccessModes: []v1.PersistentVolumeAccessMode{
						v1.ReadWriteOnce,
					},
				},
			}
		default:
			errs = multierror.Append(errs, fmt.Errorf("unknown value: %q", value))
		}

		ret = append(ret, binding)
	}

	return ret, errs
}

func getParams(req CreatePipelineRunRequest) []pipelinev1.Param {
	ret := make([]pipelinev1.Param, 0, len(req.Params))

	for key, val := range req.Params {
		ret = append(ret, pipelinev1.Param{
			Name: key,
			Value: pipelinev1.ParamValue{
				StringVal: val,
				Type:      pipelinev1.ParamTypeString,
			},
		})
	}

	return ret
}
