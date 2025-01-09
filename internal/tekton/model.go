package tekton

type CreatePipelineRunRequest struct {
	Namespace       string
	PipelineRunName string
	PipelineName    string

	Params            map[string]string
	WorkspaceBindings map[string]map[string]string
}
