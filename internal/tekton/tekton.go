package tekton

import (
	"context"
	"errors"
	"fmt"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	valSecret           = "secret"
	keyType             = "type"
	valVolume           = "volume"
	keyStorageClassName = "storageClassName"
	keySecretName       = "secretName"
)

var (
	ErrTektonCreatePipelineRunJobUnauthorized = errors.New("unauthorized to create tekton pipelinerun resource")
	ErrTektonPipelineNotFound                 = errors.New("tekton pipeline not found")
	ErrTektonGetPipelineForbidden             = errors.New("forbidden to get tekton pipeline")
	ErrTektonInvalidPipelineRunSpec           = errors.New("can not build PipelineRun from gollum spec")
)

type TektonPipelineRunner struct {
	client *versioned.Clientset
}

func NewTektonPipelineRunner(client *versioned.Clientset) (*TektonPipelineRunner, error) {
	return &TektonPipelineRunner{
		client: client,
	}, nil
}

func (t *TektonPipelineRunner) GetPipeline(ctx context.Context, namespace, name string) (*pipelinev1.Pipeline, error) {
	opts := metav1.GetOptions{}
	run, err := t.client.TektonV1().Pipelines(namespace).Get(ctx, name, opts)
	if err == nil {
		return run, nil
	}

	if k8serrors.IsNotFound(err) {
		return nil, ErrTektonPipelineNotFound
	}

	if k8serrors.IsForbidden(err) {
		return nil, ErrTektonGetPipelineForbidden
	}

	return nil, err
}

func (t *TektonPipelineRunner) GetPipelineRun(ctx context.Context, namespace, name string) (*pipelinev1.PipelineRun, error) {
	opts := metav1.GetOptions{}
	run, err := t.client.TektonV1().PipelineRuns(namespace).Get(ctx, name, opts)
	if err == nil {
		return run, nil
	}

	if k8serrors.IsNotFound(err) {
		return nil, ErrTektonPipelineNotFound
	}

	if k8serrors.IsForbidden(err) {
		return nil, ErrTektonGetPipelineForbidden
	}

	return nil, err
}

func (t *TektonPipelineRunner) CreatePipelineRun(ctx context.Context, req CreatePipelineRunRequest) (*pipelinev1.PipelineRun, error) {
	pipelineSpec, err := getPipelineRunSpec(req)
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", ErrTektonInvalidPipelineRunSpec, err)
	}

	pipelineRun := &pipelinev1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    req.Namespace,
			GenerateName: fmt.Sprintf("%s-", req.PipelineRunName),
		},
		Spec: *pipelineSpec,
	}

	run, err := t.client.TektonV1().PipelineRuns(req.Namespace).Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsForbidden(err) {
			return nil, ErrTektonCreatePipelineRunJobUnauthorized
		}
	}

	return run, err
}
