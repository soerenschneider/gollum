# gollum

## Overview
Gollum is a Kubernetes operator designed to monitor GitHub release assets. If the required assets for a specified version are missing, Gollum triggers a Tekton `PipelineRun` to generate them.

## Features
- Monitors GitHub releases for specified repositories.
- Detects missing release assets for a given version.
- Triggers a Tekton `PipelineRun` to build and publish the missing assets.
- Configurable through Kubernetes Custom Resources.

## Installation
### Prerequisites
- Kubernetes cluster (v1.20+ recommended)
- Tekton installed in the cluster
- GitHub API access token (if monitoring private repositories)

### Deploying Gollum
1. Clone the repository:
   ```sh
   git clone https://github.com/your-org/gollum.git
   cd gollum
   ```
2. Apply the CRDs and Operator:
   ```sh
   kubectl apply -f deploy/crds/
   kubectl apply -f deploy/operator.yaml
   ```
3. Verify installation:
   ```sh
   kubectl get pods -n gollum
   ```
   Ensure the operator is running properly.

## Usage
### Defining a Custom Resource
To monitor a repository, create a `GollumReleaseMonitor` Custom Resource:

```yaml
apiVersion: gollum.soeren.cloud/v1alpha1
kind: Repository
metadata:
   labels:
      app.kubernetes.io/name: gollum
      app.kubernetes.io/managed-by: kustomize
   name: soerenschneider-tunnelguard
spec:
   owner: "soerenschneider"
   repo: "tunnelguard"
   cloneUsingSsh: false
   pipelineRunName: "gollum"
   pipelineNames:
      assets: "build-gh-release"
   versionFilter:
      impl: "semver"
      arg: ">= v1.0.0"
   workspaces:
      signify:
         type: "secret"
         secretName: "signify"
      shared-data:
         type: "volume"
         storageClassName: "openebs-hostpath"
      github-token:
         type: "secret"
         secretName: "github"
```

Apply it using:
```sh
kubectl apply -f example-repo-monitor.yaml
```

### How It Works
1. Gollum continuously monitors the specified GitHub repository for new releases.
2. If a release is missing any required assets, Gollum triggers the specified Tekton `PipelineRun`.
3. Tekton builds and uploads the missing assets to the release.
4. Gollum updates the status of the `GollumReleaseMonitor` resource.

## Configuration
- **GitHub Authentication**: Use a Kubernetes secret to store a GitHub personal access token (PAT) for private repositories.
- **Tekton Integration**: Specify an existing Tekton pipeline reference in the CR.
- **Polling Interval**: Configure how frequently Gollum checks GitHub releases.

## Development
### Running Locally
```sh
go run main.go
```
### Building the Docker Image
```sh
docker build -t your-org/gollum:latest .
```
### Deploying the Operator
```sh
kubectl apply -f deploy/
```
